package server

import (
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
)

type DBApi struct {
	data *data.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

// InitDB
// @Tags     InitDB
// @Summary  初始化用户数据库
// @Produce  application/json
// @Param    data  body      dto.InitDB                  true  "初始化数据库参数"
// @Success  200   {object}  rhttp.Response{data=string}  "初始化用户数据库"
// @Router   /init/initdb [post]
func (i *DBApi) InitDB(c *gin.Context) {
	if svc.data.SqlClient != nil {
		svr.log.Error("已存在数据库配置!")
		rhttp.FailWithMessage("已存在数据库配置", c)
		return
	}
	var dbInfo dto.InitDB
	if err := c.ShouldBindJSON(&dbInfo); err != nil {
		svr.log.Error("参数校验不通过!", zap.Error(err))
		rhttp.FailWithMessage("参数校验不通过", c)
		return
	}
	if err := initDBService.InitDB(dbInfo); err != nil {
		svr.log.Error("自动创建数据库失败!", zap.Error(err))
		rhttp.FailWithMessage("自动创建数据库失败，请查看后台日志，检查后在进行初始化", c)
		return
	}
	rhttp.OkWithMessage("自动创建数据库成功", c)
}

// CheckDB
// @Tags     CheckDB
// @Summary  初始化用户数据库
// @Produce  application/json
// @Success  200  {object}  rhttp.Response{data=map[string]interface{},msg=string}  "初始化用户数据库"
// @Router   /init/checkdb [post]
func (i *DBApi) CheckDB(c *gin.Context) {
	var (
		message  = "前往初始化数据库"
		needInit = true
	)

	if svc.data.SqlClient != nil {
		message = "数据库无需初始化"
		needInit = false
	}
	svr.log.Info(message)
	rhttp.OkWithDetailed(gin.H{"needInit": needInit}, message, c)
}
