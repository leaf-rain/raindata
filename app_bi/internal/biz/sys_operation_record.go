package biz

import (
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"github.com/leaf-rain/raindata/app_bi/internal/data/dto"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
)

//@function: CreateSysOperationRecord
//@description: 创建记录
//@param: sysOperationRecord data.SysOperationRecord
//@return: err error

type OperationRecord struct {
	*Business
}

func NewOperationRecord(biz *Business) *OperationRecord {
	return &OperationRecord{
		biz,
	}
}

func (b *OperationRecord) CreateSysOperationRecord(sysOperationRecord data.SysOperationRecord) (err error) {
	err = b.data.SqlClient.Create(&sysOperationRecord).Error
	return err
}

//@function: DeleteSysOperationRecordByIds
//@description: 批量删除记录
//@param: ids rhttp.IdsReq
//@return: err error

func (b *OperationRecord) DeleteSysOperationRecordByIds(ids rhttp.IdsReq) (err error) {
	err = b.data.SqlClient.Delete(&[]data.SysOperationRecord{}, "id in (?)", ids.Ids).Error
	return err
}

//@function: DeleteSysOperationRecord
//@description: 删除操作记录
//@param: sysOperationRecord data.SysOperationRecord
//@return: err error

func (b *OperationRecord) DeleteSysOperationRecord(sysOperationRecord data.SysOperationRecord) (err error) {
	err = b.data.SqlClient.Delete(&sysOperationRecord).Error
	return err
}

//@function: GetSysOperationRecord
//@description: 根据id获取单条操作记录
//@param: id uint
//@return: sysOperationRecord data.SysOperationRecord, err error

func (b *OperationRecord) GetSysOperationRecord(id uint) (sysOperationRecord data.SysOperationRecord, err error) {
	err = b.data.SqlClient.Where("id = ?", id).First(&sysOperationRecord).Error
	return
}

//@function: GetSysOperationRecordInfoList
//@description: 分页获取操作记录列表
//@param: info dto.SysOperationRecordSearch
//@return: list interface{}, total int64, err error

func (b *OperationRecord) GetSysOperationRecordInfoList(info dto.SysOperationRecordSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := b.data.SqlClient.Model(&data.SysOperationRecord{})
	var sysOperationRecords []data.SysOperationRecord
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Method != "" {
		db = db.Where("method = ?", info.Method)
	}
	if info.Path != "" {
		db = db.Where("path LIKE ?", "%"+info.Path+"%")
	}
	if info.Status != 0 {
		db = db.Where("status = ?", info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("id desc").Limit(limit).Offset(offset).Preload("User").Find(&sysOperationRecords).Error
	return sysOperationRecords, total, err
}
