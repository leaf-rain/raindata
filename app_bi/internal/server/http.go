package server

import (
	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"go.uber.org/zap"
	"net/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, logger *zap.Logger, svr *Server) *khttp.Server {
	engine := gin.Default()
	// 使用kratos中间件
	engine.Use(kgin.Middlewares(recovery.Recovery()))

	publicGroup := engine.Group("/")
	privateGroup := engine.Group("/")

	// todo: private接口校验权限
	//privateGroup.Use()

	{
		// 健康监测
		publicGroup.GET("/health", func(c *gin.Context) {
			var result = make(map[string]interface{}, 6)
			result["status"] = "ok"
			// 获取 CPU 使用率
			cpuPercent, err := cpu.Percent(0, false)
			if err != nil {
				logger.Error("Error getting CPU percent.", zap.Error(err))
			} else {
				logger.Info("CPU Usage:", zap.Float64("percent", cpuPercent[0]))
				result["cpu_percent"] = cpuPercent[0]
			}
			// 获取内存使用情况
			vmStat, err := mem.VirtualMemory()
			if err != nil {
				logger.Error("Error getting memory stats:", zap.Error(err))
			} else {
				result["total_memory"] = vmStat.Total
				result["available_memory"] = vmStat.Available
				result["used_memory"] = vmStat.Used
				result["memory_usage"] = vmStat.UsedPercent
				logger.Info("Total Memory", zap.Uint64("total", vmStat.Total))
				logger.Info("Available Memory", zap.Uint64("available", vmStat.Available))
				logger.Info("Used Memory", zap.Uint64("used", vmStat.Used))
				logger.Info("Memory Usage", zap.Float64("percent", vmStat.UsedPercent))
			}
			c.JSON(http.StatusOK, map[string]interface{}{
				"status": "ok",
			})
		})
	}

	{
		InitUserAuthRouter(svr, publicGroup) // 注册基础功能路由 不做鉴权
	}

	// todo:服务注册
	{
		InitUserRouter(svr, privateGroup)
	}

	httpSrv := khttp.NewServer(khttp.Address(":8000"))
	httpSrv.HandlePrefix("/", engine)
	return httpSrv
}
