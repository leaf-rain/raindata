package server

import (
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"github.com/leaf-rain/raindata/app_bi/internal/data/dto"
	"github.com/leaf-rain/raindata/app_bi/internal/service"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"github.com/leaf-rain/raindata/app_bi/third_party/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthorityApi struct {
	*Server
}

// CreateAuthority
// @Tags      Authority
// @Summary   创建角色
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysAuthority                                                true  "权限id, 权限名, 父角色id"
// @Success   200   {object}  rhttp.Response{data=systemRes.SysAuthorityResponse,msg=string}  "创建角色,返回包括系统角色详情"
// @Router    /authority/createAuthority [post]
func (svr *AuthorityApi) CreateAuthority(c *gin.Context) {
	var authority, authBack data.SysAuthority
	var err error

	if err = c.ShouldBindJSON(&authority); err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}

	if err = utils.Verify(authority, utils.AuthorityVerify); err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	authService := service.NewAuthorityService(svr.svc)
	if authBack, err = authService.CreateAuthority(authority); err != nil {
		svr.logger.Error("创建失败!", zap.Error(err))
		rhttp.FailWithMessage("创建失败"+err.Error(), c)
		return
	}
	casbinService := service.NewCasbinService(svr.svc)
	err = casbinService.FreshCasbin()
	if err != nil {
		svr.logger.Error("创建成功，权限刷新失败。", zap.Error(err))
		rhttp.FailWithMessage("创建成功，权限刷新失败。"+err.Error(), c)
		return
	}
	rhttp.OkWithDetailed(dto.SysAuthorityResponse{Authority: authBack}, "创建成功", c)
}

// DeleteAuthority
// @Tags      Authority
// @Summary   删除角色
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysAuthority            true  "删除角色"
// @Success   200   {object}  rhttp.Response{msg=string}  "删除角色"
// @Router    /authority/deleteAuthority [post]
func (svr *AuthorityApi) DeleteAuthority(c *gin.Context) {
	var authority data.SysAuthority
	var err error
	if err = c.ShouldBindJSON(&authority); err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	if err = utils.Verify(authority, utils.AuthorityIdVerify); err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	// 删除角色之前需要判断是否有用户正在使用此角色
	authService := service.NewAuthorityService(svr.svc)
	if err = authService.DeleteAuthority(&authority); err != nil {
		svr.logger.Error("删除失败!", zap.Error(err))
		rhttp.FailWithMessage("删除失败"+err.Error(), c)
		return
	}
	casbinService := service.NewCasbinService(svr.svc)
	_ = casbinService.FreshCasbin()
	rhttp.OkWithMessage("删除成功", c)
}

// UpdateAuthority
// @Tags      Authority
// @Summary   更新角色信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysAuthority                                                true  "权限id, 权限名, 父角色id"
// @Success   200   {object}  rhttp.Response{data=systemRes.SysAuthorityResponse,msg=string}  "更新角色信息,返回包括系统角色详情"
// @Router    /authority/updateAuthority [post]
func (svr *AuthorityApi) UpdateAuthority(c *gin.Context) {
	var auth data.SysAuthority
	err := c.ShouldBindJSON(&auth)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(auth, utils.AuthorityVerify)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	authService := service.NewAuthorityService(svr.svc)
	authority, err := authService.UpdateAuthority(auth)
	if err != nil {
		svr.logger.Error("更新失败!", zap.Error(err))
		rhttp.FailWithMessage("更新失败"+err.Error(), c)
		return
	}
	rhttp.OkWithDetailed(dto.SysAuthorityResponse{Authority: authority}, "更新成功", c)
}

// GetAuthorityList
// @Tags      Authority
// @Summary   分页获取角色列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      dto.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  rhttp.Response{data=rhttp.PageResult,msg=string}  "分页获取角色列表,返回包括列表,总数,页码,每页数量"
// @Router    /authority/getAuthorityList [post]
func (svr *AuthorityApi) GetAuthorityList(c *gin.Context) {
	var pageInfo rhttp.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	authService := service.NewAuthorityService(svr.svc)
	list, total, err := authService.GetAuthorityInfoList(pageInfo)
	if err != nil {
		svr.logger.Error("获取失败!", zap.Error(err))
		rhttp.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	rhttp.OkWithDetailed(rhttp.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// SetDataAuthority
// @Tags      Authority
// @Summary   设置角色资源权限
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysAuthority            true  "设置角色资源权限"
// @Success   200   {object}  rhttp.Response{msg=string}  "设置角色资源权限"
// @Router    /authority/setDataAuthority [post]
func (svr *AuthorityApi) SetDataAuthority(c *gin.Context) {
	var auth data.SysAuthority
	err := c.ShouldBindJSON(&auth)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(auth, utils.AuthorityIdVerify)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	authService := service.NewAuthorityService(svr.svc)
	err = authService.SetDataAuthority(auth)
	if err != nil {
		svr.logger.Error("设置失败!", zap.Error(err))
		rhttp.FailWithMessage("设置失败"+err.Error(), c)
		return
	}
	rhttp.OkWithMessage("设置成功", c)
}
