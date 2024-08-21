package service

import (
	"github.com/google/uuid"
	"github.com/leaf-rain/raindata/app_bi/internal/biz"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
)

type UserService struct {
	*Service
}

// NewUserService new a greeter service.
func NewUserService(service *Service) *UserService {
	return &UserService{
		service,
	}
}

func (svc *UserService) Register(u data.SysUser) (userInter data.SysUser, err error) {
	b := biz.NewUser(svc.biz)
	return b.Register(u)
}

//@function: Login
//@description: 用户登录
//@param: u *data.SysUser
//@return: err error, userInter *data.SysUser

func (svc *UserService) Login(u *data.SysUser) (userInter *data.SysUser, err error) {
	b := biz.NewUser(svc.biz)
	return b.Login(u)
}

//@function: ChangePassword
//@description: 修改用户密码
//@param: u *data.SysUser, newPassword string
//@return: userInter *data.SysUser,err error

func (svc *UserService) ChangePassword(u *data.SysUser, newPassword string) (userInter *data.SysUser, err error) {
	b := biz.NewUser(svc.biz)
	return b.ChangePassword(u, newPassword)
}

//@function: GetUserInfoList
//@description: 分页获取数据
//@param: info dto.PageInfo
//@return: err error, list interface{}, total int64

func (svc *UserService) GetUserInfoList(info rhttp.PageInfo) (list interface{}, total int64, err error) {
	b := biz.NewUser(svc.biz)
	return b.GetUserInfoList(info)
}

//@function: SetUserAuthority
//@description: 设置一个用户的权限
//@param: uuid uuid.UUID, authorityId string
//@return: err error

func (svc *UserService) SetUserAuthority(id uint, authorityId uint) (err error) {
	b := biz.NewUser(svc.biz)
	return b.SetUserAuthority(id, authorityId)
}

//@function: SetUserAuthorities
//@description: 设置一个用户的权限
//@param: id uint, authorityIds []string
//@return: err error

func (svc *UserService) SetUserAuthorities(id uint, authorityIds []uint) (err error) {
	b := biz.NewUser(svc.biz)
	return b.SetUserAuthorities(id, authorityIds)
}

//@function: DeleteUser
//@description: 删除用户
//@param: id float64
//@return: err error

func (svc *UserService) DeleteUser(id int) (err error) {
	b := biz.NewUser(svc.biz)
	return b.DeleteUser(id)
}

//@function: SetUserInfo
//@description: 设置用户信息
//@param: reqUser data.SysUser
//@return: err error, user data.SysUser

func (svc *UserService) SetUserInfo(req data.SysUser) error {
	b := biz.NewUser(svc.biz)
	return b.SetUserInfo(req)
}

//@function: SetSelfInfo
//@description: 设置用户信息
//@param: reqUser data.SysUser
//@return: err error, user data.SysUser

func (svc *UserService) SetSelfInfo(req data.SysUser) error {
	b := biz.NewUser(svc.biz)
	return b.SetSelfInfo(req)
}

//@function: GetUserInfo
//@description: 获取用户信息
//@param: uuid uuid.UUID
//@return: err error, user data.SysUser

func (svc *UserService) GetUserInfo(uuid uuid.UUID) (user data.SysUser, err error) {
	b := biz.NewUser(svc.biz)
	return b.GetUserInfo(uuid)
}

//@function: FindUserById
//@description: 通过id获取用户信息
//@param: id int
//@return: err error, user *data.SysUser

func (svc *UserService) FindUserById(id int) (user *data.SysUser, err error) {
	b := biz.NewUser(svc.biz)
	return b.FindUserById(id)
}

//@function: FindUserByUuid
//@description: 通过uuid获取用户信息
//@param: uuid string
//@return: err error, user *data.SysUser

func (svc *UserService) FindUserByUuid(uuid string) (user *data.SysUser, err error) {
	b := biz.NewUser(svc.biz)
	return b.FindUserByUuid(uuid)
}

//@function: ResetPassword
//@description: 修改用户密码
//@param: ID uint
//@return: err error

func (svc *UserService) ResetPassword(ID uint) (err error) {
	b := biz.NewUser(svc.biz)
	return b.ResetPassword(ID)
}
