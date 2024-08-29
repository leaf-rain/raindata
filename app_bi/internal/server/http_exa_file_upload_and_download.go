package server

import (
	"github.com/gin-gonic/gin"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"github.com/leaf-rain/raindata/app_bi/internal/service"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"go.uber.org/zap"
)

func newFileUploadAndDownloadApi(server *Server) *fileUploadAndDownloadApi {
	return &fileUploadAndDownloadApi{
		Server: server,
	}
}

type fileUploadAndDownloadApi struct {
	*Server
}

func (svr *fileUploadAndDownloadApi) InitRouter(Router *gin.RouterGroup) {
	fileUploadAndDownloadRouter := Router.Group("fileUploadAndDownload")
	{
		fileUploadAndDownloadRouter.POST("upload", svr.UploadFile)         // 上传文件
		fileUploadAndDownloadRouter.POST("getFileList", svr.GetFileList)   // 获取上传文件列表
		fileUploadAndDownloadRouter.POST("deleteFile", svr.DeleteFile)     // 删除指定文件
		fileUploadAndDownloadRouter.POST("editFileName", svr.EditFileName) // 编辑文件名或者备注
	}
}

// UploadFile
// @Tags      ExaFileUploadAndDownload
// @Summary   上传文件示例
// @Security  ApiKeyAuth
// @accept    multipart/form-data
// @Produce   application/json
// @Param     file  formData  file                                                           true  "上传文件示例"
// @Success   200   {object}  rhttp.Response{data=data.ExaFileUploadAndDownload,msg=string}  "上传文件示例,返回包括文件详情"
// @Router    /fileUploadAndDownload/upload [post]
func (svr *fileUploadAndDownloadApi) UploadFile(c *gin.Context) {
	var file data.ExaFileUploadAndDownload
	noSave := c.DefaultQuery("noSave", "0")
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		svr.logger.Error("接收文件失败!", zap.Error(err))
		rhttp.FailWithMessage("接收文件失败", c)
		return
	}
	fileService := service.NewFileUploadAndDownloadService(svr.svc)
	file, err = fileService.UploadFile(header, noSave) // 文件上传后拿到文件路径
	if err != nil {
		svr.logger.Error("修改数据库链接失败!", zap.Error(err))
		rhttp.FailWithMessage("修改数据库链接失败", c)
		return
	}
	rhttp.OkWithDetailed(file, "上传成功", c)
}

// EditFileName 编辑文件名或者备注
func (svr *fileUploadAndDownloadApi) EditFileName(c *gin.Context) {
	var file data.ExaFileUploadAndDownload
	err := c.ShouldBindJSON(&file)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	fileService := service.NewFileUploadAndDownloadService(svr.svc)
	err = fileService.EditFileName(file)
	if err != nil {
		svr.logger.Error("编辑失败!", zap.Error(err))
		rhttp.FailWithMessage("编辑失败", c)
		return
	}
	rhttp.OkWithMessage("编辑成功", c)
}

// DeleteFile
// @Tags      ExaFileUploadAndDownload
// @Summary   删除文件
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      data.ExaFileUploadAndDownload  true  "传入文件里面id即可"
// @Success   200   {object}  rhttp.Response{msg=string}     "删除文件"
// @Router    /fileUploadAndDownload/deleteFile [post]
func (svr *fileUploadAndDownloadApi) DeleteFile(c *gin.Context) {
	var file data.ExaFileUploadAndDownload
	err := c.ShouldBindJSON(&file)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	fileService := service.NewFileUploadAndDownloadService(svr.svc)
	if err = fileService.DeleteFile(file); err != nil {
		svr.logger.Error("删除失败!", zap.Error(err))
		rhttp.FailWithMessage("删除失败", c)
		return
	}
	rhttp.OkWithMessage("删除成功", c)
}

// GetFileList
// @Tags      ExaFileUploadAndDownload
// @Summary   分页文件列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      rhttp.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  rhttp.Response{data=rhttp.PageResult,msg=string}  "分页文件列表,返回包括列表,总数,页码,每页数量"
// @Router    /fileUploadAndDownload/getFileList [post]
func (svr *fileUploadAndDownloadApi) GetFileList(c *gin.Context) {
	var pageInfo rhttp.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	fileService := service.NewFileUploadAndDownloadService(svr.svc)
	list, total, err := fileService.GetFileRecordInfoList(pageInfo)
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
