package dto

import (
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
)

type SysOperationRecordSearch struct {
	data.SysOperationRecord
	rhttp.PageInfo
}
