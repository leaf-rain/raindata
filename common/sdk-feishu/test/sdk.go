package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	larkauthen "github.com/larksuite/oapi-sdk-go/v3/service/authen/v1"
	"golang.org/x/oauth2"
	"log"
	"math/rand"
	"net/http"
)

// SDK 使用文档：https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/server-side-sdk/golang-sdk-guide/preparations
// 复制该 Demo 后, 需要将 "YOUR_APP_ID", "YOUR_APP_SECRET" 替换为自己应用的 APP_ID, APP_SECRET.
// 以下示例代码默认根据文档示例值填充，如果存在代码问题，请在 API 调试台填上相关必要参数后再复制代码使用

var oauthEndpoint = oauth2.Endpoint{
	AuthURL:  "https://accounts.feishu.cn/open-apis/authen/v1/authorize",
	TokenURL: "https://open.feishu.cn/open-apis/authen/v2/oauth/token",
}

func NewFeishuSdk(cfg oauth2.Config, options ...lark.ClientOptionFunc) *SDKFeishu {
	if cfg.Endpoint.AuthURL == "" {
		cfg.Endpoint = oauthEndpoint
	}
	sdk := &SDKFeishu{
		cfg: cfg,
	}
	sdk.client = lark.NewClient(cfg.ClientID, cfg.ClientSecret, options...)
	return sdk
}

type SDKFeishu struct {
	client *lark.Client
	cfg    oauth2.Config
}

func (sdk *SDKFeishu) GetUserInfoV1(accessToken string) (*larkauthen.GetUserInfoResp, error) {
	return sdk.client.Authen.V1.UserInfo.Get(context.Background(), larkcore.WithUserAccessToken(accessToken))
}

func main() {
	sdk := NewFeishuSdk(oauth2.Config{
		ClientID:     "cli_a72d9f7014f89013",
		ClientSecret: "GxB2VAGUYx93eNSHIyZOzfyxP3hlrfBT",
		Endpoint:     oauthEndpoint,
		RedirectURL:  "http://120.76.192.134:3000/callback/",
	})
	r := gin.Default()
	// 使用 Cookie 存储 session
	store := cookie.NewStore([]byte("secret")) // 此处仅为示例，务必不要硬编码密钥
	r.Use(sessions.Sessions("mysession", store))
	r.RedirectTrailingSlash = false
	r.Use(CORS())
	r.Use(gin.Logger())
	r.GET("/", sdk.indexController)
	r.GET("/login", sdk.loginController)
	r.GET("/callback/", sdk.oauthCallbackController)
	log.Fatal(r.Run(":3000"))
}

func CORS() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.OptionsResponseStatusCode = http.StatusOK
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{
		"Content-Type",
		"Authorization",
		"accept-encoding",
		"authorization",
		"content-type",
		"dnt",
		"origin",
		"user-agent",
		"x-csrftoken",
		"x-requested-with",
	}
	return cors.New(corsConfig)
}

func (sdk *SDKFeishu) indexController(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	var username string
	session := sessions.Default(c)
	if session.Get("user") != nil {
		username = session.Get("user").(string)
	}
	html := fmt.Sprintf(`<html><head><style>body{font-family:Arial,sans-serif;background:#f4f4f4;margin:0;display:flex;justify-content:center;align-items:center;height:100vh}.container{text-align:center;background:#fff;padding:30px;border-radius:10px;box-shadow:0 0 10px rgba(0,0,0,0.1)}a{padding:10px 20px;font-size:16px;color:#fff;background:#007bff;border-radius:5px;text-decoration:none;transition:0.3s}a:hover{background:#0056b3}}</style></head><body><div class="container"><h2>欢迎%s！</h2><a href="/login">使用飞书登录</a></div></body></html>`, username)
	c.String(http.StatusOK, html)
}

func (sdk *SDKFeishu) loginController(c *gin.Context) {
	session := sessions.Default(c)
	// 生成随机 state 字符串，你也可以用其他有意义的信息来构建 state
	state := fmt.Sprintf("%d", rand.Int())
	// 将 state 存入 session 中
	session.Set("state", state)
	// 生成 PKCE 需要的 code verifier
	verifier := oauth2.GenerateVerifier()
	// 将 code verifier 存入 session 中
	session.Set("code_verifier", verifier)
	session.Save()
	url := sdk.cfg.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))
	// 用户点击登录时，重定向到授权页面
	c.Redirect(http.StatusTemporaryRedirect, url)
}
func (sdk *SDKFeishu) oauthCallbackController(c *gin.Context) {
	session := sessions.Default(c)
	ctx := context.Background()
	// 从 session 中获取 state
	expectedState := session.Get("state")
	state := c.Query("state")
	// 如果 state 不匹配，说明是 CSRF 攻击，拒绝处理
	if state != expectedState {
		//log.Printf("invalid oauth state, expected '%s', got '%s'\n", expectedState, state)
		//c.Redirect(http.StatusTemporaryRedirect, "/")
		//return
	}
	code := c.Query("code")
	// 如果 code 为空，说明用户拒绝了授权
	if code == "" {
		log.Printf("error: %s", c.Query("error"))
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	codeVerifier, _ := session.Get("code_verifier").(string)
	// 使用获取到的 code 获取 token
	token, err := sdk.cfg.Exchange(ctx, code, oauth2.VerifierOption(codeVerifier))
	if err != nil {
		log.Printf("oauthConfig.Exchange() failed with '%s'\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	userInfo, err := sdk.GetUserInfoV1(token.AccessToken)
	if err != nil {
		log.Printf("sdk.GetUserInfoV1() failed with '%s'\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	} else {
		js, _ := json.Marshal(userInfo)
		log.Printf("oauthCallbackController success, clientIp:%s, remoteIp:%s, userInfo:%s", c.ClientIP(), c.RemoteIP(), string(js))
	}
	// 后续可以用获取到的用户信息构建登录态，此处仅为示例，请勿直接使用
	session.Set("user", userInfo.Data.Name)
	session.Save()
	c.Header("Content-Type", "text/html; charset=utf-8")
	html := fmt.Sprintf(`<html><head><style>body{font-family:Arial,sans-serif;background:#f4f4f4;margin:0;display:flex;justify-content:center;align-items:center;height:100vh}.container{text-align:center;background:#fff;padding:30px;border-radius:10px;box-shadow:0 0 10px rgba(0,0,0,0.1)}a{padding:10px 20px;font-size:16px;color:#fff;background:#007bff;border-radius:5px;text-decoration:none;transition:0.3s}a:hover{background:#0056b3}}</style></head><body><div class="container"><h2>你好，%s！</h2><p>你已成功完成授权登录流程。</p><a href="/">返回主页</a></div></body></html>`, *userInfo.Data.Name)
	c.String(http.StatusOK, html)
}
