package dto

import (
	"github.com/leaf-rain/raindata/app_bi/internal/data/entity"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
)

// api分页条件查询及排序结构体
type SearchApiParams struct {
	entity.SysApi
	rhttp.PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

type SysAPIResponse struct {
	Api entity.SysApi `json:"api"`
}

type SysAPIListResponse struct {
	Apis []entity.SysApi `json:"apis"`
}

type SysSyncApis struct {
	NewApis    []entity.SysApi `json:"newApis"`
	DeleteApis []entity.SysApi `json:"deleteApis"`
}
