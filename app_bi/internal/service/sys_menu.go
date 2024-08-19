package service

import (
	"errors"
	"github.com/leaf-rain/raindata/app_bi/internal/data/entity"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"

	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//@function: getMenuTreeMap
//@description: 获取路由总树map
//@param: authorityId string
//@return: treeMap map[string][]data.SysMenu, err error

type MenuService struct {
	data *entity.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

var MenuServiceApp = new(MenuService)

func (svc *MenuService) getMenuTreeMap(authorityId uint) (treeMap map[uint][]entity.SysMenu, err error) {
	var allMenus []entity.SysMenu
	var baseMenu []entity.SysBaseMenu
	var btns []entity.SysAuthorityBtn
	treeMap = make(map[uint][]entity.SysMenu)

	var SysAuthorityMenus []entity.SysAuthorityMenu
	err = svc.data.SqlClient.Where("sys_authority_authority_id = ?", authorityId).Find(&SysAuthorityMenus).Error
	if err != nil {
		return
	}

	var MenuIds []string

	for i := range SysAuthorityMenus {
		MenuIds = append(MenuIds, SysAuthorityMenus[i].MenuId)
	}

	err = svc.data.SqlClient.Where("id in (?)", MenuIds).Order("sort").Preload("Parameters").Find(&baseMenu).Error
	if err != nil {
		return
	}

	for i := range baseMenu {
		allMenus = append(allMenus, entity.SysMenu{
			SysBaseMenu: baseMenu[i],
			AuthorityId: authorityId,
			MenuId:      baseMenu[i].ID,
			Parameters:  baseMenu[i].Parameters,
		})
	}

	err = svc.data.SqlClient.Where("authority_id = ?", authorityId).Preload("SysBaseMenuBtn").Find(&btns).Error
	if err != nil {
		return
	}
	var btnMap = make(map[uint]map[string]uint)
	for _, v := range btns {
		if btnMap[v.SysMenuID] == nil {
			btnMap[v.SysMenuID] = make(map[string]uint)
		}
		btnMap[v.SysMenuID][v.SysBaseMenuBtn.Name] = authorityId
	}
	for _, v := range allMenus {
		v.Btns = btnMap[v.SysBaseMenu.ID]
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}
	return treeMap, err
}

//@function: GetMenuTree
//@description: 获取动态菜单树
//@param: authorityId string
//@return: menus []data.SysMenu, err error

func (svc *MenuService) GetMenuTree(authorityId uint) (menus []entity.SysMenu, err error) {
	menuTree, err := svc.getMenuTreeMap(authorityId)
	menus = menuTree[0]
	for i := 0; i < len(menus); i++ {
		err = svc.getChildrenList(&menus[i], menuTree)
	}
	return menus, err
}

//@function: getChildrenList
//@description: 获取子菜单
//@param: menu *data.SysMenu, treeMap map[string][]data.SysMenu
//@return: err error

func (svc *MenuService) getChildrenList(menu *entity.SysMenu, treeMap map[uint][]entity.SysMenu) (err error) {
	menu.Children = treeMap[menu.MenuId]
	for i := 0; i < len(menu.Children); i++ {
		err = svc.getChildrenList(&menu.Children[i], treeMap)
	}
	return err
}

//@function: GetInfoList
//@description: 获取路由分页
//@return: list interface{}, total int64,err error

func (svc *MenuService) GetInfoList() (list interface{}, total int64, err error) {
	var menuList []entity.SysBaseMenu
	treeMap, err := svc.getBaseMenuTreeMap()
	menuList = treeMap[0]
	for i := 0; i < len(menuList); i++ {
		err = svc.getBaseChildrenList(&menuList[i], treeMap)
	}
	return menuList, total, err
}

//@function: getBaseChildrenList
//@description: 获取菜单的子菜单
//@param: menu *data.SysBaseMenu, treeMap map[string][]data.SysBaseMenu
//@return: err error

func (svc *MenuService) getBaseChildrenList(menu *entity.SysBaseMenu, treeMap map[uint][]entity.SysBaseMenu) (err error) {
	menu.Children = treeMap[menu.ID]
	for i := 0; i < len(menu.Children); i++ {
		err = svc.getBaseChildrenList(&menu.Children[i], treeMap)
	}
	return err
}

//@function: AddBaseMenu
//@description: 添加基础路由
//@param: menu data.SysBaseMenu
//@return: error

func (svc *MenuService) AddBaseMenu(menu entity.SysBaseMenu) error {
	if !errors.Is(svc.data.SqlClient.Where("name = ?", menu.Name).First(&entity.SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复name，请修改name")
	}
	return svc.data.SqlClient.Create(&menu).Error
}

//@function: getBaseMenuTreeMap
//@description: 获取路由总树map
//@return: treeMap map[string][]data.SysBaseMenu, err error

func (svc *MenuService) getBaseMenuTreeMap() (treeMap map[uint][]entity.SysBaseMenu, err error) {
	var allMenus []entity.SysBaseMenu
	treeMap = make(map[uint][]entity.SysBaseMenu)
	err = svc.data.SqlClient.Order("sort").Preload("MenuBtn").Preload("Parameters").Find(&allMenus).Error
	for _, v := range allMenus {
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}
	return treeMap, err
}

//@function: GetBaseMenuTree
//@description: 获取基础路由树
//@return: menus []data.SysBaseMenu, err error

func (svc *MenuService) GetBaseMenuTree() (menus []entity.SysBaseMenu, err error) {
	treeMap, err := svc.getBaseMenuTreeMap()
	menus = treeMap[0]
	for i := 0; i < len(menus); i++ {
		err = svc.getBaseChildrenList(&menus[i], treeMap)
	}
	return menus, err
}

//@function: AddMenuAuthority
//@description: 为角色增加menu树
//@param: menus []data.SysBaseMenu, authorityId string
//@return: err error

func (svc *MenuService) AddMenuAuthority(menus []entity.SysBaseMenu, authorityId uint) (err error) {
	var auth entity.SysAuthority
	auth.AuthorityId = authorityId
	auth.SysBaseMenus = menus
	err = AuthorityServiceApp.SetMenuAuthority(&auth)
	return err
}

//@function: GetMenuAuthority
//@description: 查看当前角色树
//@param: info *dto.GetAuthorityId
//@return: menus []data.SysMenu, err error

func (svc *MenuService) GetMenuAuthority(info *rhttp.GetAuthorityId) (menus []entity.SysMenu, err error) {
	var baseMenu []entity.SysBaseMenu
	var SysAuthorityMenus []entity.SysAuthorityMenu
	err = svc.data.SqlClient.Where("sys_authority_authority_id = ?", info.AuthorityId).Find(&SysAuthorityMenus).Error
	if err != nil {
		return
	}

	var MenuIds []string

	for i := range SysAuthorityMenus {
		MenuIds = append(MenuIds, SysAuthorityMenus[i].MenuId)
	}

	err = svc.data.SqlClient.Where("id in (?) ", MenuIds).Order("sort").Find(&baseMenu).Error

	for i := range baseMenu {
		menus = append(menus, entity.SysMenu{
			SysBaseMenu: baseMenu[i],
			AuthorityId: info.AuthorityId,
			MenuId:      baseMenu[i].ID,
			Parameters:  baseMenu[i].Parameters,
		})
	}
	return menus, err
}

// UserAuthorityDefaultRouter 用户角色默认路由检查
func (svc *MenuService) UserAuthorityDefaultRouter(user *entity.SysUser) {
	var menuIds []string
	err := svc.data.SqlClient.Model(&entity.SysAuthorityMenu{}).Where("sys_authority_authority_id = ?", user.AuthorityId).Pluck("sys_base_menu_id", &menuIds).Error
	if err != nil {
		return
	}
	var am entity.SysBaseMenu
	err = svc.data.SqlClient.First(&am, "name = ? and id in (?)", user.Authority.DefaultRouter, menuIds).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user.Authority.DefaultRouter = "404"
	}
}
