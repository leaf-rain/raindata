package server

import (
	"github.com/gin-gonic/gin"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"go.uber.org/zap"
)

type SystemApi struct {
	data *data.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

// GetSystemConfig
// @Tags      System
// @Summary   获取配置文件内容
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200  {object}  rhttp.Response{data=systemRes.SysConfigResponse,msg=string}  "获取配置文件内容,返回包括系统配置"
// @Router    /system/getSystemConfig [post]
func (s *SystemApi) GetSystemConfig(c *gin.Context) {
	config, err := systemConfigService.GetSystemConfig()
	if err != nil {
		svr.log.Error("获取失败!", zap.Error(err))
		rhttp.FailWithMessage("获取失败", c)
		return
	}
	rhttp.OkWithDetailed(systemRes.SysConfigResponse{Config: config}, "获取成功", c)
}

// SetSystemConfig
// @Tags      System
// @Summary   设置配置文件内容
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      system.System                   true  "设置配置文件内容"
// @Success   200   {object}  rhttp.Response{data=string}  "设置配置文件内容"
// @Router    /system/setSystemConfig [post]
func (s *SystemApi) SetSystemConfig(c *gin.Context) {
	var sys system.System
	err := c.ShouldBindJSON(&sys)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = systemConfigService.SetSystemConfig(sys)
	if err != nil {
		svr.log.Error("设置失败!", zap.Error(err))
		rhttp.FailWithMessage("设置失败", c)
		return
	}
	rhttp.OkWithMessage("设置成功", c)
}

// ReloadSystem
// @Tags      System
// @Summary   重启系统
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200  {object}  rhttp.Response{msg=string}  "重启系统"
// @Router    /system/reloadSystem [post]
func (s *SystemApi) ReloadSystem(c *gin.Context) {
	err := utils.Reload()
	if err != nil {
		svr.log.Error("重启系统失败!", zap.Error(err))
		rhttp.FailWithMessage("重启系统失败", c)
		return
	}
	rhttp.OkWithMessage("重启系统成功", c)
}

// GetServerInfo
// @Tags      System
// @Summary   获取服务器信息
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200  {object}  rhttp.Response{data=map[string]interface{},msg=string}  "获取服务器信息"
// @Router    /system/getServerInfo [post]
func (s *SystemApi) GetServerInfo(c *gin.Context) {
	server, err := systemConfigService.GetServerInfo()
	if err != nil {
		svr.log.Error("获取失败!", zap.Error(err))
		rhttp.FailWithMessage("获取失败", c)
		return
	}
	rhttp.OkWithDetailed(gin.H{"server": server}, "获取成功", c)
}
