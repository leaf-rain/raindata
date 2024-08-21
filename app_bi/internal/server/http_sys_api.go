package server

import (
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"github.com/leaf-rain/raindata/app_bi/third_party/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemApiApi struct {
	data *data.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

// CreateApi
// @Tags      SysApi
// @Summary   创建基础api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      data.SysApi                  true  "api路径, api中文描述, api组, 方法"
// @Success   200   {object}  rhttp.Response{msg=string}  "创建基础api"
// @Router    /api/createApi [post]
func (svr *SystemApiApi) CreateApi(c *gin.Context) {
	var api data.SysApi
	err := c.ShouldBindJSON(&api)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(api, utils.ApiVerify)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.CreateApi(api)
	if err != nil {
		svr.log.Error("创建失败!", zap.Error(err))
		rhttp.FailWithMessage("创建失败", c)
		return
	}
	rhttp.OkWithMessage("创建成功", c)
}

// SyncApi
// @Tags      SysApi
// @Summary   同步API
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  rhttp.Response{msg=string}  "同步API"
// @Router    /api/syncApi [get]
func (svr *SystemApiApi) SyncApi(c *gin.Context) {
	newApis, deleteApis, ignoreApis, err := apiService.SyncApi()
	if err != nil {
		svr.log.Error("同步失败!", zap.Error(err))
		rhttp.FailWithMessage("同步失败", c)
		return
	}
	rhttp.OkWithData(gin.H{
		"newApis":    newApis,
		"deleteApis": deleteApis,
		"ignoreApis": ignoreApis,
	}, c)
}

// GetApiGroups
// @Tags      SysApi
// @Summary   获取API分组
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  rhttp.Response{msg=string}  "获取API分组"
// @Router    /api/getApiGroups [post]
func (svr *SystemApiApi) GetApiGroups(c *gin.Context) {
	groups, apiGroupMap, err := apiService.GetApiGroups()
	if err != nil {
		svr.log.Error("获取失败!", zap.Error(err))
		rhttp.FailWithMessage("获取失败", c)
		return
	}
	rhttp.OkWithData(gin.H{
		"groups":      groups,
		"apiGroupMap": apiGroupMap,
	}, c)
}

// IgnoreApi
// @Tags      IgnoreApi
// @Summary   忽略API
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  rhttp.Response{msg=string}  "同步API"
// @Router    /api/ignoreApi [post]
func (svr *SystemApiApi) IgnoreApi(c *gin.Context) {
	var ignoreApi data.SysIgnoreApi
	err := c.ShouldBindJSON(&ignoreApi)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.IgnoreApi(ignoreApi)
	if err != nil {
		svr.log.Error("忽略失败!", zap.Error(err))
		rhttp.FailWithMessage("忽略失败", c)
		return
	}
	rhttp.Ok(c)
}

// EnterSyncApi
// @Tags      SysApi
// @Summary   确认同步API
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  rhttp.Response{msg=string}  "确认同步API"
// @Router    /api/enterSyncApi [post]
func (svr *SystemApiApi) EnterSyncApi(c *gin.Context) {
	var syncApi systemRes.SysSyncApis
	err := c.ShouldBindJSON(&syncApi)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.EnterSyncApi(syncApi)
	if err != nil {
		svr.log.Error("忽略失败!", zap.Error(err))
		rhttp.FailWithMessage("忽略失败", c)
		return
	}
	rhttp.Ok(c)
}

// DeleteApi
// @Tags      SysApi
// @Summary   删除api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      data.SysApi                  true  "ID"
// @Success   200   {object}  rhttp.Response{msg=string}  "删除api"
// @Router    /api/deleteApi [post]
func (svr *SystemApiApi) DeleteApi(c *gin.Context) {
	var api data.SysApi
	err := c.ShouldBindJSON(&api)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(api.gorm.Model, utils.IdVerify)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.DeleteApi(api)
	if err != nil {
		svr.log.Error("删除失败!", zap.Error(err))
		rhttp.FailWithMessage("删除失败", c)
		return
	}
	rhttp.OkWithMessage("删除成功", c)
}

// GetApiList
// @Tags      SysApi
// @Summary   分页获取API列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.SearchApiParams                               true  "分页获取API列表"
// @Success   200   {object}  rhttp.Response{data=rhttp.PageResult,msg=string}  "分页获取API列表,返回包括列表,总数,页码,每页数量"
// @Router    /api/getApiList [post]
func (svr *SystemApiApi) GetApiList(c *gin.Context) {
	var pageInfo systemReq.SearchApiParams
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := apiService.GetAPIInfoList(pageInfo.SysApi, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc)
	if err != nil {
		svr.log.Error("获取失败!", zap.Error(err))
		rhttp.FailWithMessage("获取失败", c)
		return
	}
	rhttp.OkWithDetailed(rhttp.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetApiById
// @Tags      SysApi
// @Summary   根据id获取api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      dto.GetById                                   true  "根据id获取api"
// @Success   200   {object}  rhttp.Response{data=systemRes.SysAPIResponse}  "根据id获取api,返回包括api详情"
// @Router    /api/getApiById [post]
func (svr *SystemApiApi) GetApiById(c *gin.Context) {
	var idInfo dto.GetById
	err := c.ShouldBindJSON(&idInfo)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(idInfo, utils.IdVerify)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	api, err := apiService.GetApiById(idInfo.ID)
	if err != nil {
		svr.log.Error("获取失败!", zap.Error(err))
		rhttp.FailWithMessage("获取失败", c)
		return
	}
	rhttp.OkWithDetailed(systemRes.SysAPIResponse{Api: api}, "获取成功", c)
}

// UpdateApi
// @Tags      SysApi
// @Summary   修改基础api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      data.SysApi                  true  "api路径, api中文描述, api组, 方法"
// @Success   200   {object}  rhttp.Response{msg=string}  "修改基础api"
// @Router    /api/updateApi [post]
func (svr *SystemApiApi) UpdateApi(c *gin.Context) {
	var api data.SysApi
	err := c.ShouldBindJSON(&api)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(api, utils.ApiVerify)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.UpdateApi(api)
	if err != nil {
		svr.log.Error("修改失败!", zap.Error(err))
		rhttp.FailWithMessage("修改失败", c)
		return
	}
	rhttp.OkWithMessage("修改成功", c)
}

// GetAllApis
// @Tags      SysApi
// @Summary   获取所有的Api 不分页
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  rhttp.Response{data=systemRes.SysAPIListResponse,msg=string}  "获取所有的Api 不分页,返回包括api列表"
// @Router    /api/getAllApis [post]
func (svr *SystemApiApi) GetAllApis(c *gin.Context) {
	apis, err := apiService.GetAllApis()
	if err != nil {
		svr.log.Error("获取失败!", zap.Error(err))
		rhttp.FailWithMessage("获取失败", c)
		return
	}
	rhttp.OkWithDetailed(systemRes.SysAPIListResponse{Apis: apis}, "获取成功", c)
}

// DeleteApisByIds
// @Tags      SysApi
// @Summary   删除选中Api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      dto.IdsReq                 true  "ID"
// @Success   200   {object}  rhttp.Response{msg=string}  "删除选中Api"
// @Router    /api/deleteApisByIds [delete]
func (svr *SystemApiApi) DeleteApisByIds(c *gin.Context) {
	var ids dto.IdsReq
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.DeleteApisByIds(ids)
	if err != nil {
		svr.log.Error("删除失败!", zap.Error(err))
		rhttp.FailWithMessage("删除失败", c)
		return
	}
	rhttp.OkWithMessage("删除成功", c)
}

// FreshCasbin
// @Tags      SysApi
// @Summary   刷新casbin缓存
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  rhttp.Response{msg=string}  "刷新成功"
// @Router    /api/freshCasbin [get]
func (svr *SystemApiApi) FreshCasbin(c *gin.Context) {
	err := casbinService.FreshCasbin()
	if err != nil {
		svr.log.Error("刷新失败!", zap.Error(err))
		rhttp.FailWithMessage("刷新失败", c)
		return
	}
	rhttp.OkWithMessage("刷新成功", c)
}
