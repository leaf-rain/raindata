package service

import (
	"github.com/leaf-rain/raindata/app_bi/internal/data"
	"github.com/leaf-rain/raindata/app_bi/internal/data/entity"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
)

type CustomerService struct {
	data *entity.Data
}

var CustomerServiceApp = new(CustomerService)

//@function: CreateExaCustomer
//@description: 创建客户
//@param: e data.ExaCustomer
//@return: err error

func (exa *CustomerService) CreateExaCustomer(e data.ExaCustomer) (err error) {
	err = exa.data.SqlClient.Create(&e).Error
	return err
}

//@function: DeleteFileChunk
//@description: 删除客户
//@param: e data.ExaCustomer
//@return: err error

func (exa *CustomerService) DeleteExaCustomer(e data.ExaCustomer) (err error) {
	err = exa.data.SqlClient.Delete(&e).Error
	return err
}

//@function: UpdateExaCustomer
//@description: 更新客户
//@param: e *data.ExaCustomer
//@return: err error

func (exa *CustomerService) UpdateExaCustomer(e *data.ExaCustomer) (err error) {
	err = exa.data.SqlClient.Save(e).Error
	return err
}

//@function: GetExaCustomer
//@description: 获取客户信息
//@param: id uint
//@return: customer data.ExaCustomer, err error

func (exa *CustomerService) GetExaCustomer(id uint) (customer data.ExaCustomer, err error) {
	err = exa.data.SqlClient.Where("id = ?", id).First(&customer).Error
	return
}

//@function: GetCustomerInfoList
//@description: 分页获取客户列表
//@param: sysUserAuthorityID string, info dto.PageInfo
//@return: list interface{}, total int64, err error

func (exa *CustomerService) GetCustomerInfoList(sysUserAuthorityID uint, info rhttp.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := exa.data.SqlClient.Model(&data.ExaCustomer{})
	var a entity.SysAuthority
	a.AuthorityId = sysUserAuthorityID
	auth, err := systemService.AuthorityServiceApp.GetAuthorityInfo(a)
	if err != nil {
		return
	}
	var dataId []uint
	for _, v := range auth.DataAuthorityId {
		dataId = append(dataId, v.AuthorityId)
	}
	var CustomerList []data.ExaCustomer
	err = db.Where("sys_user_authority_id in ?", dataId).Count(&total).Error
	if err != nil {
		return CustomerList, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Preload("SysUser").Where("sys_user_authority_id in ?", dataId).Find(&CustomerList).Error
	}
	return CustomerList, total, err
}
