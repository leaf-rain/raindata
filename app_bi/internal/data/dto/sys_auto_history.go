package dto

import (
	"github.com/leaf-rain/raindata/app_bi/internal/data/entity"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
)

type SysAutoHistoryCreate struct {
	Table       string            // 表名
	Package     string            // 模块名/插件名
	Request     string            // 前端传入的结构化信息
	StructName  string            // 结构体名称
	BusinessDB  string            // 业务库
	Description string            // Struct中文名称
	Injections  map[string]string // 注入路径
	Templates   map[string]string // 模板信息
	ApiIDs      []uint            // api表注册内容
	MenuID      uint              // 菜单ID
}

func (r *SysAutoHistoryCreate) Create() entity.SysAutoCodeHistory {
	entity := entity.SysAutoCodeHistory{
		Package:     r.Package,
		Request:     r.Request,
		Table:       r.Table,
		StructName:  r.StructName,
		BusinessDB:  r.BusinessDB,
		Description: r.Description,
		Injections:  r.Injections,
		Templates:   r.Templates,
		ApiIDs:      r.ApiIDs,
		MenuID:      r.MenuID,
	}
	if entity.Table == "" {
		entity.Table = r.StructName
	}
	return entity
}

type SysAutoHistoryRollBack struct {
	rhttp.GetById
	DeleteApi   bool `json:"deleteApi" form:"deleteApi"`     // 是否删除接口
	DeleteMenu  bool `json:"deleteMenu" form:"deleteMenu"`   // 是否删除菜单
	DeleteTable bool `json:"deleteTable" form:"deleteTable"` // 是否删除表
}

func (r *SysAutoHistoryRollBack) ApiIds(entity entity.SysAutoCodeHistory) rhttp.IdsReq {
	length := len(entity.ApiIDs)
	ids := make([]int, 0)
	for i := 0; i < length; i++ {
		ids = append(ids, int(entity.ApiIDs[i]))
	}
	return rhttp.IdsReq{Ids: ids}
}
