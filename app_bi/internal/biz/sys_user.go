package biz

import (
	"errors"
	"github.com/google/uuid"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"github.com/leaf-rain/raindata/app_bi/third_party/hash"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"github.com/leaf-rain/raindata/app_bi/third_party/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type User struct {
	*Business
}

// NewUserService new a greeter service.
func NewUser(biz *Business) *User {
	return &User{biz}
}

func (b *User) Register(u data.SysUser) (userInter data.SysUser, err error) {
	var et = data.NewEntitySysUser(b.data)
	// 否则 附加uuid 密码hash加密 注册
	u.Password = hash.BcryptHash(u.Password)
	u.UUID = uuid.Must(uuid.NewV7())
	et.Model = &u
	err = et.CreateUser()
	if err != nil {
		b.logger.Error("[Register] 创建用户失败", zap.Error(err), zap.String("username", u.Username))
	}
	return u, err
}

//@function: Login
//@description: 用户登录
//@param: u *data.SysUser
//@return: err error, userInter *data.SysUser

func (b *User) Login(u *data.SysUser) (userInter *data.SysUser, err error) {
	var entityUser = data.NewEntitySysUser(b.data)
	entityUser.Model = u
	err = entityUser.ReloadByDb()
	if err != nil {
		b.logger.Error("[Register] 用户登陆失败", zap.Error(err), zap.String("username", u.Username))
	}
	return u, err
}

//@function: ChangePassword
//@description: 修改用户密码
//@param: u *data.SysUser, newPassword string
//@return: userInter *data.SysUser,err error

func (b *User) ChangePassword(u *data.SysUser, newPassword string) (userInter *data.SysUser, err error) {
	var user data.SysUser
	if err = b.data.SqlClient.Where("id = ?", u.ID).First(&user).Error; err != nil {
		return nil, err
	}
	if ok := hash.BcryptCheck(u.Password, user.Password); !ok {
		return nil, errors.New("原密码错误")
	}
	user.Password = hash.BcryptHash(newPassword)
	err = b.data.SqlClient.Save(&user).Error
	return &user, err

}

//@function: GetUserInfoList
//@description: 分页获取数据
//@param: info dto.PageInfo
//@return: err error, list interface{}, total int64

func (b *User) GetUserInfoList(info rhttp.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := b.data.SqlClient.Model(&data.SysUser{})
	var userList []data.SysUser
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

func (b *User) SetUserAuthority(id uint, authorityId uint) (err error) {
	assignErr := b.data.SqlClient.Where("sys_user_id = ? AND sys_authority_authority_id = ?", id, authorityId).First(&data.SysUserAuthority{}).Error
	if errors.Is(assignErr, gorm.ErrRecordNotFound) {
		return errors.New("该用户无此角色")
	}
	err = b.data.SqlClient.Model(&data.SysUser{}).Where("id = ?", id).Update("authority_id", authorityId).Error
	return err
}

//@function: SetUserAuthorities
//@description: 设置一个用户的权限
//@param: id uint, authorityIds []string
//@return: err error

func (b *User) SetUserAuthorities(id uint, authorityIds []uint) (err error) {
	return b.data.SqlClient.Transaction(func(tx *gorm.DB) error {
		var user data.SysUser
		TxErr := tx.Where("id = ?", id).First(&user).Error
		if TxErr != nil {
			b.logger.Error("[SetUserAuthorities]查询用户数据失败", zap.Error(TxErr))
			return errors.New("查询用户数据失败")
		}
		TxErr = tx.Delete(&[]data.SysUserAuthority{}, "sys_user_id = ?", id).Error
		if TxErr != nil {
			return TxErr
		}
		var useAuthority []data.SysUserAuthority
		for _, v := range authorityIds {
			useAuthority = append(useAuthority, data.SysUserAuthority{
				SysUserId: id, SysAuthorityId: v,
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

func (b *User) DeleteUser(id int) (err error) {
	return b.data.SqlClient.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&data.SysUser{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&[]data.SysUserAuthority{}, "sys_user_id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}

//@function: SetUserInfo
//@description: 设置用户信息
//@param: reqUser data.SysUser
//@return: err error, user data.SysUser

func (b *User) SetUserInfo(req data.SysUser) error {
	return b.data.SqlClient.Model(&data.SysUser{}).
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

func (b *User) SetSelfInfo(req data.SysUser) error {
	return b.data.SqlClient.Model(&data.SysUser{}).
		Where("id=?", req.ID).
		Updates(req).Error
}

//@function: GetUserInfo
//@description: 获取用户信息
//@param: uuid uuid.UUID
//@return: err error, user data.SysUser

func (b *User) GetUserInfo(uuid uuid.UUID) (user data.SysUser, err error) {
	var reqUser data.SysUser
	err = b.data.SqlClient.Preload("Authorities").Preload("Authority").First(&reqUser, "uuid = ?", uuid).Error
	if err != nil {
		return reqUser, err
	}
	return reqUser, err
}

//@function: FindUserById
//@description: 通过id获取用户信息
//@param: id int
//@return: err error, user *data.SysUser

func (b *User) FindUserById(id int) (user *data.SysUser, err error) {
	var u data.SysUser
	err = b.data.SqlClient.Where("id = ?", id).First(&u).Error
	return &u, err
}

//@function: FindUserByUuid
//@description: 通过uuid获取用户信息
//@param: uuid string
//@return: err error, user *data.SysUser

func (b *User) FindUserByUuid(uuid string) (user *data.SysUser, err error) {
	var u data.SysUser
	if err = b.data.SqlClient.Where("uuid = ?", uuid).First(&u).Error; err != nil {
		return &u, errors.New("用户不存在")
	}
	return &u, nil
}

//@function: ResetPassword
//@description: 修改用户密码
//@param: ID uint
//@return: err error

func (b *User) ResetPassword(ID uint) (err error) {
	err = b.data.SqlClient.Model(&data.SysUser{}).Where("id = ?", ID).Update("password", utils.BcryptHash("123456")).Error
	return err
}
