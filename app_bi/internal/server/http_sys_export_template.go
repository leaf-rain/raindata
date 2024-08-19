package server

import (
	"fmt"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data/entity"
	"github.com/leaf-rain/raindata/app_bi/third_party/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"go.uber.org/zap"
)

type SysExportTemplateApi struct {
	data *entity.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

var sysExportTemplateService = service.ServiceGroupApp.SystemServiceGroup.SysExportTemplateService

// CreateSysExportTemplate 创建导出模板
// @Tags SysExportTemplate
// @Summary 创建导出模板
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysExportTemplate true "创建导出模板"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /sysExportTemplate/createSysExportTemplate [post]
func (sysExportTemplateApi *SysExportTemplateApi) CreateSysExportTemplate(c *gin.Context) {
	var sysExportTemplate system.SysExportTemplate
	err := c.ShouldBindJSON(&sysExportTemplate)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	verify := utils.Rules{
		"Name": {utils.NotEmpty()},
	}
	if err := utils.Verify(sysExportTemplate, verify); err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	if err := sysExportTemplateService.CreateSysExportTemplate(&sysExportTemplate); err != nil {
		svr.log.Error("创建失败!", zap.Error(err))
		rhttp.FailWithMessage("创建失败", c)
	} else {
		rhttp.OkWithMessage("创建成功", c)
	}
}

// DeleteSysExportTemplate 删除导出模板
// @Tags SysExportTemplate
// @Summary 删除导出模板
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysExportTemplate true "删除导出模板"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysExportTemplate/deleteSysExportTemplate [delete]
func (sysExportTemplateApi *SysExportTemplateApi) DeleteSysExportTemplate(c *gin.Context) {
	var sysExportTemplate system.SysExportTemplate
	err := c.ShouldBindJSON(&sysExportTemplate)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	if err := sysExportTemplateService.DeleteSysExportTemplate(sysExportTemplate); err != nil {
		svr.log.Error("删除失败!", zap.Error(err))
		rhttp.FailWithMessage("删除失败", c)
	} else {
		rhttp.OkWithMessage("删除成功", c)
	}
}

// DeleteSysExportTemplateByIds 批量删除导出模板
// @Tags SysExportTemplate
// @Summary 批量删除导出模板
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dto.IdsReq true "批量删除导出模板"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /sysExportTemplate/deleteSysExportTemplateByIds [delete]
func (sysExportTemplateApi *SysExportTemplateApi) DeleteSysExportTemplateByIds(c *gin.Context) {
	var IDS dto.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	if err := sysExportTemplateService.DeleteSysExportTemplateByIds(IDS); err != nil {
		svr.log.Error("批量删除失败!", zap.Error(err))
		rhttp.FailWithMessage("批量删除失败", c)
	} else {
		rhttp.OkWithMessage("批量删除成功", c)
	}
}

// UpdateSysExportTemplate 更新导出模板
// @Tags SysExportTemplate
// @Summary 更新导出模板
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysExportTemplate true "更新导出模板"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sysExportTemplate/updateSysExportTemplate [put]
func (sysExportTemplateApi *SysExportTemplateApi) UpdateSysExportTemplate(c *gin.Context) {
	var sysExportTemplate system.SysExportTemplate
	err := c.ShouldBindJSON(&sysExportTemplate)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	verify := utils.Rules{
		"Name": {utils.NotEmpty()},
	}
	if err := utils.Verify(sysExportTemplate, verify); err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	if err := sysExportTemplateService.UpdateSysExportTemplate(sysExportTemplate); err != nil {
		svr.log.Error("更新失败!", zap.Error(err))
		rhttp.FailWithMessage("更新失败", c)
	} else {
		rhttp.OkWithMessage("更新成功", c)
	}
}

// FindSysExportTemplate 用id查询导出模板
// @Tags SysExportTemplate
// @Summary 用id查询导出模板
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query system.SysExportTemplate true "用id查询导出模板"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /sysExportTemplate/findSysExportTemplate [get]
func (sysExportTemplateApi *SysExportTemplateApi) FindSysExportTemplate(c *gin.Context) {
	var sysExportTemplate system.SysExportTemplate
	err := c.ShouldBindQuery(&sysExportTemplate)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	if resysExportTemplate, err := sysExportTemplateService.GetSysExportTemplate(sysExportTemplate.ID); err != nil {
		svr.log.Error("查询失败!", zap.Error(err))
		rhttp.FailWithMessage("查询失败", c)
	} else {
		rhttp.OkWithData(gin.H{"resysExportTemplate": resysExportTemplate}, c)
	}
}

// GetSysExportTemplateList 分页获取导出模板列表
// @Tags SysExportTemplate
// @Summary 分页获取导出模板列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query systemReq.SysExportTemplateSearch true "分页获取导出模板列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sysExportTemplate/getSysExportTemplateList [get]
func (sysExportTemplateApi *SysExportTemplateApi) GetSysExportTemplateList(c *gin.Context) {
	var pageInfo systemReq.SysExportTemplateSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := sysExportTemplateService.GetSysExportTemplateInfoList(pageInfo); err != nil {
		svr.log.Error("获取失败!", zap.Error(err))
		rhttp.FailWithMessage("获取失败", c)
	} else {
		rhttp.OkWithDetailed(rhttp.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// ExportExcel 导出表格
// @Tags SysExportTemplate
// @Summary 导出表格
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Router /sysExportTemplate/exportExcel [get]
func (sysExportTemplateApi *SysExportTemplateApi) ExportExcel(c *gin.Context) {
	templateID := c.Query("templateID")
	queryParams := c.dto.URL.Query()
	if templateID == "" {
		rhttp.FailWithMessage("模板ID不能为空", c)
		return
	}
	if file, name, err := sysExportTemplateService.ExportExcel(templateID, queryParams); err != nil {
		svr.log.Error("获取失败!", zap.Error(err))
		rhttp.FailWithMessage("获取失败", c)
	} else {
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name+utils.RandomString(6)+".xlsx")) // 对下载的文件重命名
		c.Header("success", "true")
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", file.Bytes())
	}
}

// ExportTemplate 导出表格模板
// @Tags SysExportTemplate
// @Summary 导出表格模板
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Router /sysExportTemplate/exportExcel [get]
func (sysExportTemplateApi *SysExportTemplateApi) ExportTemplate(c *gin.Context) {
	templateID := c.Query("templateID")
	if templateID == "" {
		rhttp.FailWithMessage("模板ID不能为空", c)
		return
	}
	if file, name, err := sysExportTemplateService.ExportTemplate(templateID); err != nil {
		svr.log.Error("获取失败!", zap.Error(err))
		rhttp.FailWithMessage("获取失败", c)
	} else {
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name+"模板.xlsx")) // 对下载的文件重命名
		c.Header("success", "true")
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", file.Bytes())
	}
}

// ImportExcel 导入表格
// @Tags SysImportTemplate
// @Summary 导入表格
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Router /sysExportTemplate/importExcel [post]
func (sysExportTemplateApi *SysExportTemplateApi) ImportExcel(c *gin.Context) {
	templateID := c.Query("templateID")
	if templateID == "" {
		rhttp.FailWithMessage("模板ID不能为空", c)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		svr.log.Error("文件获取失败!", zap.Error(err))
		rhttp.FailWithMessage("文件获取失败", c)
		return
	}
	if err := sysExportTemplateService.ImportExcel(templateID, file); err != nil {
		svr.log.Error(err.Error(), zap.Error(err))
		rhttp.FailWithMessage(err.Error(), c)
	} else {
		rhttp.OkWithMessage("导入成功", c)
	}
}
