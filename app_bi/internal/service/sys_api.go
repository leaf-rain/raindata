package service

import (
	"errors"
	"fmt"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data"
	"github.com/leaf-rain/raindata/app_bi/internal/data/dto"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

//@function: CreateApi
//@description: 新增基础api
//@param: api data.SysApi
//@return: err error

type ApiService struct {
	data *data.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

var ApiServiceApp = new(ApiService)

func (svc *ApiService) CreateApi(api data.SysApi) (err error) {
	if !errors.Is(svc.data.SqlClient.Where("path = ? AND method = ?", api.Path, api.Method).First(&data.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同api")
	}
	return svc.data.SqlClient.Create(&api).Error
}

func (svc *ApiService) GetApiGroups() (groups []string, groupApiMap map[string]string, err error) {
	var apis []data.SysApi
	err = svc.data.SqlClient.Find(&apis).Error
	if err != nil {
		return
	}
	groupApiMap = make(map[string]string, 0)
	for i := range apis {
		pathArr := strings.Split(apis[i].Path, "/")
		newGroup := true
		for i2 := range groups {
			if groups[i2] == apis[i].ApiGroup {
				newGroup = false
			}
		}
		if newGroup {
			groups = append(groups, apis[i].ApiGroup)
		}
		groupApiMap[pathArr[1]] = apis[i].ApiGroup
	}
	return
}

func (svc *ApiService) SyncApi() (newApis, deleteApis, ignoreApis []data.SysApi, err error) {
	newApis = make([]data.SysApi, 0)
	deleteApis = make([]data.SysApi, 0)
	ignoreApis = make([]data.SysApi, 0)
	var apis []data.SysApi
	err = svc.data.SqlClient.Find(&apis).Error
	if err != nil {
		return
	}
	var ignores []data.SysIgnoreApi
	err = svc.data.SqlClient.Find(&ignores).Error
	if err != nil {
		return
	}

	for i := range ignores {
		ignoreApis = append(ignoreApis, data.SysApi{
			Path:        ignores[i].Path,
			Description: "",
			ApiGroup:    "",
			Method:      ignores[i].Method,
		})
	}

	var cacheApis []data.SysApi
	for i := range global.GVA_ROUTERS {
		ignoresFlag := false
		for j := range ignores {
			if ignores[j].Path == global.GVA_ROUTERS[i].Path && ignores[j].Method == global.GVA_ROUTERS[i].Method {
				ignoresFlag = true
			}
		}
		if !ignoresFlag {
			cacheApis = append(cacheApis, data.SysApi{
				Path:   global.GVA_ROUTERS[i].Path,
				Method: global.GVA_ROUTERS[i].Method,
			})
		}
	}

	//对比数据库中的api和内存中的api，如果数据库中的api不存在于内存中，则把api放入删除数组，如果内存中的api不存在于数据库中，则把api放入新增数组
	for i := range cacheApis {
		var flag bool
		// 如果存在于内存不存在于api数组中
		for j := range apis {
			if cacheApis[i].Path == apis[j].Path && cacheApis[i].Method == apis[j].Method {
				flag = true
			}
		}
		if !flag {
			newApis = append(newApis, data.SysApi{
				Path:        cacheApis[i].Path,
				Description: "",
				ApiGroup:    "",
				Method:      cacheApis[i].Method,
			})
		}
	}

	for i := range apis {
		var flag bool
		// 如果存在于api数组不存在于内存
		for j := range cacheApis {
			if cacheApis[j].Path == apis[i].Path && cacheApis[j].Method == apis[i].Method {
				flag = true
			}
		}
		if !flag {
			deleteApis = append(deleteApis, apis[i])
		}
	}
	return
}

func (svc *ApiService) IgnoreApi(ignoreApi data.SysIgnoreApi) (err error) {
	if ignoreApi.Flag {
		return svc.data.SqlClient.Create(&ignoreApi).Error
	}
	return svc.data.SqlClient.Unscoped().Delete(&ignoreApi, "path = ? AND method = ?", ignoreApi.Path, ignoreApi.Method).Error
}

func (svc *ApiService) EnterSyncApi(syncApis dto.SysSyncApis) (err error) {
	return svc.data.SqlClient.Transaction(func(tx *gorm.DB) error {
		var txErr error
		if syncApis.NewApis != nil && len(syncApis.NewApis) > 0 {
			txErr = tx.Create(&syncApis.NewApis).Error
			if txErr != nil {
				return txErr
			}
		}
		for i := range syncApis.DeleteApis {
			CasbinServiceApp.ClearCasbin(1, syncApis.DeleteApis[i].Path, syncApis.DeleteApis[i].Method)
			txErr = tx.Delete(&data.SysApi{}, "path = ? AND method = ?", syncApis.DeleteApis[i].Path, syncApis.DeleteApis[i].Method).Error
			if txErr != nil {
				return txErr
			}
		}
		return nil
	})
}

//@function: DeleteApi
//@description: 删除基础api
//@param: api data.SysApi
//@return: err error

func (svc *ApiService) DeleteApi(api data.SysApi) (err error) {
	var entity data.SysApi
	err = svc.data.SqlClient.First(&entity, "id = ?", api.ID).Error // 根据id查询api记录
	if errors.Is(err, gorm.ErrRecordNotFound) {                     // api记录不存在
		return err
	}
	err = svc.data.SqlClient.Delete(&entity).Error
	if err != nil {
		return err
	}
	CasbinServiceApp.ClearCasbin(1, entity.Path, entity.Method)
	return nil
}

//@function: GetAPIInfoList
//@description: 分页获取数据,
//@param: api data.SysApi, info dto.PageInfo, order string, desc bool
//@return: list interface{}, total int64, err error

func (svc *ApiService) GetAPIInfoList(api data.SysApi, info rhttp.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := svc.data.SqlClient.Model(&data.SysApi{})
	var apiList []data.SysApi

	if api.Path != "" {
		db = db.Where("path LIKE ?", "%"+api.Path+"%")
	}

	if api.Description != "" {
		db = db.Where("description LIKE ?", "%"+api.Description+"%")
	}

	if api.Method != "" {
		db = db.Where("method = ?", api.Method)
	}

	if api.ApiGroup != "" {
		db = db.Where("api_group = ?", api.ApiGroup)
	}

	err = db.Count(&total).Error

	if err != nil {
		return apiList, total, err
	}

	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	if order != "" {
		orderMap := make(map[string]bool, 5)
		orderMap["id"] = true
		orderMap["path"] = true
		orderMap["api_group"] = true
		orderMap["description"] = true
		orderMap["method"] = true
		if !orderMap[order] {
			err = fmt.Errorf("非法的排序字段: %v", order)
			return apiList, total, err
		}
		OrderStr = order
		if desc {
			OrderStr = order + " desc"
		}
	}
	err = db.Order(OrderStr).Find(&apiList).Error
	return apiList, total, err
}

//@function: GetAllApis
//@description: 获取所有的api
//@return:  apis []data.SysApi, err error

func (svc *ApiService) GetAllApis() (apis []data.SysApi, err error) {
	err = svc.data.SqlClient.Order("id desc").Find(&apis).Error
	return
}

//@function: GetApiById
//@description: 根据id获取api
//@param: id float64
//@return: api data.SysApi, err error

func (svc *ApiService) GetApiById(id int) (api data.SysApi, err error) {
	err = svc.data.SqlClient.First(&api, "id = ?", id).Error
	return
}

//@function: UpdateApi
//@description: 根据id更新api
//@param: api data.SysApi
//@return: err error

func (svc *ApiService) UpdateApi(api data.SysApi) (err error) {
	var oldA data.SysApi
	err = svc.data.SqlClient.First(&oldA, "id = ?", api.ID).Error
	if oldA.Path != api.Path || oldA.Method != api.Method {
		var duplicateApi data.SysApi
		if ferr := svc.data.SqlClient.First(&duplicateApi, "path = ? AND method = ?", api.Path, api.Method).Error; ferr != nil {
			if !errors.Is(ferr, gorm.ErrRecordNotFound) {
				return ferr
			}
		} else {
			if duplicateApi.ID != api.ID {
				return errors.New("存在相同api路径")
			}
		}

	}
	if err != nil {
		return err
	}

	err = CasbinServiceApp.UpdateCasbinApi(oldA.Path, api.Path, oldA.Method, api.Method)
	if err != nil {
		return err
	}

	return svc.data.SqlClient.Save(&api).Error
}

//@function: DeleteApisByIds
//@description: 删除选中API
//@param: apis []data.SysApi
//@return: err error

func (svc *ApiService) DeleteApisByIds(ids rhttp.IdsReq) (err error) {
	return svc.data.SqlClient.Transaction(func(tx *gorm.DB) error {
		var apis []data.SysApi
		err = tx.Find(&apis, "id in ?", ids.Ids).Error
		if err != nil {
			return err
		}
		err = tx.Delete(&[]data.SysApi{}, "id in ?", ids.Ids).Error
		if err != nil {
			return err
		}
		for _, sysApi := range apis {
			CasbinServiceApp.ClearCasbin(1, sysApi.Path, sysApi.Method)
		}
		return err
	})
}
