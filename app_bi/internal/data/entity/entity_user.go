package entity

import (
	"errors"
	"github.com/leaf-rain/raindata/app_bi/third_party/hash"
)

var (
	ErrUserPassword = errors.New("密码错误")
)

func NewEntityUser(user *SysUser, data *Data) *EntityUser {
	return &EntityUser{user: user, data: data}
}

type EntityUser struct {
	user *SysUser
	data *Data
}

func (entity *EntityUser) ReloadByDb() error {
	var user *SysUser
	err := entity.data.SqlClient.Where("username = ?", entity.user.Username).Preload("Authorities").Preload("Authority").First(&user).Error
	if err == nil {
		if ok := hash.BcryptCheck(entity.user.Password, user.Password); !ok {
			return ErrUserPassword
		}
	}
	entity.user = user
	return err
}
