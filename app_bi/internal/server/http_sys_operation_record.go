package server

import (
	"github.com/gin-gonic/gin"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"github.com/leaf-rain/raindata/app_bi/internal/data/dto"
	"github.com/leaf-rain/raindata/app_bi/internal/service"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"github.com/leaf-rain/raindata/app_bi/third_party/utils"
	"go.uber.org/zap"
)

type operationRecordApi struct {
	*Server
}

func newOperationRecordApi(server *Server) *operationRecordApi {
	return &operationRecordApi{
		Server: server,
	}
}

func (svr *operationRecordApi) InitRouter(Router *gin.RouterGroup) {
	operationRecordRouter := Router.Group("sysOperationRecord")
	{
		operationRecordRouter.POST("createSysOperationRecord", svr.CreateSysOperationRecord)             // 新建SysOperationRecord
		operationRecordRouter.DELETE("deleteSysOperationRecord", svr.DeleteSysOperationRecord)           // 删除SysOperationRecord
		operationRecordRouter.DELETE("deleteSysOperationRecordByIds", svr.DeleteSysOperationRecordByIds) // 批量删除SysOperationRecord
		operationRecordRouter.GET("findSysOperationRecord", svr.FindSysOperationRecord)                  // 根据ID获取SysOperationRecord
		operationRecordRouter.GET("getSysOperationRecordList", svr.GetSysOperationRecordList)            // 获取SysOperationRecord列表

	}
}

// CreateSysOperationRecord
// @Tags      SysOperationRecord
// @Summary   创建SysOperationRecord
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      data.SysOperationRecord      true  "创建SysOperationRecord"
// @Success   200   {object}  rhttp.Response{msg=string}  "创建SysOperationRecord"
// @Router    /sysOperationRecord/createSysOperationRecord [post]
func (svr *operationRecordApi) CreateSysOperationRecord(c *gin.Context) {
	var sysOperationRecord data.SysOperationRecord
	err := c.ShouldBindJSON(&sysOperationRecord)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	operationRecordService := service.NewOperationRecordService(svr.svc)
	err = operationRecordService.CreateSysOperationRecord(sysOperationRecord)
	if err != nil {
		svr.logger.Error("创建失败!", zap.Error(err))
		rhttp.FailWithMessage("创建失败", c)
		return
	}
	rhttp.OkWithMessage("创建成功", c)
}

// DeleteSysOperationRecord
// @Tags      SysOperationRecord
// @Summary   删除SysOperationRecord
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      data.SysOperationRecord      true  "SysOperationRecord模型"
// @Success   200   {object}  rhttp.Response{msg=string}  "删除SysOperationRecord"
// @Router    /sysOperationRecord/deleteSysOperationRecord [delete]
func (svr *operationRecordApi) DeleteSysOperationRecord(c *gin.Context) {
	var sysOperationRecord data.SysOperationRecord
	err := c.ShouldBindJSON(&sysOperationRecord)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	operationRecordService := service.NewOperationRecordService(svr.svc)
	err = operationRecordService.DeleteSysOperationRecord(sysOperationRecord)
	if err != nil {
		svr.logger.Error("删除失败!", zap.Error(err))
		rhttp.FailWithMessage("删除失败", c)
		return
	}
	rhttp.OkWithMessage("删除成功", c)
}

// DeleteSysOperationRecordByIds
// @Tags      SysOperationRecord
// @Summary   批量删除SysOperationRecord
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      rhttp.IdsReq                 true  "批量删除SysOperationRecord"
// @Success   200   {object}  rhttp.Response{msg=string}  "批量删除SysOperationRecord"
// @Router    /sysOperationRecord/deleteSysOperationRecordByIds [delete]
func (svr *operationRecordApi) DeleteSysOperationRecordByIds(c *gin.Context) {
	var IDS rhttp.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	operationRecordService := service.NewOperationRecordService(svr.svc)
	err = operationRecordService.DeleteSysOperationRecordByIds(IDS)
	if err != nil {
		svr.logger.Error("批量删除失败!", zap.Error(err))
		rhttp.FailWithMessage("批量删除失败", c)
		return
	}
	rhttp.OkWithMessage("批量删除成功", c)
}

// FindSysOperationRecord
// @Tags      SysOperationRecord
// @Summary   用id查询SysOperationRecord
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     data.SysOperationRecord                                  true  "Id"
// @Success   200   {object}  rhttp.Response{data=map[string]interface{},msg=string}  "用id查询SysOperationRecord"
// @Router    /sysOperationRecord/findSysOperationRecord [get]
func (svr *operationRecordApi) FindSysOperationRecord(c *gin.Context) {
	var sysOperationRecord data.SysOperationRecord
	err := c.ShouldBindQuery(&sysOperationRecord)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(sysOperationRecord, utils.IdVerify)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	operationRecordService := service.NewOperationRecordService(svr.svc)
	reSysOperationRecord, err := operationRecordService.GetSysOperationRecord(sysOperationRecord.ID)
	if err != nil {
		svr.logger.Error("查询失败!", zap.Error(err))
		rhttp.FailWithMessage("查询失败", c)
		return
	}
	rhttp.OkWithDetailed(gin.H{"reSysOperationRecord": reSysOperationRecord}, "查询成功", c)
}

// GetSysOperationRecordList
// @Tags      SysOperationRecord
// @Summary   分页获取SysOperationRecord列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     dto.SysOperationRecordSearch                        true  "页码, 每页大小, 搜索条件"
// @Success   200   {object}  rhttp.Response{data=rhttp.PageResult,msg=string}  "分页获取SysOperationRecord列表,返回包括列表,总数,页码,每页数量"
// @Router    /sysOperationRecord/getSysOperationRecordList [get]
func (svr *operationRecordApi) GetSysOperationRecordList(c *gin.Context) {
	var pageInfo dto.SysOperationRecordSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	operationRecordService := service.NewOperationRecordService(svr.svc)
	list, total, err := operationRecordService.GetSysOperationRecordInfoList(pageInfo)
	if err != nil {
		svr.logger.Error("获取失败!", zap.Error(err))
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
