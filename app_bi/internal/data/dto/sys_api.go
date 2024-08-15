package dto

import (
	"github.com/leaf-rain/raindata/app_bi/internal/data"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
)

// api分页条件查询及排序结构体
type SearchApiParams struct {
	data.SysApi
	rhttp.PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

type SysAPIResponse struct {
	Api data.SysApi `json:"api"`
}

type SysAPIListResponse struct {
	Apis []data.SysApi `json:"apis"`
}

type SysSyncApis struct {
	NewApis    []data.SysApi `json:"newApis"`
	DeleteApis []data.SysApi `json:"deleteApis"`
}
