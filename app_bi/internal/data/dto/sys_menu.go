package dto

import "github.com/leaf-rain/raindata/app_bi/internal/data"

// Add menu authority info structure
type AddMenuAuthorityInfo struct {
	Menus       []data.SysBaseMenu `json:"menus"`
	AuthorityId uint               `json:"authorityId"` // 角色ID
}

func DefaultMenu() []data.SysBaseMenu {
	return []data.SysBaseMenu{{
		GVA_MODEL: data.GVA_MODEL{ID: 1},
		ParentId:  0,
		Path:      "dashboard",
		Name:      "dashboard",
		Component: "view/dashboard/index.vue",
		Sort:      1,
		Meta: data.Meta{
			Title: "仪表盘",
			Icon:  "setting",
		},
	}}
}

type SysMenusResponse struct {
	Menus []data.SysMenu `json:"menus"`
}

type SysBaseMenusResponse struct {
	Menus []data.SysBaseMenu `json:"menus"`
}

type SysBaseMenuResponse struct {
	Menu data.SysBaseMenu `json:"menu"`
}
