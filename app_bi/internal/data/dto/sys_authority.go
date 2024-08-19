package dto

import (
	"github.com/leaf-rain/raindata/app_bi/internal/data/entity"
)

type SysAuthorityResponse struct {
	Authority entity.SysAuthority `json:"authority"`
}

type SysAuthorityCopyResponse struct {
	Authority      entity.SysAuthority `json:"authority"`
	OldAuthorityId uint                `json:"oldAuthorityId"` // 旧角色ID
}
