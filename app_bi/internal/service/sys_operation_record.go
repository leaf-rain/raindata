package service

import (
	"github.com/leaf-rain/raindata/app_bi/internal/biz"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"github.com/leaf-rain/raindata/app_bi/internal/data/dto"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
)

//@function: CreateSysOperationRecord
//@description: 创建记录
//@param: sysOperationRecord data.SysOperationRecord
//@return: err error

type OperationRecordService struct {
	*Service
}

func NewOperationRecordService(service *Service) *OperationRecordService {
	return &OperationRecordService{
		service,
	}
}

func (svc *OperationRecordService) CreateSysOperationRecord(sysOperationRecord data.SysOperationRecord) (err error) {
	b := biz.NewOperationRecord(svc.biz)
	return b.CreateSysOperationRecord(sysOperationRecord)
}

//@function: DeleteSysOperationRecordByIds
//@description: 批量删除记录
//@param: ids dto.IdsReq
//@return: err error

func (svc *OperationRecordService) DeleteSysOperationRecordByIds(ids rhttp.IdsReq) (err error) {
	b := biz.NewOperationRecord(svc.biz)
	return b.DeleteSysOperationRecordByIds(ids)
}

//@function: DeleteSysOperationRecord
//@description: 删除操作记录
//@param: sysOperationRecord data.SysOperationRecord
//@return: err error

func (svc *OperationRecordService) DeleteSysOperationRecord(sysOperationRecord data.SysOperationRecord) (err error) {
	b := biz.NewOperationRecord(svc.biz)
	return b.DeleteSysOperationRecord(sysOperationRecord)
}

//@function: GetSysOperationRecord
//@description: 根据id获取单条操作记录
//@param: id uint
//@return: sysOperationRecord data.SysOperationRecord, err error

func (svc *OperationRecordService) GetSysOperationRecord(id uint) (sysOperationRecord data.SysOperationRecord, err error) {
	b := biz.NewOperationRecord(svc.biz)
	return b.GetSysOperationRecord(id)
}

//@function: GetSysOperationRecordInfoList
//@description: 分页获取操作记录列表
//@param: info systemReq.SysOperationRecordSearch
//@return: list interface{}, total int64, err error

func (svc *OperationRecordService) GetSysOperationRecordInfoList(info dto.SysOperationRecordSearch) (list interface{}, total int64, err error) {
	b := biz.NewOperationRecord(svc.biz)
	return b.GetSysOperationRecordInfoList(info)
}
