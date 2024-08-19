package service

import (
	"errors"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data/dto"
	"github.com/leaf-rain/raindata/app_bi/internal/data/entity"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
)

var ErrRoleExistence = errors.New("存在相同角色id")

//@function: CreateAuthority
//@description: 创建一个角色
//@param: auth data.SysAuthority
//@return: authority data.SysAuthority, err error

type AuthorityService struct {
	data *entity.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

var AuthorityServiceApp = new(AuthorityService)

func (svc *AuthorityService) CreateAuthority(auth entity.SysAuthority) (authority entity.SysAuthority, err error) {

	if err = svc.data.SqlClient.Where("authority_id = ?", auth.AuthorityId).First(&entity.SysAuthority{}).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return auth, ErrRoleExistence
	}

	e := svc.data.SqlClient.Transaction(func(tx *gorm.DB) error {

		if err = tx.Create(&auth).Error; err != nil {
			return err
		}

		auth.SysBaseMenus = dto.DefaultMenu()
		if err = tx.Model(&auth).Association("SysBaseMenus").Replace(&auth.SysBaseMenus); err != nil {
			return err
		}
		casbinInfos := dto.DefaultCasbin()
		authorityId := strconv.Itoa(int(auth.AuthorityId))
		rules := [][]string{}
		for _, v := range casbinInfos {
			rules = append(rules, []string{authorityId, v.Path, v.Method})
		}
		return CasbinServiceApp.AddPolicies(tx, rules)
	})

	return auth, e
}

//@function: CopyAuthority
//@description: 复制一个角色
//@param: copyInfo dto.SysAuthorityCopyResponse
//@return: authority data.SysAuthority, err error

func (svc *AuthorityService) CopyAuthority(copyInfo dto.SysAuthorityCopyResponse) (authority entity.SysAuthority, err error) {
	var authorityBox entity.SysAuthority
	if !errors.Is(svc.data.SqlClient.Where("authority_id = ?", copyInfo.Authority.AuthorityId).First(&authorityBox).Error, gorm.ErrRecordNotFound) {
		return authority, ErrRoleExistence
	}
	copyInfo.Authority.Children = []entity.SysAuthority{}
	menus, err := MenuServiceApp.GetMenuAuthority(&rhttp.GetAuthorityId{AuthorityId: copyInfo.OldAuthorityId})
	if err != nil {
		return
	}
	var baseMenu []entity.SysBaseMenu
	for _, v := range menus {
		intNum := v.MenuId
		v.SysBaseMenu.ID = uint(intNum)
		baseMenu = append(baseMenu, v.SysBaseMenu)
	}
	copyInfo.Authority.SysBaseMenus = baseMenu
	err = svc.data.SqlClient.Create(&copyInfo.Authority).Error
	if err != nil {
		return
	}

	var btns []entity.SysAuthorityBtn

	err = svc.data.SqlClient.Find(&btns, "authority_id = ?", copyInfo.OldAuthorityId).Error
	if err != nil {
		return
	}
	if len(btns) > 0 {
		for i := range btns {
			btns[i].AuthorityId = copyInfo.Authority.AuthorityId
		}
		err = svc.data.SqlClient.Create(&btns).Error

		if err != nil {
			return
		}
	}
	paths := CasbinServiceApp.GetPolicyPathByAuthorityId(copyInfo.OldAuthorityId)
	err = CasbinServiceApp.UpdateCasbin(copyInfo.Authority.AuthorityId, paths)
	if err != nil {
		_ = AuthorityService.DeleteAuthority(&copyInfo.Authority)
	}
	return copyInfo.Authority, err
}

//@function: UpdateAuthority
//@description: 更改一个角色
//@param: auth data.SysAuthority
//@return: authority data.SysAuthority, err error

func (svc *AuthorityService) UpdateAuthority(auth entity.SysAuthority) (authority entity.SysAuthority, err error) {
	var oldAuthority entity.SysAuthority
	err = svc.data.SqlClient.Where("authority_id = ?", auth.AuthorityId).First(&oldAuthority).Error
	if err != nil {
		svc.log.Debug(err.Error())
		return entity.SysAuthority{}, errors.New("查询角色数据失败")
	}
	err = svc.data.SqlClient.Model(&oldAuthority).Updates(&auth).Error
	return auth, err
}

//@function: DeleteAuthority
//@description: 删除角色
//@param: auth *data.SysAuthority
//@return: err error

func (svc *AuthorityService) DeleteAuthority(auth *entity.SysAuthority) error {
	if errors.Is(svc.data.SqlClient.Debug().Preload("Users").First(&auth).Error, gorm.ErrRecordNotFound) {
		return errors.New("该角色不存在")
	}
	if len(auth.Users) != 0 {
		return errors.New("此角色有用户正在使用禁止删除")
	}
	if !errors.Is(svc.data.SqlClient.Where("authority_id = ?", auth.AuthorityId).First(&entity.SysUser{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("此角色有用户正在使用禁止删除")
	}
	if !errors.Is(svc.data.SqlClient.Where("parent_id = ?", auth.AuthorityId).First(&entity.SysAuthority{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("此角色存在子角色不允许删除")
	}

	return svc.data.SqlClient.Transaction(func(tx *gorm.DB) error {
		var err error
		if err = tx.Preload("SysBaseMenus").Preload("DataAuthorityId").Where("authority_id = ?", auth.AuthorityId).First(auth).Unscoped().Delete(auth).Error; err != nil {
			return err
		}

		if len(auth.SysBaseMenus) > 0 {
			if err = tx.Model(auth).Association("SysBaseMenus").Delete(auth.SysBaseMenus); err != nil {
				return err
			}
			// err = db.Association("SysBaseMenus").Delete(&auth)
		}
		if len(auth.DataAuthorityId) > 0 {
			if err = tx.Model(auth).Association("DataAuthorityId").Delete(auth.DataAuthorityId); err != nil {
				return err
			}
		}

		if err = tx.Delete(&entity.SysUserAuthority{}, "sys_authority_authority_id = ?", auth.AuthorityId).Error; err != nil {
			return err
		}
		if err = tx.Where("authority_id = ?", auth.AuthorityId).Delete(&[]entity.SysAuthorityBtn{}).Error; err != nil {
			return err
		}

		authorityId := strconv.Itoa(int(auth.AuthorityId))

		if err = CasbinServiceApp.RemoveFilteredPolicy(tx, authorityId); err != nil {
			return err
		}

		return nil
	})
}

//@function: GetAuthorityInfoList
//@description: 分页获取数据
//@param: info dto.PageInfo
//@return: list interface{}, total int64, err error

func (svc *AuthorityService) GetAuthorityInfoList(info rhttp.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := svc.data.SqlClient.Model(&entity.SysAuthority{})
	if err = db.Where("parent_id = ?", "0").Count(&total).Error; total == 0 || err != nil {
		return
	}
	var authority []entity.SysAuthority
	err = db.Limit(limit).Offset(offset).Preload("DataAuthorityId").Where("parent_id = ?", "0").Find(&authority).Error
	for k := range authority {
		err = AuthorityService.findChildrenAuthority(&authority[k])
	}
	return authority, total, err
}

//@function: GetAuthorityInfo
//@description: 获取所有角色信息
//@param: auth data.SysAuthority
//@return: sa data.SysAuthority, err error

func (svc *AuthorityService) GetAuthorityInfo(auth entity.SysAuthority) (sa entity.SysAuthority, err error) {
	err = svc.data.SqlClient.Preload("DataAuthorityId").Where("authority_id = ?", auth.AuthorityId).First(&sa).Error
	return sa, err
}

//@function: SetDataAuthority
//@description: 设置角色资源权限
//@param: auth data.SysAuthority
//@return: error

func (svc *AuthorityService) SetDataAuthority(auth entity.SysAuthority) error {
	var s entity.SysAuthority
	svc.data.SqlClient.Preload("DataAuthorityId").First(&s, "authority_id = ?", auth.AuthorityId)
	err := svc.data.SqlClient.Model(&s).Association("DataAuthorityId").Replace(&auth.DataAuthorityId)
	return err
}

//@function: SetMenuAuthority
//@description: 菜单与角色绑定
//@param: auth *data.SysAuthority
//@return: error

func (svc *AuthorityService) SetMenuAuthority(auth *entity.SysAuthority) error {
	var s entity.SysAuthority
	svc.data.SqlClient.Preload("SysBaseMenus").First(&s, "authority_id = ?", auth.AuthorityId)
	err := svc.data.SqlClient.Model(&s).Association("SysBaseMenus").Replace(&auth.SysBaseMenus)
	return err
}

//@function: findChildrenAuthority
//@description: 查询子角色
//@param: authority *data.SysAuthority
//@return: err error

func (svc *AuthorityService) findChildrenAuthority(authority *entity.SysAuthority) (err error) {
	err = svc.data.SqlClient.Preload("DataAuthorityId").Where("parent_id = ?", authority.AuthorityId).Find(&authority.Children).Error
	if len(authority.Children) > 0 {
		for k := range authority.Children {
			err = AuthorityService.findChildrenAuthority(&authority.Children[k])
		}
	}
	return err
}
