package service

import (
	"errors"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data"
	"github.com/leaf-rain/raindata/app_bi/internal/data/dto"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthorityBtnService struct {
	data *data.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

var AuthorityBtnServiceApp = new(AuthorityBtnService)

func (svc *AuthorityBtnService) GetAuthorityBtn(req dto.SysAuthorityBtnReq) (res dto.SysAuthorityBtnRes, err error) {
	var authorityBtn []data.SysAuthorityBtn
	err = svc.data.SqlClient.Find(&authorityBtn, "authority_id = ? and sys_menu_id = ?", req.AuthorityId, req.MenuID).Error
	if err != nil {
		return
	}
	var selected []uint
	for _, v := range authorityBtn {
		selected = append(selected, v.SysBaseMenuBtnID)
	}
	res.Selected = selected
	return res, err
}

func (svc *AuthorityBtnService) SetAuthorityBtn(req dto.SysAuthorityBtnReq) (err error) {
	return svc.data.SqlClient.Transaction(func(tx *gorm.DB) error {
		var authorityBtn []data.SysAuthorityBtn
		err = tx.Delete(&[]data.SysAuthorityBtn{}, "authority_id = ? and sys_menu_id = ?", req.AuthorityId, req.MenuID).Error
		if err != nil {
			return err
		}
		for _, v := range req.Selected {
			authorityBtn = append(authorityBtn, data.SysAuthorityBtn{
				AuthorityId:      req.AuthorityId,
				SysMenuID:        req.MenuID,
				SysBaseMenuBtnID: v,
			})
		}
		if len(authorityBtn) > 0 {
			err = tx.Create(&authorityBtn).Error
		}
		if err != nil {
			return err
		}
		return err
	})
}

func (svc *AuthorityBtnService) CanRemoveAuthorityBtn(ID string) (err error) {
	fErr := svc.data.SqlClient.First(&data.SysAuthorityBtn{}, "sys_base_menu_btn_id = ?", ID).Error
	if errors.Is(fErr, gorm.ErrRecordNotFound) {
		return nil
	}
	return errors.New("此按钮正在被使用无法删除")
}
