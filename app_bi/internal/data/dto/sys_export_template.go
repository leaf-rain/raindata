package dto

import (
	"github.com/leaf-rain/raindata/app_bi/internal/data"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"time"
)

type SysExportTemplateSearch struct {
	data.SysExportTemplate
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	rhttp.PageInfo
}
