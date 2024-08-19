package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/leaf-rain/raindata/app_bi/internal/data/entity"
	"github.com/leaf-rain/raindata/app_bi/third_party/hash"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"github.com/leaf-rain/raindata/app_bi/third_party/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type UserService struct {
	log  *zap.Logger
	data *entity.Data
}

// NewUserService new a greeter service.
func NewUserService(logger *zap.Logger, data *entity.Data) *UserService {
	return &UserService{log: logger}
}

func (svc *UserService) Register(u entity.SysUser) (userInter entity.SysUser, err error) {
	var user entity.SysUser
	if !errors.Is(svc.data.SqlClient.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return userInter, errors.New("用户名已注册")
	}
	// 否则 附加uuid 密码hash加密 注册
	u.Password = hash.BcryptHash(u.Password)
	u.UUID = uuid.Must(uuid.NewV6())
	err = svc.data.SqlClient.Create(&u).Error
	return u, err
}

//@function: Login
//@description: 用户登录
//@param: u *data.SysUser
//@return: err error, userInter *data.SysUser

func (svc *UserService) Login(u *entity.SysUser) (userInter *entity.SysUser, err error) {
	if nil == svc.data.SqlClient {
		return nil, fmt.Errorf("db not init")
	}
	var entityUser = entity.NewEntityUser(u, svc.data)
	err = entityUser.ReloadByDb()
	if err == nil {
		MenuServiceApp.UserAuthorityDefaultRouter(&user)
	}
	return &user, err
}

//@function: ChangePassword
//@description: 修改用户密码
//@param: u *data.SysUser, newPassword string
//@return: userInter *data.SysUser,err error

func (svc *UserService) ChangePassword(u *entity.SysUser, newPassword string) (userInter *entity.SysUser, err error) {
	var user entity.SysUser
	if err = svc.data.SqlClient.Where("id = ?", u.ID).First(&user).Error; err != nil {
		return nil, err
	}
	if ok := hash.BcryptCheck(u.Password, user.Password); !ok {
		return nil, errors.New("原密码错误")
	}
	user.Password = hash.BcryptHash(newPassword)
	err = svc.data.SqlClient.Save(&user).Error
	return &user, err

}

//@function: GetUserInfoList
//@description: 分页获取数据
//@param: info dto.PageInfo
//@return: err error, list interface{}, total int64

func (svc *UserService) GetUserInfoList(info rhttp.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := svc.data.SqlClient.Model(&entity.SysUser{})
	var userList []entity.SysUser
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Preload("Authorities").Preload("Authority").Find(&userList).Error
	return userList, total, err
}

//@function: SetUserAuthority
//@description: 设置一个用户的权限
//@param: uuid uuid.UUID, authorityId string
//@return: err error

func (svc *UserService) SetUserAuthority(id uint, authorityId uint) (err error) {
	assignErr := svc.data.SqlClient.Where("sys_user_id = ? AND sys_authority_authority_id = ?", id, authorityId).First(&entity.SysUserAuthority{}).Error
	if errors.Is(assignErr, gorm.ErrRecordNotFound) {
		return errors.New("该用户无此角色")
	}
	err = svc.data.SqlClient.Model(&entity.SysUser{}).Where("id = ?", id).Update("authority_id", authorityId).Error
	return err
}

//@function: SetUserAuthorities
//@description: 设置一个用户的权限
//@param: id uint, authorityIds []string
//@return: err error

func (svc *UserService) SetUserAuthorities(id uint, authorityIds []uint) (err error) {
	return svc.data.SqlClient.Transaction(func(tx *gorm.DB) error {
		var user entity.SysUser
		TxErr := tx.Where("id = ?", id).First(&user).Error
		if TxErr != nil {
			svc.log.Error("[SetUserAuthorities]查询用户数据失败", zap.Error(TxErr))
			return errors.New("查询用户数据失败")
		}
		TxErr = tx.Delete(&[]entity.SysUserAuthority{}, "sys_user_id = ?", id).Error
		if TxErr != nil {
			return TxErr
		}
		var useAuthority []entity.SysUserAuthority
		for _, v := range authorityIds {
			useAuthority = append(useAuthority, entity.SysUserAuthority{
				SysUserId: id, SysAuthorityAuthorityId: v,
			})
		}
		TxErr = tx.Create(&useAuthority).Error
		if TxErr != nil {
			return TxErr
		}
		TxErr = tx.Model(&user).Update("authority_id", authorityIds[0]).Error
		if TxErr != nil {
			return TxErr
		}
		// 返回 nil 提交事务
		return nil
	})
}

//@function: DeleteUser
//@description: 删除用户
//@param: id float64
//@return: err error

func (svc *UserService) DeleteUser(id int) (err error) {
	return svc.data.SqlClient.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&entity.SysUser{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&[]entity.SysUserAuthority{}, "sys_user_id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}

//@function: SetUserInfo
//@description: 设置用户信息
//@param: reqUser data.SysUser
//@return: err error, user data.SysUser

func (svc *UserService) SetUserInfo(req entity.SysUser) error {
	return svc.data.SqlClient.Model(&entity.SysUser{}).
		Select("updated_at", "nick_name", "header_img", "phone", "email", "sideMode", "enable").
		Where("id=?", req.ID).
		Updates(map[string]interface{}{
			"updated_at": time.Now(),
			"nick_name":  req.NickName,
			"header_img": req.HeaderImg,
			"phone":      req.Phone,
			"email":      req.Email,
			"side_mode":  req.SideMode,
			"enable":     req.Enable,
		}).Error
}

//@function: SetSelfInfo
//@description: 设置用户信息
//@param: reqUser data.SysUser
//@return: err error, user data.SysUser

func (svc *UserService) SetSelfInfo(req entity.SysUser) error {
	return svc.data.SqlClient.Model(&entity.SysUser{}).
		Where("id=?", req.ID).
		Updates(req).Error
}

//@function: GetUserInfo
//@description: 获取用户信息
//@param: uuid uuid.UUID
//@return: err error, user data.SysUser

func (svc *UserService) GetUserInfo(uuid uuid.UUID) (user entity.SysUser, err error) {
	var reqUser entity.SysUser
	err = svc.data.SqlClient.Preload("Authorities").Preload("Authority").First(&reqUser, "uuid = ?", uuid).Error
	if err != nil {
		return reqUser, err
	}
	MenuServiceApp.UserAuthorityDefaultRouter(&reqUser)
	return reqUser, err
}

//@function: FindUserById
//@description: 通过id获取用户信息
//@param: id int
//@return: err error, user *data.SysUser

func (svc *UserService) FindUserById(id int) (user *entity.SysUser, err error) {
	var u entity.SysUser
	err = svc.data.SqlClient.Where("id = ?", id).First(&u).Error
	return &u, err
}

//@function: FindUserByUuid
//@description: 通过uuid获取用户信息
//@param: uuid string
//@return: err error, user *data.SysUser

func (svc *UserService) FindUserByUuid(uuid string) (user *entity.SysUser, err error) {
	var u entity.SysUser
	if err = svc.data.SqlClient.Where("uuid = ?", uuid).First(&u).Error; err != nil {
		return &u, errors.New("用户不存在")
	}
	return &u, nil
}

//@function: ResetPassword
//@description: 修改用户密码
//@param: ID uint
//@return: err error

func (svc *UserService) ResetPassword(ID uint) (err error) {
	err = svc.data.SqlClient.Model(&entity.SysUser{}).Where("id = ?", ID).Update("password", utils.BcryptHash("123456")).Error
	return err
}
