package entity

type SysUserResponse struct {
	User SysUser `json:"user"`
}

type LoginResponse struct {
	User      SysUser `json:"user"`
	Token     string  `json:"token"`
	ExpiresAt int64   `json:"expiresAt"`
}
