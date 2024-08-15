package server

import (
	"github.com/gin-gonic/gin"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data"
	"github.com/leaf-rain/raindata/app_bi/internal/data/dto"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"github.com/leaf-rain/raindata/app_bi/third_party/utils"
	"go.uber.org/zap"
)

type CasbinApi struct {
	data *data.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

// UpdateCasbin
// @Tags      Casbin
// @Summary   更新角色api权限
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      dto.CasbinInReceive        true  "权限id, 权限模型列表"
// @Success   200   {object}  rhttp.Response{msg=string}  "更新角色api权限"
// @Router    /casbin/UpdateCasbin [post]
func (cas *CasbinApi) UpdateCasbin(c *gin.Context) {
	var cmr dto.CasbinInReceive
	err := c.ShouldBindJSON(&cmr)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(cmr, utils.AuthorityIdVerify)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = casbinService.UpdateCasbin(cmr.AuthorityId, cmr.CasbinInfos)
	if err != nil {
		svr.log.Error("更新失败!", zap.Error(err))
		rhttp.FailWithMessage("更新失败", c)
		return
	}
	rhttp.OkWithMessage("更新成功", c)
}

// GetPolicyPathByAuthorityId
// @Tags      Casbin
// @Summary   获取权限列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      dto.CasbinInReceive                                          true  "权限id, 权限模型列表"
// @Success   200   {object}  rhttp.Response{data=systemRes.PolicyPathResponse,msg=string}  "获取权限列表,返回包括casbin详情列表"
// @Router    /casbin/getPolicyPathByAuthorityId [post]
func (cas *CasbinApi) GetPolicyPathByAuthorityId(c *gin.Context) {
	var casbin dto.CasbinInReceive
	err := c.ShouldBindJSON(&casbin)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(casbin, utils.AuthorityIdVerify)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	paths := casbinService.GetPolicyPathByAuthorityId(casbin.AuthorityId)
	rhttp.OkWithDetailed(systemRes.PolicyPathResponse{Paths: paths}, "获取成功", c)
}
