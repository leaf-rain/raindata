package service

import (
	"errors"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data/dto"
	"github.com/leaf-rain/raindata/app_bi/internal/data/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthorityBtnService struct {
	data *entity.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

var AuthorityBtnServiceApp = new(AuthorityBtnService)

func (svc *AuthorityBtnService) GetAuthorityBtn(req dto.SysAuthorityBtnReq) (res dto.SysAuthorityBtnRes, err error) {
	var authorityBtn []entity.SysAuthorityBtn
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
		var authorityBtn []entity.SysAuthorityBtn
		err = tx.Delete(&[]entity.SysAuthorityBtn{}, "authority_id = ? and sys_menu_id = ?", req.AuthorityId, req.MenuID).Error
		if err != nil {
			return err
		}
		for _, v := range req.Selected {
			authorityBtn = append(authorityBtn, entity.SysAuthorityBtn{
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
	fErr := svc.data.SqlClient.First(&entity.SysAuthorityBtn{}, "sys_base_menu_btn_id = ?", ID).Error
	if errors.Is(fErr, gorm.ErrRecordNotFound) {
		return nil
	}
	return errors.New("此按钮正在被使用无法删除")
}
