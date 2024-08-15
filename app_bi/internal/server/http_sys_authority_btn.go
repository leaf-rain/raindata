package server

import (
	"github.com/gin-gonic/gin"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data"
	"github.com/leaf-rain/raindata/app_bi/internal/data/dto"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"go.uber.org/zap"
)

type AuthorityBtnApi struct {
	data *data.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

// GetAuthorityBtn
// @Tags      AuthorityBtn
// @Summary   获取权限按钮
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      dto.SysAuthorityBtnReq                                      true  "菜单id, 角色id, 选中的按钮id"
// @Success   200   {object}  rhttp.Response{data=rhttp.SysAuthorityBtnRes,msg=string}  "返回列表成功"
// @Router    /authorityBtn/getAuthorityBtn [post]
func (svr *AuthorityBtnApi) GetAuthorityBtn(c *gin.Context) {
	var req dto.SysAuthorityBtnReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	res, err := authorityBtnService.GetAuthorityBtn(req)
	if err != nil {
		svr.log.Error("查询失败!", zap.Error(err))
		rhttp.FailWithMessage("查询失败", c)
		return
	}
	rhttp.OkWithDetailed(res, "查询成功", c)
}

// SetAuthorityBtn
// @Tags      AuthorityBtn
// @Summary   设置权限按钮
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      dto.SysAuthorityBtnReq     true  "菜单id, 角色id, 选中的按钮id"
// @Success   200   {object}  rhttp.Response{msg=string}  "返回列表成功"
// @Router    /authorityBtn/setAuthorityBtn [post]
func (svr *AuthorityBtnApi) SetAuthorityBtn(c *gin.Context) {
	var req dto.SysAuthorityBtnReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = authorityBtnService.SetAuthorityBtn(req)
	if err != nil {
		svr.log.Error("分配失败!", zap.Error(err))
		rhttp.FailWithMessage("分配失败", c)
		return
	}
	rhttp.OkWithMessage("分配成功", c)
}

// CanRemoveAuthorityBtn
// @Tags      AuthorityBtn
// @Summary   设置权限按钮
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  rhttp.Response{msg=string}  "删除成功"
// @Router    /authorityBtn/canRemoveAuthorityBtn [post]
func (svr *AuthorityBtnApi) CanRemoveAuthorityBtn(c *gin.Context) {
	id := c.Query("id")
	err := authorityBtnService.CanRemoveAuthorityBtn(id)
	if err != nil {
		svr.log.Error("删除失败!", zap.Error(err))
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	rhttp.OkWithMessage("删除成功", c)
}
