package service

import (
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data"
	"github.com/leaf-rain/raindata/app_bi/internal/data/dto"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"go.uber.org/zap"
)

//@function: CreateSysOperationRecord
//@description: 创建记录
//@param: sysOperationRecord data.SysOperationRecord
//@return: err error

type OperationRecordService struct {
	data *data.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

var OperationRecordServiceApp = new(OperationRecordService)

func (svc *OperationRecordService) CreateSysOperationRecord(sysOperationRecord data.SysOperationRecord) (err error) {
	err = svc.data.SqlClient.Create(&sysOperationRecord).Error
	return err
}

//@function: DeleteSysOperationRecordByIds
//@description: 批量删除记录
//@param: ids dto.IdsReq
//@return: err error

func (svc *OperationRecordService) DeleteSysOperationRecordByIds(ids rhttp.IdsReq) (err error) {
	err = svc.data.SqlClient.Delete(&[]data.SysOperationRecord{}, "id in (?)", ids.Ids).Error
	return err
}

//@function: DeleteSysOperationRecord
//@description: 删除操作记录
//@param: sysOperationRecord data.SysOperationRecord
//@return: err error

func (svc *OperationRecordService) DeleteSysOperationRecord(sysOperationRecord data.SysOperationRecord) (err error) {
	err = svc.data.SqlClient.Delete(&sysOperationRecord).Error
	return err
}

//@function: GetSysOperationRecord
//@description: 根据id获取单条操作记录
//@param: id uint
//@return: sysOperationRecord data.SysOperationRecord, err error

func (svc *OperationRecordService) GetSysOperationRecord(id uint) (sysOperationRecord data.SysOperationRecord, err error) {
	err = svc.data.SqlClient.Where("id = ?", id).First(&sysOperationRecord).Error
	return
}

//@function: GetSysOperationRecordInfoList
//@description: 分页获取操作记录列表
//@param: info systemReq.SysOperationRecordSearch
//@return: list interface{}, total int64, err error

func (svc *OperationRecordService) GetSysOperationRecordInfoList(info dto.SysOperationRecordSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := svc.data.SqlClient.Model(&data.SysOperationRecord{})
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
