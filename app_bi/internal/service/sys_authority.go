package service

import (
	"errors"
	"github.com/leaf-rain/raindata/app_bi/internal/biz"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
)

var ErrRoleExistence = errors.New("存在相同角色id")

//@function: CreateAuthority
//@description: 创建一个角色
//@param: auth data.SysAuthority
//@return: authority data.SysAuthority, err error

type AuthorityService struct {
	*Service
}

func NewAuthorityService(service *Service) *AuthorityService {
	return &AuthorityService{
		service,
	}
}
func (svc *AuthorityService) CreateAuthority(auth data.SysAuthority) (authority data.SysAuthority, err error) {
	b := biz.NewAuthority(svc.biz)
	return b.CreateAuthority(auth)
}

//@function: UpdateAuthority
//@description: 更改一个角色
//@param: auth data.SysAuthority
//@return: authority data.SysAuthority, err error

func (svc *AuthorityService) UpdateAuthority(auth data.SysAuthority) (authority data.SysAuthority, err error) {
	b := biz.NewAuthority(svc.biz)
	return b.UpdateAuthority(auth)
}

//@function: DeleteAuthority
//@description: 删除角色
//@param: auth *data.SysAuthority
//@return: err error

func (svc *AuthorityService) DeleteAuthority(auth *data.SysAuthority) error {
	b := biz.NewAuthority(svc.biz)
	return b.DeleteAuthority(auth)
}

//@function: GetAuthorityInfoList
//@description: 分页获取数据
//@param: info dto.PageInfo
//@return: list interface{}, total int64, err error

func (svc *AuthorityService) GetAuthorityInfoList(info rhttp.PageInfo) (list interface{}, total int64, err error) {
	b := biz.NewAuthority(svc.biz)
	return b.GetAuthorityInfoList(info)
}

//@function: GetAuthorityInfo
//@description: 获取所有角色信息
//@param: auth data.SysAuthority
//@return: sa data.SysAuthority, err error

func (svc *AuthorityService) GetAuthorityInfo(auth data.SysAuthority) (sa data.SysAuthority, err error) {
	b := biz.NewAuthority(svc.biz)
	return b.GetAuthorityInfo(auth)
}

//@function: SetDataAuthority
//@description: 设置角色资源权限
//@param: auth data.SysAuthority
//@return: error

func (svc *AuthorityService) SetDataAuthority(auth data.SysAuthority) error {
	b := biz.NewAuthority(svc.biz)
	return b.SetDataAuthority(auth)
}
