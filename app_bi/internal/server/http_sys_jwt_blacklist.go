package server

import (
	"github.com/gin-gonic/gin"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"github.com/leaf-rain/raindata/app_bi/internal/service"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"go.uber.org/zap"
)

// JsonInBlacklist
// @Tags      Jwt
// @Summary   jwt加入黑名单
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  rhttp.Response{msg=string}  "jwt加入黑名单"
// @Router    /jwt.proto/jsonInBlacklist [post]
func (svr *UserApi) JsonInBlacklist(c *gin.Context) {
	token := svr.GetToken(c)
	jwt := data.JwtBlacklist{Jwt: token}
	jwtService := service.NewJwtService(svr.svc)
	err := jwtService.JsonInBlacklist(jwt)
	if err != nil {
		svr.logger.Error("jwt作废失败!", zap.Error(err))
		rhttp.FailWithMessage("jwt作废失败", c)
		return
	}
	svr.ClearToken(c)
	rhttp.OkWithMessage("jwt作废成功", c)
}
