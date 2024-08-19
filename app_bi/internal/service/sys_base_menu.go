package service

import (
	"errors"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BaseMenuService struct {
	data *entity.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

//@function: DeleteBaseMenu
//@description: 删除基础路由
//@param: id float64
//@return: err error

var BaseMenuServiceApp = new(BaseMenuService)

func (svc *BaseMenuService) DeleteBaseMenu(id int) (err error) {
	err = svc.data.SqlClient.First(&entity.SysBaseMenu{}, "parent_id = ?", id).Error
	if err != nil {
		return svc.data.SqlClient.Transaction(func(tx *gorm.DB) error {

			err = tx.Delete(&entity.SysBaseMenu{}, "id = ?", id).Error
			if err != nil {
				return err
			}

			err = tx.Delete(&entity.SysBaseMenuParameter{}, "sys_base_menu_id = ?", id).Error
			if err != nil {
				return err
			}

			err = tx.Delete(&entity.SysBaseMenuBtn{}, "sys_base_menu_id = ?", id).Error
			if err != nil {
				return err
			}
			err = tx.Delete(&entity.SysAuthorityBtn{}, "sys_menu_id = ?", id).Error
			if err != nil {
				return err
			}

			err = tx.Delete(&entity.SysAuthorityMenu{}, "sys_base_menu_id = ?", id).Error
			if err != nil {
				return err
			}
			return nil
		})
	}
	return errors.New("此菜单存在子菜单不可删除")
}

//@function: UpdateBaseMenu
//@description: 更新路由
//@param: menu data.SysBaseMenu
//@return: err error

func (svc *BaseMenuService) UpdateBaseMenu(menu entity.SysBaseMenu) (err error) {
	var oldMenu entity.SysBaseMenu
	upDateMap := make(map[string]interface{})
	upDateMap["keep_alive"] = menu.KeepAlive
	upDateMap["close_tab"] = menu.CloseTab
	upDateMap["default_menu"] = menu.DefaultMenu
	upDateMap["parent_id"] = menu.ParentId
	upDateMap["path"] = menu.Path
	upDateMap["name"] = menu.Name
	upDateMap["hidden"] = menu.Hidden
	upDateMap["component"] = menu.Component
	upDateMap["title"] = menu.Title
	upDateMap["active_name"] = menu.ActiveName
	upDateMap["icon"] = menu.Icon
	upDateMap["sort"] = menu.Sort

	err = svc.data.SqlClient.Transaction(func(tx *gorm.DB) error {
		tx.Where("id = ?", menu.ID).Find(&oldMenu)
		if oldMenu.Name != menu.Name {
			if !errors.Is(tx.Where("id <> ? AND name = ?", menu.ID, menu.Name).First(&entity.SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
				svc.log.Debug("存在相同name修改失败")
				return errors.New("存在相同name修改失败")
			}
		}
		txErr := tx.Unscoped().Delete(&entity.SysBaseMenuParameter{}, "sys_base_menu_id = ?", menu.ID).Error
		if txErr != nil {
			svc.log.Debug(txErr.Error())
			return txErr
		}
		txErr = tx.Unscoped().Delete(&entity.SysBaseMenuBtn{}, "sys_base_menu_id = ?", menu.ID).Error
		if txErr != nil {
			svc.log.Debug(txErr.Error())
			return txErr
		}
		if len(menu.Parameters) > 0 {
			for k := range menu.Parameters {
				menu.Parameters[k].SysBaseMenuID = menu.ID
			}
			txErr = tx.Create(&menu.Parameters).Error
			if txErr != nil {
				svc.log.Debug(txErr.Error())
				return txErr
			}
		}

		if len(menu.MenuBtn) > 0 {
			for k := range menu.MenuBtn {
				menu.MenuBtn[k].SysBaseMenuID = menu.ID
			}
			txErr = tx.Create(&menu.MenuBtn).Error
			if txErr != nil {
				svc.log.Debug(txErr.Error())
				return txErr
			}
		}

		txErr = tx.Model(&oldMenu).Updates(upDateMap).Error
		if txErr != nil {
			svc.log.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return err
}

//@function: GetBaseMenuById
//@description: 返回当前选中menu
//@param: id float64
//@return: menu data.SysBaseMenu, err error

func (svc *BaseMenuService) GetBaseMenuById(id int) (menu entity.SysBaseMenu, err error) {
	err = svc.data.SqlClient.Preload("MenuBtn").Preload("Parameters").Where("id = ?", id).First(&menu).Error
	return
}
