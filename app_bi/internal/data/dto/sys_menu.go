package dto

import (
	"github.com/leaf-rain/raindata/app_bi/internal/data/entity"
)

// Add menu authority info structure
type AddMenuAuthorityInfo struct {
	Menus       []entity.SysBaseMenu `json:"menus"`
	AuthorityId uint                 `json:"authorityId"` // 角色ID
}

func DefaultMenu() []entity.SysBaseMenu {
	return []entity.SysBaseMenu{{
		gorm.Model: entity.gorm.Model{ID: 1},
		ParentId:  0,
		Path:      "dashboard",
		Name:      "dashboard",
		Component: "view/dashboard/index.vue",
		Sort:      1,
		Meta: entity.Meta{
			Title: "仪表盘",
			Icon:  "setting",
		},
	}}
}

type SysMenusResponse struct {
	Menus []entity.SysMenu `json:"menus"`
}

type SysBaseMenusResponse struct {
	Menus []entity.SysBaseMenu `json:"menus"`
}

type SysBaseMenuResponse struct {
	Menu entity.SysBaseMenu `json:"menu"`
}
