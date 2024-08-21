package server

import (
	"github.com/gin-gonic/gin"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"github.com/leaf-rain/raindata/app_bi/internal/data/dto"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"github.com/leaf-rain/raindata/app_bi/third_party/utils"
	"go.uber.org/zap"
)

type DictionaryDetailApi struct {
	data *data.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

// CreateSysDictionaryDetail
// @Tags      SysDictionaryDetail
// @Summary   创建SysDictionaryDetail
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysDictionaryDetail     true  "SysDictionaryDetail模型"
// @Success   200   {object}  rhttp.Response{msg=string}  "创建SysDictionaryDetail"
// @Router    /sysDictionaryDetail/createSysDictionaryDetail [post]
func (s *DictionaryDetailApi) CreateSysDictionaryDetail(c *gin.Context) {
	var detail system.SysDictionaryDetail
	err := c.ShouldBindJSON(&detail)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = dictionaryDetailService.CreateSysDictionaryDetail(detail)
	if err != nil {
		svr.log.Error("创建失败!", zap.Error(err))
		rhttp.FailWithMessage("创建失败", c)
		return
	}
	rhttp.OkWithMessage("创建成功", c)
}

// DeleteSysDictionaryDetail
// @Tags      SysDictionaryDetail
// @Summary   删除SysDictionaryDetail
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysDictionaryDetail     true  "SysDictionaryDetail模型"
// @Success   200   {object}  rhttp.Response{msg=string}  "删除SysDictionaryDetail"
// @Router    /sysDictionaryDetail/deleteSysDictionaryDetail [delete]
func (s *DictionaryDetailApi) DeleteSysDictionaryDetail(c *gin.Context) {
	var detail system.SysDictionaryDetail
	err := c.ShouldBindJSON(&detail)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = dictionaryDetailService.DeleteSysDictionaryDetail(detail)
	if err != nil {
		svr.log.Error("删除失败!", zap.Error(err))
		rhttp.FailWithMessage("删除失败", c)
		return
	}
	rhttp.OkWithMessage("删除成功", c)
}

// UpdateSysDictionaryDetail
// @Tags      SysDictionaryDetail
// @Summary   更新SysDictionaryDetail
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysDictionaryDetail     true  "更新SysDictionaryDetail"
// @Success   200   {object}  rhttp.Response{msg=string}  "更新SysDictionaryDetail"
// @Router    /sysDictionaryDetail/updateSysDictionaryDetail [put]
func (s *DictionaryDetailApi) UpdateSysDictionaryDetail(c *gin.Context) {
	var detail system.SysDictionaryDetail
	err := c.ShouldBindJSON(&detail)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = dictionaryDetailService.UpdateSysDictionaryDetail(&detail)
	if err != nil {
		svr.log.Error("更新失败!", zap.Error(err))
		rhttp.FailWithMessage("更新失败", c)
		return
	}
	rhttp.OkWithMessage("更新成功", c)
}

// FindSysDictionaryDetail
// @Tags      SysDictionaryDetail
// @Summary   用id查询SysDictionaryDetail
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     system.SysDictionaryDetail                                 true  "用id查询SysDictionaryDetail"
// @Success   200   {object}  rhttp.Response{data=map[string]interface{},msg=string}  "用id查询SysDictionaryDetail"
// @Router    /sysDictionaryDetail/findSysDictionaryDetail [get]
func (s *DictionaryDetailApi) FindSysDictionaryDetail(c *gin.Context) {
	var detail system.SysDictionaryDetail
	err := c.ShouldBindQuery(&detail)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(detail, utils.IdVerify)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	reSysDictionaryDetail, err := dictionaryDetailService.GetSysDictionaryDetail(detail.ID)
	if err != nil {
		svr.log.Error("查询失败!", zap.Error(err))
		rhttp.FailWithMessage("查询失败", c)
		return
	}
	rhttp.OkWithDetailed(gin.H{"reSysDictionaryDetail": reSysDictionaryDetail}, "查询成功", c)
}

// GetSysDictionaryDetailList
// @Tags      SysDictionaryDetail
// @Summary   分页获取SysDictionaryDetail列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     dto.SysDictionaryDetailSearch                       true  "页码, 每页大小, 搜索条件"
// @Success   200   {object}  rhttp.Response{data=rhttp.PageResult,msg=string}  "分页获取SysDictionaryDetail列表,返回包括列表,总数,页码,每页数量"
// @Router    /sysDictionaryDetail/getSysDictionaryDetailList [get]
func (s *DictionaryDetailApi) GetSysDictionaryDetailList(c *gin.Context) {
	var pageInfo dto.SysDictionaryDetailSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := dictionaryDetailService.GetSysDictionaryDetailInfoList(pageInfo)
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
