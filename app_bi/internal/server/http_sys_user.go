package server

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/leaf-rain/raindata/app_bi/internal/data"
	"github.com/leaf-rain/raindata/app_bi/internal/service"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"go.uber.org/zap"
)

func InitUserAuthRouter(Router *gin.RouterGroup, userSvc *service.UserService, log *zap.Logger, data *data.Data) {
	api := NewUserApi(userSvc, log, data)

	baseRouter := Router.Group("base")
	{
		baseRouter.POST("login", api.Login)
		//baseRouter.POST("captcha", api.Captcha)
	}
}

func InitUserRouter(Router *gin.RouterGroup, userSvc *service.UserService, log *zap.Logger, data *data.Data) {
	api := NewUserApi(userSvc, log, data)

	//userRouter := Router.Group("user").Use(middleware.OperationRecord())
	userRouter := Router.Group("user")
	userRouterWithoutRecord := Router.Group("user")
	{
		userRouter.POST("admin_register", api.Register)               // 管理员注册账号
		userRouter.POST("changePassword", api.ChangePassword)         // 用户修改密码
		userRouter.POST("setUserAuthority", api.SetUserAuthority)     // 设置用户权限
		userRouter.DELETE("deleteUser", api.DeleteUser)               // 删除用户
		userRouter.PUT("setUserInfo", api.SetUserInfo)                // 设置用户信息
		userRouter.PUT("setSelfInfo", api.SetSelfInfo)                // 设置自身信息
		userRouter.POST("setUserAuthorities", api.SetUserAuthorities) // 设置用户权限组
		userRouter.POST("resetPassword", api.ResetPassword)           // 设置用户权限组
	}
	{
		userRouterWithoutRecord.POST("getUserList", api.GetUserList) // 分页获取用户列表
		userRouterWithoutRecord.GET("getUserInfo", api.GetUserInfo)  // 获取自身信息
	}
}

type UserApi struct {
	userSvc *service.UserService
	log     *zap.Logger
	data    *data.Data
}

func NewUserApi(userSvc *service.UserService, log *zap.Logger, data *data.Data) *UserApi {
	return &UserApi{
		userSvc: userSvc,
		log:     log,
		data:    data,
	}
}

// Login
// @Tags     Base
// @Summary  用户登录
// @Produce   application/json
// @Param    data  body      data.LoginReq                                             true  "用户名, 密码, 验证码"
// @Success  200   {object}  rhttp.Response{data=systemRes.LoginResponse,msg=string}  "返回包括用户信息,token,过期时间"
// @Router   /base/login [post]
func (b *UserApi) Login(c *gin.Context) {
	var l data.LoginReq
	err := c.ShouldBindJSON(&l)
	key := c.ClientIP()
	b.log.Info("[Login]登陆请求", zap.String("ip", key), zap.String("username", l.Username), zap.String("password", l.Password))
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
	user, err = b.userSvc.Login(u)
	if err != nil {
		b.log.Error("[Login] 登陆失败! 用户名不存在或者密码错误!", zap.Error(err))
		// 验证码次数+1
		rhttp.FailWithMessage("用户名不存在或者密码错误", c)
		return
	}
	if user.Enable != 1 {
		b.log.Error("[Login]登陆失败! 用户被禁止登录!", zap.String("userId", user.UUID.String()))
		rhttp.FailWithMessage("用户被禁止登录", c)
		return
	}
	b.TokenNext(c, *user)
	return
}

// TokenNext 登录以后签发jwt
func (b *UserApi) TokenNext(c *gin.Context, user data.SysUser) {
	jwt, token, claims, err := data.LoginToken(&user, b.data)
	if err != nil {
		b.log.Error("获取token失败!", zap.Error(err))
		rhttp.FailWithMessage("获取token失败", c)
		return
	}
	if !b.data.Config.GetJwt().GetUseMultipoint() {
		b.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		rhttp.OkWithDetailed(data.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
		return
	}

	if jwtStr, err := jwt.GetRedisJWT(c, user.Username); errors.Is(err, redis.Nil) {
		if err := jwt.SetRedisJWT(token, user.Username); err != nil {
			b.log.Error("设置登录状态失败!", zap.Error(err))
			rhttp.FailWithMessage("设置登录状态失败", c)
			return
		}
		b.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		rhttp.OkWithDetailed(data.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
	} else if err != nil {
		b.log.Error("设置登录状态失败!", zap.Error(err))
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
		b.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		rhttp.OkWithDetailed(data.LoginResponse{
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
// @Param    data  body      systemReq.Register                                            true  "用户名, 昵称, 密码, 角色ID"
// @Success  200   {object}  rhttp.Response{data=systemRes.SysUserResponse,msg=string}  "用户注册账号,返回包括用户信息"
// @Router   /user/admin_register [post]
func (b *UserApi) Register(c *gin.Context) {
	var r data.Register
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
	userReturn, err := b.userSvc.Register(*user)
	if err != nil {
		b.log.Error("注册失败!", zap.Error(err))
		rhttp.FailWithDetailed(data.SysUserResponse{User: userReturn}, "注册失败", c)
		return
	}
	rhttp.OkWithDetailed(data.SysUserResponse{User: userReturn}, "注册成功", c)
}

// ChangePassword
// @Tags      SysUser
// @Summary   用户修改密码
// @Security  ApiKeyAuth
// @Produce  application/json
// @Param     data  body      systemReq.ChangePasswordReq    true  "用户名, 原密码, 新密码"
// @Success   200   {object}  rhttp.Response{msg=string}  "用户修改密码"
// @Router    /user/changePassword [post]
func (b *UserApi) ChangePassword(c *gin.Context) {
	var req data.ChangePasswordReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}

	err = req.Verify()
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	uid := b.GetUserID(c)
	u := &data.SysUser{GVA_MODEL: data.GVA_MODEL{ID: uid}, Password: req.Password}
	_, err = b.userSvc.ChangePassword(u, req.NewPassword)
	if err != nil {
		b.log.Error("修改失败!", zap.Error(err))
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
// @Param     data  body      dto.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  rhttp.Response{data=rhttp.PageResult,msg=string}  "分页获取用户列表,返回包括列表,总数,页码,每页数量"
// @Router    /user/getUserList [post]
func (b *UserApi) GetUserList(c *gin.Context) {
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
	list, total, err := b.userSvc.GetUserInfoList(pageInfo)
	if err != nil {
		b.log.Error("获取失败!", zap.Error(err))
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
// @Param     data  body      systemReq.SetUserAuth          true  "用户UUID, 角色ID"
// @Success   200   {object}  rhttp.Response{msg=string}  "设置用户权限"
// @Router    /user/setUserAuthority [post]
func (b *UserApi) SetUserAuthority(c *gin.Context) {
	var sua data.SetUserAuth
	err := c.ShouldBindJSON(&sua)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	if UserVerifyErr := sua.Verify(); UserVerifyErr != nil {
		rhttp.FailWithMessage(UserVerifyErr.Error(), c)
		return
	}
	userID := b.GetUserID(c)
	err = b.userSvc.SetUserAuthority(userID, sua.AuthorityId)
	if err != nil {
		b.log.Error("修改失败!", zap.Error(err))
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	claims := b.GetUserInfoByCtx(c)
	j := data.NewJWT(b.data) // 唯一签名
	claims.AuthorityId = sua.AuthorityId
	if token, err := j.CreateToken(*claims); err != nil {
		b.log.Error("修改失败!", zap.Error(err))
		rhttp.FailWithMessage(err.Error(), c)
	} else {
		c.Header("new-token", token)
		c.Header("new-expires-at", strconv.FormatInt(claims.ExpiresAt.Unix(), 10))
		b.SetToken(c, token, int((claims.ExpiresAt.Unix()-time.Now().Unix())/60))
		rhttp.OkWithMessage("修改成功", c)
	}
}

// SetUserAuthorities
// @Tags      SysUser
// @Summary   设置用户权限
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.SetUserAuthorities   true  "用户UUID, 角色ID"
// @Success   200   {object}  rhttp.Response{msg=string}  "设置用户权限"
// @Router    /user/setUserAuthorities [post]
func (b *UserApi) SetUserAuthorities(c *gin.Context) {
	var sua data.SetUserAuthorities
	err := c.ShouldBindJSON(&sua)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = b.userSvc.SetUserAuthorities(sua.ID, sua.AuthorityIds)
	if err != nil {
		b.log.Error("修改失败!", zap.Error(err))
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
// @Param     data  body      dto.GetById                true  "用户ID"
// @Success   200   {object}  rhttp.Response{msg=string}  "删除用户"
// @Router    /user/deleteUser [delete]
func (b *UserApi) DeleteUser(c *gin.Context) {
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
	jwtId := b.GetUserID(c)
	if jwtId == uint(reqId.ID) {
		rhttp.FailWithMessage("删除失败, 自杀失败", c)
		return
	}
	err = b.userSvc.DeleteUser(reqId.ID)
	if err != nil {
		b.log.Error("删除失败!", zap.Error(err))
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
// @Param     data  body      system.SysUser                                             true  "ID, 用户名, 昵称, 头像链接"
// @Success   200   {object}  rhttp.Response{data=map[string]interface{},msg=string}  "设置用户信息"
// @Router    /user/setUserInfo [put]
func (b *UserApi) SetUserInfo(c *gin.Context) {
	var user data.ChangeUserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = user.Verify()
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}

	if len(user.AuthorityIds) != 0 {
		err = b.userSvc.SetUserAuthorities(user.ID, user.AuthorityIds)
		if err != nil {
			b.log.Error("设置失败!", zap.Error(err))
			rhttp.FailWithMessage("设置失败", c)
			return
		}
	}
	err = b.userSvc.SetUserInfo(data.SysUser{
		GVA_MODEL: data.GVA_MODEL{
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
		b.log.Error("设置失败!", zap.Error(err))
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
// @Param     data  body      system.SysUser                                             true  "ID, 用户名, 昵称, 头像链接"
// @Success   200   {object}  rhttp.Response{data=map[string]interface{},msg=string}  "设置用户信息"
// @Router    /user/SetSelfInfo [put]
func (b *UserApi) SetSelfInfo(c *gin.Context) {
	var user data.ChangeUserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	user.ID = b.GetUserID(c)
	err = b.userSvc.SetSelfInfo(data.SysUser{
		GVA_MODEL: data.GVA_MODEL{
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
		b.log.Error("设置失败!", zap.Error(err))
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
func (b *UserApi) GetUserInfo(c *gin.Context) {
	uuid := b.GetUserUuid(c)
	ReqUser, err := b.userSvc.GetUserInfo(uuid)
	if err != nil {
		b.log.Error("获取失败!", zap.Error(err))
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
// @Param     data  body      system.SysUser                 true  "ID"
// @Success   200   {object}  rhttp.Response{msg=string}  "重置用户密码"
// @Router    /user/resetPassword [post]
func (b *UserApi) ResetPassword(c *gin.Context) {
	var user data.SysUser
	err := c.ShouldBindJSON(&user)
	if err != nil {
		rhttp.FailWithMessage(err.Error(), c)
		return
	}
	err = b.userSvc.ResetPassword(user.ID)
	if err != nil {
		b.log.Error("重置失败!", zap.Error(err))
		rhttp.FailWithMessage("重置失败"+err.Error(), c)
		return
	}
	rhttp.OkWithMessage("重置成功", c)
}
