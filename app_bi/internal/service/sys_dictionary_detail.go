package service

import (
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data"
	"github.com/leaf-rain/raindata/app_bi/internal/data/dto"
	"go.uber.org/zap"
)

//@function: CreateSysDictionaryDetail
//@description: 创建字典详情数据
//@param: sysDictionaryDetail data.SysDictionaryDetail
//@return: err error

type DictionaryDetailService struct {
	data *data.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

var DictionaryDetailServiceApp = new(DictionaryDetailService)

func (svc *DictionaryDetailService) CreateSysDictionaryDetail(sysDictionaryDetail data.SysDictionaryDetail) (err error) {
	err = svc.data.SqlClient.Create(&sysDictionaryDetail).Error
	return err
}

//@function: DeleteSysDictionaryDetail
//@description: 删除字典详情数据
//@param: sysDictionaryDetail data.SysDictionaryDetail
//@return: err error

func (svc *DictionaryDetailService) DeleteSysDictionaryDetail(sysDictionaryDetail data.SysDictionaryDetail) (err error) {
	err = svc.data.SqlClient.Delete(&sysDictionaryDetail).Error
	return err
}

//@function: UpdateSysDictionaryDetail
//@description: 更新字典详情数据
//@param: sysDictionaryDetail *data.SysDictionaryDetail
//@return: err error

func (svc *DictionaryDetailService) UpdateSysDictionaryDetail(sysDictionaryDetail *data.SysDictionaryDetail) (err error) {
	err = svc.data.SqlClient.Save(sysDictionaryDetail).Error
	return err
}

//@function: GetSysDictionaryDetail
//@description: 根据id获取字典详情单条数据
//@param: id uint
//@return: sysDictionaryDetail data.SysDictionaryDetail, err error

func (svc *DictionaryDetailService) GetSysDictionaryDetail(id uint) (sysDictionaryDetail data.SysDictionaryDetail, err error) {
	err = svc.data.SqlClient.Where("id = ?", id).First(&sysDictionaryDetail).Error
	return
}

//@function: GetSysDictionaryDetailInfoList
//@description: 分页获取字典详情列表
//@param: info dto.SysDictionaryDetailSearch
//@return: list interface{}, total int64, err error

func (svc *DictionaryDetailService) GetSysDictionaryDetailInfoList(info dto.SysDictionaryDetailSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := svc.data.SqlClient.Model(&data.SysDictionaryDetail{})
	var sysDictionaryDetails []data.SysDictionaryDetail
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Label != "" {
		db = db.Where("label LIKE ?", "%"+info.Label+"%")
	}
	if info.Value != "" {
		db = db.Where("value = ?", info.Value)
	}
	if info.Status != nil {
		db = db.Where("status = ?", info.Status)
	}
	if info.SysDictionaryID != 0 {
		db = db.Where("sys_dictionary_id = ?", info.SysDictionaryID)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("sort").Find(&sysDictionaryDetails).Error
	return sysDictionaryDetails, total, err
}

// 按照字典id获取字典全部内容的方法
func (svc *DictionaryDetailService) GetDictionaryList(dictionaryID uint) (list []data.SysDictionaryDetail, err error) {
	var sysDictionaryDetails []data.SysDictionaryDetail
	err = svc.data.SqlClient.Find(&sysDictionaryDetails, "sys_dictionary_id = ?", dictionaryID).Error
	return sysDictionaryDetails, err
}

// 按照字典type获取字典全部内容的方法
func (svc *DictionaryDetailService) GetDictionaryListByType(t string) (list []data.SysDictionaryDetail, err error) {
	var sysDictionaryDetails []data.SysDictionaryDetail
	db := svc.data.SqlClient.Model(&data.SysDictionaryDetail{}).Joins("JOIN sys_dictionaries ON sys_dictionaries.id = sys_dictionary_details.sys_dictionary_id")
	err = db.Debug().Find(&sysDictionaryDetails, "type = ?", t).Error
	return sysDictionaryDetails, err
}

// 按照字典id+字典内容value获取单条字典内容
func (svc *DictionaryDetailService) GetDictionaryInfoByValue(dictionaryID uint, value string) (detail data.SysDictionaryDetail, err error) {
	var sysDictionaryDetail data.SysDictionaryDetail
	err = svc.data.SqlClient.First(&sysDictionaryDetail, "sys_dictionary_id = ? and value = ?", dictionaryID, value).Error
	return sysDictionaryDetail, err
}

// 按照字典type+字典内容value获取单条字典内容
func (svc *DictionaryDetailService) GetDictionaryInfoByTypeValue(t string, value string) (detail data.SysDictionaryDetail, err error) {
	var sysDictionaryDetails data.SysDictionaryDetail
	db := svc.data.SqlClient.Model(&data.SysDictionaryDetail{}).Joins("JOIN sys_dictionaries ON sys_dictionaries.id = sys_dictionary_details.sys_dictionary_id")
	err = db.First(&sysDictionaryDetails, "sys_dictionaries.type = ? and sys_dictionary_details.value = ?", t, value).Error
	return sysDictionaryDetails, err
}
