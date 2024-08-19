package dto

import (
	"github.com/leaf-rain/raindata/app_bi/internal/data/entity"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
)

type SysDictionaryDetailSearch struct {
	entity.SysDictionaryDetail
	rhttp.PageInfo
}
