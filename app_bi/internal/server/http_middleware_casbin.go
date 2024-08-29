package server

import (
	"github.com/gin-gonic/gin"
	"github.com/leaf-rain/raindata/app_bi/internal/service"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"strconv"
)

// CasbinHandler 拦截器
func (mid *middleware) CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userApi = NewUserApi(mid.Server)
		waitUse, _ := userApi.GetClaims(c)
		//获取请求的PATH
		path := c.Request.URL.Path
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		sub := strconv.Itoa(int(waitUse.AuthorityId))
		casbinService := service.NewCasbinService(mid.svc)
		e := casbinService.Casbin() // 判断策略中是否存在
		success, _ := e.Enforce(sub, path, act)
		if !success {
			rhttp.FailWithDetailed(gin.H{}, "权限不足", c)
			c.Abort()
			return
		}
		c.Next()
	}
}
