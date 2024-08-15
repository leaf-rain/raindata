package dto

import (
	"github.com/leaf-rain/raindata/app_bi/internal/data"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
)

type SysDictionaryDetailSearch struct {
	data.SysDictionaryDetail
	rhttp.PageInfo
}
