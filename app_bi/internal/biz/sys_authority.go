package biz

import (
	"errors"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"github.com/leaf-rain/raindata/app_bi/internal/data/dto"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"gorm.io/gorm"
	"strconv"
)

var ErrRoleExistence = errors.New("存在相同角色id")

//@function: CreateAuthority
//@description: 创建一个角色
//@param: auth data.SysAuthority
//@return: authority data.SysAuthority, err error

type Authority struct {
	*Business
}

func NewAuthority(biz *Business) *Authority {
	return &Authority{
		biz,
	}
}

func (b *Authority) CreateAuthority(auth data.SysAuthority) (authority data.SysAuthority, err error) {

	if err = b.data.SqlClient.Where("authority_id = ?", auth.AuthorityId).First(&data.SysAuthority{}).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return auth, ErrRoleExistence
	}

	e := b.data.SqlClient.Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(&auth).Error; err != nil {
			return err
		}
		casbinInfos := dto.DefaultCasbin()
		authorityId := strconv.Itoa(int(auth.AuthorityId))
		rules := [][]string{}
		for _, v := range casbinInfos {
			rules = append(rules, []string{authorityId, v.Path, v.Method})
		}
		cs := NewCasbin(b.Business)
		return cs.AddPolicies(tx, rules)
	})

	return auth, e
}

//@function: UpdateAuthority
//@description: 更改一个角色
//@param: auth data.SysAuthority
//@return: authority data.SysAuthority, err error

func (b *Authority) UpdateAuthority(auth data.SysAuthority) (authority data.SysAuthority, err error) {
	var oldAuthority data.SysAuthority
	err = b.data.SqlClient.Where("authority_id = ?", auth.AuthorityId).First(&oldAuthority).Error
	if err != nil {
		b.logger.Debug(err.Error())
		return data.SysAuthority{}, errors.New("查询角色数据失败")
	}
	err = b.data.SqlClient.Model(&oldAuthority).Updates(&auth).Error
	return auth, err
}

//@function: DeleteAuthority
//@description: 删除角色
//@param: auth *data.SysAuthority
//@return: err error

func (b *Authority) DeleteAuthority(auth *data.SysAuthority) error {
	if errors.Is(b.data.SqlClient.Debug().Preload("Users").First(&auth).Error, gorm.ErrRecordNotFound) {
		return errors.New("该角色不存在")
	}
	if len(auth.Users) != 0 {
		return errors.New("此角色有用户正在使用禁止删除")
	}
	if !errors.Is(b.data.SqlClient.Where("authority_id = ?", auth.AuthorityId).First(&data.SysUser{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("此角色有用户正在使用禁止删除")
	}
	if !errors.Is(b.data.SqlClient.Where("parent_id = ?", auth.AuthorityId).First(&data.SysAuthority{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("此角色存在子角色不允许删除")
	}

	return b.data.SqlClient.Transaction(func(tx *gorm.DB) error {
		var err error
		if err = tx.Preload("SysBaseMenus").Preload("DataAuthorityId").Where("authority_id = ?", auth.AuthorityId).First(auth).Unscoped().Delete(auth).Error; err != nil {
			return err
		}
		if len(auth.DataAuthorityId) > 0 {
			if err = tx.Model(auth).Association("DataAuthorityId").Delete(auth.DataAuthorityId); err != nil {
				return err
			}
		}

		if err = tx.Delete(&data.SysUserAuthority{}, "sys_authority_authority_id = ?", auth.AuthorityId).Error; err != nil {
			return err
		}
		authorityId := strconv.Itoa(int(auth.AuthorityId))
		cs := NewCasbin(b.Business)
		if err = cs.RemoveFilteredPolicy(tx, authorityId); err != nil {
			return err
		}

		return nil
	})
}

//@function: GetAuthorityInfoList
//@description: 分页获取数据
//@param: info dto.PageInfo
//@return: list interface{}, total int64, err error

func (b *Authority) GetAuthorityInfoList(info rhttp.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := b.data.SqlClient.Model(&data.SysAuthority{})
	if err = db.Where("parent_id = ?", "0").Count(&total).Error; total == 0 || err != nil {
		return
	}
	var authority []data.SysAuthority
	err = db.Limit(limit).Offset(offset).Preload("DataAuthorityId").Where("parent_id = ?", "0").Find(&authority).Error
	for k := range authority {
		err = b.findChildrenAuthority(&authority[k])
	}
	return authority, total, err
}

//@function: GetAuthorityInfo
//@description: 获取所有角色信息
//@param: auth data.SysAuthority
//@return: sa data.SysAuthority, err error

func (b *Authority) GetAuthorityInfo(auth data.SysAuthority) (sa data.SysAuthority, err error) {
	err = b.data.SqlClient.Preload("DataAuthorityId").Where("authority_id = ?", auth.AuthorityId).First(&sa).Error
	return sa, err
}

//@function: SetDataAuthority
//@description: 设置角色资源权限
//@param: auth data.SysAuthority
//@return: error

func (b *Authority) SetDataAuthority(auth data.SysAuthority) error {
	var s data.SysAuthority
	b.data.SqlClient.Preload("DataAuthorityId").First(&s, "authority_id = ?", auth.AuthorityId)
	err := b.data.SqlClient.Model(&s).Association("DataAuthorityId").Replace(&auth.DataAuthorityId)
	return err
}

//@function: findChildrenAuthority
//@description: 查询子角色
//@param: authority *data.SysAuthority
//@return: err error

func (b *Authority) findChildrenAuthority(authority *data.SysAuthority) (err error) {
	err = b.data.SqlClient.Preload("DataAuthorityId").Where("parent_id = ?", authority.AuthorityId).Find(&authority.Children).Error
	if len(authority.Children) > 0 {
		for k := range authority.Children {
			err = b.findChildrenAuthority(&authority.Children[k])
		}
	}
	return err
}
