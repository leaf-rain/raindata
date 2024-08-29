package server

import (
	"errors"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"github.com/leaf-rain/raindata/app_bi/internal/data/dto"
	"github.com/leaf-rain/raindata/app_bi/internal/service"
	"gorm.io/gorm"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"go.uber.org/zap"
)

type UserApi struct {
	*Server
}

func NewUserApi(server *Server) *UserApi {
	return &UserApi{
		Server: server,
	}
}

func (svr *UserApi) InitUserAuthRouter(server *Server, router *gin.RouterGroup) {
	api := NewUserApi(server)
	baseRouter := router.Group("base")
	{
		baseRouter.POST("login", api.Login)
	}
}

func (svr *UserApi) InitRouter(server *Server, router *gin.RouterGroup) {
	userServer := NewUserApi(server)
	mid := newMiddleware(server)
	userRouter := router.Group("user").Use(mid.OperationRecord())
	userRouterWithoutRecord := router.Group("user")
	{
		userRouter.POST("admin_register", userServer.Register)               // 管理员注册账号
		userRouter.POST("changePassword", userServer.ChangePassword)         // 用户修改密码
		userRouter.POST("setUserAuthority", userServer.SetUserAuthority)     // 设置用户权限
		userRouter.DELETE("deleteUser", userServer.DeleteUser)               // 删除用户
		userRouter.PUT("setUserInfo", userServer.SetUserInfo)                // 设置用户信息
		userRouter.PUT("setSelfInfo", userServer.SetSelfInfo)                // 设置自身信息
		userRouter.POST("setUserAuthorities", userServer.SetUserAuthorities) // 设置用户权限组
		userRouter.POST("resetPassword", userServer.ResetPassword)           // 重置用户密码
	}
	{
		userRouterWithoutRecord.POST("getUserList", userServer.GetUserList) // 分页获取用户列表
		userRouterWithoutRecord.GET("getUserInfo", userServer.GetUserInfo)  // 获取自身信息
	}
}

// Login
// @Tags     Base
// @Summary  用户登录
// @Produce   application/json
// @Param    data  body      dto.Login                                             true  "用户名, 密码, 验证码"
// @Success  200   {object}  rhttp.Response{data=dto.LoginResponse,msg=string}  "返回包括用户信息,token,过期时间"
// @Router   /base/login [post]
func (svr *UserApi) Login(c *gin.Context) {
	var l dto.Login
	err := c.ShouldBindJSON(&l)
	key := c.ClientIP()
	svr.logger.Info("[Login]登陆请求", zap.String("ip", key), zap.String("username", l.Username), zap.String("password", l.Password))
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = l.Verify()
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}

	u := &data.SysUser{Username: l.Username, Password: l.Password}
	var user *data.SysUser
	userService := service.NewUserService(svr.svc)
	user, err = userService.Login(u)
	if err != nil {
		svr.logger.Error("[Login] 登陆失败! 用户名不存在或者密码错误!", zap.Error(err))
		// 验证码次数+1
		rhttp.FailWithMessage("用户名不存在或者密码错误", c)
		return
	}
	if user.Enable != 1 {
		svr.logger.Error("[Login]登陆失败! 用户被禁止登录!", zap.Uint("userId", user.ID))
		rhttp.FailWithMessage("用户被禁止登录", c)
		return
	}
	svr.TokenNext(c, *user)
	return
}

// TokenNext 登录以后签发jwt
func (svr *UserApi) TokenNext(c *gin.Context, user data.SysUser) {
	jwt, token, claims, err := data.LoginToken(&user, svr.data)
	if err != nil {
		svr.logger.Error("获取token失败!", zap.Error(err))
		rhttp.FailWithMessage("获取token失败", c)
		return
	}
	if !svr.data.Config.GetJwt().GetUseMultipoint() {
		svr.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		rhttp.OkWithDetailed(dto.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
		return
	}

	if jwtStr, err := jwt.GetRedisJWT(c, user.Username); errors.Is(err, redis.Nil) {
		if err := jwt.SetRedisJWT(token, user.Username); err != nil {
			svr.logger.Error("设置登录状态失败!", zap.Error(err))
			rhttp.FailWithMessage("设置登录状态失败", c)
			return
		}
		svr.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		rhttp.OkWithDetailed(dto.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
	} else if err != nil {
		svr.logger.Error("设置登录状态失败!", zap.Error(err))
		rhttp.FailWithMessage("设置登录状态失败", c)
	} else {
		var blackJWT data.JwtBlacklist
		blackJWT.Jwt = jwtStr
		if err := jwt.JsonInBlacklist(blackJWT); err != nil {
			rhttp.FailWithMessage("jwt作废失败", c)
			return
		}
		if err := jwt.SetRedisJWT(token, user.GetUsername()); err != nil {
			rhttp.FailWithMessage("设置登录状态失败", c)
			return
		}
		svr.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		rhttp.OkWithDetailed(dto.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
	}
}

// Register
// @Tags     SysUser
// @Summary  用户注册账号
// @Produce   application/json
// @Param    data  body      dto.Register                                            true  "用户名, 昵称, 密码, 角色ID"
// @Success  200   {object}  rhttp.Response{data=dto.SysUserResponse,msg=string}  "用户注册账号,返回包括用户信息"
// @Router   /user/admin_register [post]
func (svr *UserApi) Register(c *gin.Context) {
	var r dto.Register
	err := c.ShouldBindJSON(&r)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	var authorities []data.SysAuthority
	for _, v := range r.AuthorityIds {
		authorities = append(authorities, data.SysAuthority{
			AuthorityId: v,
		})
	}
	user := &data.SysUser{Username: r.Username, NickName: r.NickName, Password: r.Password, HeaderImg: r.HeaderImg, AuthorityId: r.AuthorityId, Authorities: authorities, Enable: r.Enable, Phone: r.Phone, Email: r.Email}
	userService := service.NewUserService(svr.svc)
	userReturn, err := userService.Register(*user)
	if err != nil {
		svr.logger.Error("注册失败!", zap.Error(err))
		rhttp.FailWithDetailed(dto.SysUserResponse{User: userReturn}, "注册失败", c)
		return
	}
	rhttp.OkWithDetailed(dto.SysUserResponse{User: userReturn}, "注册成功", c)
}

// ChangePassword
// @Tags      SysUser
// @Summary   用户修改密码
// @Security  ApiKeyAuth
// @Produce  application/json
// @Param     data  body      dto.ChangePasswordReq    true  "用户名, 原密码, 新密码"
// @Success   200   {object}  rhttp.Response{msg=string}  "用户修改密码"
// @Router    /user/changePassword [post]
func (svr *UserApi) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}

	//err = req.Verify()
	//if err != nil {
	//	rhttp.FailWithMessage(err.Error(), c)
	//	return
	//}
	uid := svr.GetUserID(c)
	u := &data.SysUser{Model: gorm.Model{ID: uid}, Password: req.Password}
	userService := service.NewUserService(svr.svc)
	_, err = userService.ChangePassword(u, req.NewPassword)
	if err != nil {
		svr.logger.Error("修改失败!", zap.Error(err))
		rhttp.FailWithMessage("修改失败，原密码与当前账户不符", c)
		return
	}
	rhttp.OkWithMessage("修改成功", c)
}

// GetUserList
// @Tags      SysUser
// @Summary   分页获取用户列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      rhttp.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  rhttp.Response{data=rhttp.PageResult,msg=string}  "分页获取用户列表,返回包括列表,总数,页码,每页数量"
// @Router    /user/getUserList [post]
func (svr *UserApi) GetUserList(c *gin.Context) {
	var pageInfo rhttp.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = pageInfo.Verify()
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	userService := service.NewUserService(svr.svc)
	list, total, err := userService.GetUserInfoList(pageInfo)
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

// SetUserAuthority
// @Tags      SysUser
// @Summary   更改用户权限
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      dto.SetUserAuth          true  "用户UUID, 角色ID"
// @Success   200   {object}  rhttp.Response{msg=string}  "设置用户权限"
// @Router    /user/setUserAuthority [post]
func (svr *UserApi) SetUserAuthority(c *gin.Context) {
	var sua dto.SetUserAuth
	err := c.ShouldBindJSON(&sua)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	//if UserVerifyErr := sua.Verify(); UserVerifyErr != nil {
	//	rhttp.FailWithMessage(UserVerifyErr.Error(), c)
	//	return
	//}
	userID := svr.GetUserID(c)
	userService := service.NewUserService(svr.svc)
	err = userService.SetUserAuthority(userID, sua.AuthorityId)
	if err != nil {
		svr.logger.Error("修改失败!", zap.Error(err))
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	claims := svr.GetUserInfoByCtx(c)
	j := data.NewJWT(svr.data) // 唯一签名
	claims.AuthorityId = sua.AuthorityId
	if token, err := j.CreateToken(*claims); err != nil {
		svr.logger.Error("修改失败!", zap.Error(err))
		rhttp.FailWithMessage(err.Error(), c)
	} else {
		c.Header("new-token", token)
		c.Header("new-expires-at", strconv.FormatInt(claims.ExpiresAt.Unix(), 10))
		svr.SetToken(c, token, int((claims.ExpiresAt.Unix()-time.Now().Unix())/60))
		rhttp.OkWithMessage("修改成功", c)
	}
}

// SetUserAuthorities
// @Tags      SysUser
// @Summary   设置用户权限
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      dto.SetUserAuthorities   true  "用户UUID, 角色ID"
// @Success   200   {object}  rhttp.Response{msg=string}  "设置用户权限"
// @Router    /user/setUserAuthorities [post]
func (svr *UserApi) SetUserAuthorities(c *gin.Context) {
	var sua dto.SetUserAuthorities
	err := c.ShouldBindJSON(&sua)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	userService := service.NewUserService(svr.svc)
	err = userService.SetUserAuthorities(sua.ID, sua.AuthorityIds)
	if err != nil {
		svr.logger.Error("修改失败!", zap.Error(err))
		rhttp.FailWithMessage("修改失败", c)
		return
	}
	rhttp.OkWithMessage("修改成功", c)
}

// DeleteUser
// @Tags      SysUser
// @Summary   删除用户
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      rhttp.GetById                true  "用户ID"
// @Success   200   {object}  rhttp.Response{msg=string}  "删除用户"
// @Router    /user/deleteUser [delete]
func (svr *UserApi) DeleteUser(c *gin.Context) {
	var reqId rhttp.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = reqId.Verify()
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	jwtId := svr.GetUserID(c)
	if jwtId == uint(reqId.ID) {
		rhttp.FailWithMessage("删除失败, 自杀失败", c)
		return
	}
	userService := service.NewUserService(svr.svc)
	err = userService.DeleteUser(reqId.ID)
	if err != nil {
		svr.logger.Error("删除失败!", zap.Error(err))
		rhttp.FailWithMessage("删除失败", c)
		return
	}
	rhttp.OkWithMessage("删除成功", c)
}

// SetUserInfo
// @Tags      SysUser
// @Summary   设置用户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      data.SysUser                                             true  "ID, 用户名, 昵称, 头像链接"
// @Success   200   {object}  rhttp.Response{data=map[string]interface{},msg=string}  "设置用户信息"
// @Router    /user/setUserInfo [put]
func (svr *UserApi) SetUserInfo(c *gin.Context) {
	var user dto.ChangeUserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	//err = user.Verify()
	//if err != nil {
	//	rhttp.FailWithMessage(err.Error(), c)
	//	return
	//}
	userService := service.NewUserService(svr.svc)
	if len(user.AuthorityIds) != 0 {
		err = userService.SetUserAuthorities(user.ID, user.AuthorityIds)
		if err != nil {
			svr.logger.Error("设置失败!", zap.Error(err))
			rhttp.FailWithMessage("设置失败", c)
			return
		}
	}
	err = userService.SetUserInfo(data.SysUser{
		Model: gorm.Model{
			ID: user.ID,
		},
		NickName:  user.NickName,
		HeaderImg: user.HeaderImg,
		Phone:     user.Phone,
		Email:     user.Email,
		SideMode:  user.SideMode,
		Enable:    user.Enable,
	})
	if err != nil {
		svr.logger.Error("设置失败!", zap.Error(err))
		rhttp.FailWithMessage("设置失败", c)
		return
	}
	rhttp.OkWithMessage("设置成功", c)
}

// SetSelfInfo
// @Tags      SysUser
// @Summary   设置用户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      data.SysUser                                             true  "ID, 用户名, 昵称, 头像链接"
// @Success   200   {object}  rhttp.Response{data=map[string]interface{},msg=string}  "设置用户信息"
// @Router    /user/SetSelfInfo [put]
func (svr *UserApi) SetSelfInfo(c *gin.Context) {
	var user dto.ChangeUserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	userService := service.NewUserService(svr.svc)
	user.ID = svr.GetUserID(c)
	err = userService.SetSelfInfo(data.SysUser{
		Model: gorm.Model{
			ID: user.ID,
		},
		NickName:  user.NickName,
		HeaderImg: user.HeaderImg,
		Phone:     user.Phone,
		Email:     user.Email,
		SideMode:  user.SideMode,
		Enable:    user.Enable,
	})
	if err != nil {
		svr.logger.Error("设置失败!", zap.Error(err))
		rhttp.FailWithMessage("设置失败", c)
		return
	}
	rhttp.OkWithMessage("设置成功", c)
}

// GetUserInfo
// @Tags      SysUser
// @Summary   获取用户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  rhttp.Response{data=map[string]interface{},msg=string}  "获取用户信息"
// @Router    /user/getUserInfo [get]
func (svr *UserApi) GetUserInfo(c *gin.Context) {
	uid := svr.GetUserID(c)
	userService := service.NewUserService(svr.svc)
	ReqUser, err := userService.FindUserById(uid)
	if err != nil {
		svr.logger.Error("获取失败!", zap.Error(err))
		rhttp.FailWithMessage("获取失败", c)
		return
	}
	rhttp.OkWithDetailed(gin.H{"userInfo": ReqUser}, "获取成功", c)
}

// ResetPassword
// @Tags      SysUser
// @Summary   重置用户密码
// @Security  ApiKeyAuth
// @Produce  application/json
// @Param     data  body      data.SysUser                 true  "ID"
// @Success   200   {object}  rhttp.Response{msg=string}  "重置用户密码"
// @Router    /user/resetPassword [post]
func (svr *UserApi) ResetPassword(c *gin.Context) {
	var user data.SysUser
	err := c.ShouldBindJSON(&user)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	userService := service.NewUserService(svr.svc)
	err = userService.ResetPassword(user.ID)
	if err != nil {
		svr.logger.Error("重置失败!", zap.Error(err))
		rhttp.FailWithMessage("重置失败"+err.Error(), c)
		return
	}
	rhttp.OkWithMessage("重置成功", c)
}
