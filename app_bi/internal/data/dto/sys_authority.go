package dto

import "github.com/leaf-rain/raindata/app_bi/internal/data"

type SysAuthorityResponse struct {
	Authority data.SysAuthority `json:"authority"`
}

type SysAuthorityCopyResponse struct {
	Authority      data.SysAuthority `json:"authority"`
	OldAuthorityId uint              `json:"oldAuthorityId"` // 旧角色ID
}
