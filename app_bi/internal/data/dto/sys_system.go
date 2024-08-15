package dto

import "github.com/leaf-rain/raindata/app_bi/internal/conf"

type SysConfigResponse struct {
	Config conf.Bootstrap `json:"config"`
}
