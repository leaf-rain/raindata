package dto

type SysAuthorityBtnReq struct {
	MenuID      uint   `json:"menuID"`
	AuthorityId uint   `json:"authorityId"`
	Selected    []uint `json:"selected"`
}

type SysAuthorityBtnRes struct {
	Selected []uint `json:"selected"`
}
