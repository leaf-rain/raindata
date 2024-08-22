package dto

import (
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
)

type SysAuthorityResponse struct {
	Authority data.SysAuthority `json:"authority"`
}
