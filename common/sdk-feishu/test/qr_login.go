package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// 环境变量
var (
	APP_ID             = "cli_a72d9f7014f89013"
	APP_SECRET         = "GxB2VAGUYx93eNSHIyZOzfyxP3hlrfBT"
	LARK_PASSPORT_HOST = "https://passport.feishu.cn/suite/passport/oauth/"
	LARK_BASE_URL      = "https://open.feishu.cn/open-apis"
	REDIRECT_URI       = "http://120.76.192.134:3001/qrLogin/"
)

// QrLogin 结构体
type QrLogin struct {
	larkPassportHost string
	larkBaseURL      string
	appID            string
	appSecret        string
	tokenInfo        map[string]interface{}
	userInfo         map[string]interface{}
}

// NewQrLogin 构造函数
func NewQrLogin(appID, appSecret, larkPassportHost, larkBaseURL string) *QrLogin {
	return &QrLogin{
		larkPassportHost: larkPassportHost,
		larkBaseURL:      larkBaseURL,
		appID:            appID,
		appSecret:        appSecret,
		tokenInfo:        make(map[string]interface{}),
		userInfo:         make(map[string]interface{}),
	}
}

// AppAccessToken 获取应用访问令牌
func (q *QrLogin) AppAccessToken() (string, error) {
	param := map[string]string{
		"app_id":     q.appID,
		"app_secret": q.appSecret,
	}
	headers := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}

	body, err := q.postRequest("/auth/v3/app_access_token/internal", param, headers)
	if err != nil {
		return "", err
	}

	var tokenRes map[string]interface{}
	if err := json.Unmarshal(body, &tokenRes); err != nil {
		return "", err
	}

	tenantAccessToken, ok := tokenRes["tenant_access_token"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get tenant_access_token")
	}

	return tenantAccessToken, nil
}

// GetTokenInfo 获取令牌信息
func (q *QrLogin) GetTokenInfo(jsonParam map[string]interface{}) (map[string]interface{}, error) {
	urlStr, ok := jsonParam["url"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid url parameter")
	}

	urlParamResult, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	queryDict := urlParamResult.Query()
	log.Println("query: ---------->", urlStr)
	codeArr, ok := queryDict["code"]
	if !ok || len(codeArr) == 0 {
		return nil, nil
	}
	log.Println("code: ---------->", queryDict)

	appAccessToken, err := q.AppAccessToken()
	if err != nil {
		return nil, err
	}
	log.Println("appAccessToken: ---------->", appAccessToken)

	headers := map[string]string{
		"Content-Type":  "application/json; charset=utf-8",
		"Authorization": "Bearer " + appAccessToken,
	}

	payload := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     APP_ID,
		"client_secret": APP_SECRET,
		"redirect_uri":  REDIRECT_URI,
		"code":          codeArr[0],
	}

	body, err := q.postRequest("/authen/v2/oauth/token", payload, headers)
	if err != nil {
		return nil, err
	}
	log.Println("user access key: ---------->", string(body))

	var tokenRes map[string]interface{}
	if err := json.Unmarshal(body, &tokenRes); err != nil {
		return nil, err
	}

	if tokenRes["code"] != 0 {
		return nil, nil
	}

	tokenInfo := map[string]interface{}{
		"accessToken": tokenRes["access_token"],
		"tokenType":   tokenRes["token_type"],
	}

	q.tokenInfo = tokenInfo
	return tokenInfo, nil
}

// GetUserInfo 获取用户信息
func (q *QrLogin) GetUserInfo() (map[string]interface{}, error) {
	log.Println("accessToken: ---------->", q.tokenInfo)
	accessToken, _ := q.tokenInfo["accessToken"].(string)
	if accessToken == "" {
		return nil, fmt.Errorf("accessToken is empty")
	}
	header := map[string]string{
		"Content-Type":  "application/json;charset=UTF-8",
		"Authorization": "Bearer " + accessToken,
	}

	body, err := q.getRequest("/authen/v1/user_info", header)
	if err != nil {
		return nil, err
	}

	var userInfoObj map[string]interface{}
	if err := json.Unmarshal(body, &userInfoObj); err != nil {
		return nil, err
	}

	if userInfoObj["code"] != 0 {
		return nil, nil
	}

	data, ok := userInfoObj["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format")
	}

	qrUserInfo := map[string]interface{}{
		"name":      data["name"],
		"openId":    data["open_id"],
		"userId":    data["open_id"],
		"tenantKey": data["tenant_key"],
		"avatarUrl": data["avatar_url"],
	}

	q.userInfo = qrUserInfo
	return qrUserInfo, nil
}

// GenURL 生成URL
func (q *QrLogin) GenURL(uri string) string {
	return q.larkBaseURL + uri
}

// postRequest 发送POST请求
func (q *QrLogin) postRequest(uri string, payload map[string]string, headers map[string]string) ([]byte, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", q.GenURL(uri), strings.NewReader(string(jsonPayload)))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// getRequest 发送GET请求
func (q *QrLogin) getRequest(uri string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", q.GenURL(uri), nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// https://open.feishu.cn/open-apis/authen/v1/index?redirect_uri=http://120.76.192.134:3001&app_id=cli_a72d9f7014f89013&state=RANDOMSTRING

//  https://passport.feishu.cn/suite/passport/oauth/authorize?redirect_uri=http://120.76.192.134:3001/qrLogin&app_id=cli_a72d9f7014f89013&state=RANDOMSTRING
//  https://passport.feishu.cn/suite/passport/oauth/authorize?redirect_uri=&app_id=cli_a72d9f7014f89013&state=RANDOMSTRING
//  https://passport.feishu.cn/suite/passport/oauth/authorize?redirect_uri=http%3A%2F%2F120.76.192.134%3A3001%2FqrLogin&app_id=cli_a72d9f7014f89013&state=RANDOMSTRING

//  http://120.76.192.134:3001/qrLogin/

// qrLogin 处理QR登录请求
func qrLogin(c *gin.Context) {
	if c.Request.Method == "POST" || c.Request.Method == "GET" {
		var jsonParam map[string]interface{}
		if err := c.ShouldBindJSON(&jsonParam); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "Invalid JSON body"})
			return
		}

		qrLogin := NewQrLogin(APP_ID, APP_SECRET, LARK_PASSPORT_HOST, LARK_BASE_URL)

		tokenInfo, err := qrLogin.GetTokenInfo(jsonParam)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": err.Error()})
			return
		}

		qrUserInfo, err := qrLogin.GetUserInfo()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": err.Error()})
			return
		}

		if len(tokenInfo) == 0 && len(qrUserInfo) == 0 {
			c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "Scan qr code to obtain user information."})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":       0,
			"msg":        "get userinfo success",
			"tokenInfo":  tokenInfo,
			"qrUserInfo": qrUserInfo,
		})
	}
}

func main() {
	r := gin.Default()
	r.RedirectTrailingSlash = false
	r.Use(CORS())
	r.Use(gin.Logger())
	r.POST("/qrLogin/", qrLogin)
	r.GET("/qrLogin/", qrLogin)
	r.Run(":3000")
}

// https://passport.feishu.cn/suite/passport/oauth/authorize?redirect_uri=http://120.76.192.134:3001/qrLogin&app_id=cli_a72d9f7014f89013
// http://120.76.192.134:3001/qrLogin
//func CORS() gin.HandlerFunc {
//	corsConfig := cors.DefaultConfig()
//	corsConfig.AllowAllOrigins = true
//	corsConfig.AllowCredentials = true
//	corsConfig.OptionsResponseStatusCode = http.StatusOK
//	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
//	corsConfig.AllowHeaders = []string{
//		"Content-Type",
//		"Authorization",
//		"accept-encoding",
//		"authorization",
//		"content-type",
//		"dnt",
//		"origin",
//		"user-agent",
//		"x-csrftoken",
//		"x-requested-with",
//	}
//	return cors.New(corsConfig)
//}

// eyJhbGciOiJFUzI1NiIsImZlYXR1cmVfY29kZSI6IkZlYXR1cmVPQXV0aEpXVFNpZ25fQ04iLCJraWQiOiI3NDY5OTQ0Njk2MDE0Njg0MTY0IiwidHlwIjoiSldUIn0.eyJqdGkiOiI3NDcxMTA5MzU0NDM2MDY3MzMwIiwiaWF0IjoxNzM5NTAzMTk5LCJleHAiOjE3Mzk1MTAzOTksInZlciI6InYxIiwidHlwIjoiYWNjZXNzX3Rva2VuIiwiY2xpZW50X2lkIjoiY2xpX2E3MmQ5ZjcwMTRmODkwMTMiLCJzY29wZSI6ImF1dGg6dXNlci5pZDpyZWFkIHVzZXJfcHJvZmlsZSIsImF1dGhfaWQiOiI3NDcxMTA5MzUwMTM3OTU0MzA2IiwiYXV0aF90aW1lIjoxNzM5NTAzMTk4LCJhdXRoX2V4cCI6MTc3MTAzOTE5OCwidW5pdCI6ImV1X25jIiwidGVuYW50X3VuaXQiOiJldV9uYyIsIm9wYXF1ZSI6dHJ1ZSwiZW5jIjoiQWlRa0FRRUNBTUlEQUFFQkF3QUNBUTBBQXdzTEFBQUFBd0FBQUFkR1pXRjBkWEpsQUFBQUVHOWhkWFJvWDI5d1lYRjFaVjlxZDNRQUFBQUlWR1Z1WVc1MFNXUUFBQUFCTUFBQUFBUlVhVzFsQUFBQUNqRTNNemt4TkRVMk1EQVBBQVFNQUFBQUFRb0FBV0xGSWxtdmdBQWlDd0FDQUFBQURPWklvcVY0cUdyMDlGWnNUQXNBQXdBQUFEQ1ZOZm9oNVlXMlJZc3k0WGdSNXhxL3o4d3VlZ05WamFnb3puT29FSFY5eWxFTjhHLytPYktDaUZxLzlsRzdyLzRBQ3dBRkFBQUFCV1YxWDI1akFQbnU0SGJ4Ym5qVjBOWEhvazdFT0FzcmFqWVRFbDJDRnphZmJtaWJhVlJqbzN0ckc5QzcyUFdFQzlQK05FeXBRbzhvZnMwMHpxRkMxeU1SYWxnV2FpbGpiTmVQRHRWSzIrcUZFZVNDcC9BME5sUEhLYTJQaTdXd0JsOC92MGtPUStoWDA1cWcxa3pvQ1NkdWRpMWUrL2NpV3c9PSIsImVuY192ZXIiOiJ2MSJ9.eDM0e8aZwKpHUGwU7-bHkS5XU9nC5itr25NGO-HyMEHZ7tNVghUuI1Ix0CV0DQ512i2-HxJy0TNb8xtx5DlMtg

//curl --location --request GET 'https://open.feishu.cn/open-apis/authen/v1/user_info' \
//--header 'Authorization: Bearer eyJhbGciOiJFUzI1NiIsImZlYXR1cmVfY29kZSI6IkZlYXR1cmVPQXV0aEpXVFNpZ25fQ04iLCJraWQiOiI3NDY5OTQ0Njk2MDE0Njg0MTY0IiwidHlwIjoiSldUIn0.eyJqdGkiOiI3NDcxMTA5MzU0NDM2MDY3MzMwIiwiaWF0IjoxNzM5NTAzMTk5LCJleHAiOjE3Mzk1MTAzOTksInZlciI6InYxIiwidHlwIjoiYWNjZXNzX3Rva2VuIiwiY2xpZW50X2lkIjoiY2xpX2E3MmQ5ZjcwMTRmODkwMTMiLCJzY29wZSI6ImF1dGg6dXNlci5pZDpyZWFkIHVzZXJfcHJvZmlsZSIsImF1dGhfaWQiOiI3NDcxMTA5MzUwMTM3OTU0MzA2IiwiYXV0aF90aW1lIjoxNzM5NTAzMTk4LCJhdXRoX2V4cCI6MTc3MTAzOTE5OCwidW5pdCI6ImV1X25jIiwidGVuYW50X3VuaXQiOiJldV9uYyIsIm9wYXF1ZSI6dHJ1ZSwiZW5jIjoiQWlRa0FRRUNBTUlEQUFFQkF3QUNBUTBBQXdzTEFBQUFBd0FBQUFkR1pXRjBkWEpsQUFBQUVHOWhkWFJvWDI5d1lYRjFaVjlxZDNRQUFBQUlWR1Z1WVc1MFNXUUFBQUFCTUFBQUFBUlVhVzFsQUFBQUNqRTNNemt4TkRVMk1EQVBBQVFNQUFBQUFRb0FBV0xGSWxtdmdBQWlDd0FDQUFBQURPWklvcVY0cUdyMDlGWnNUQXNBQXdBQUFEQ1ZOZm9oNVlXMlJZc3k0WGdSNXhxL3o4d3VlZ05WamFnb3puT29FSFY5eWxFTjhHLytPYktDaUZxLzlsRzdyLzRBQ3dBRkFBQUFCV1YxWDI1akFQbnU0SGJ4Ym5qVjBOWEhvazdFT0FzcmFqWVRFbDJDRnphZmJtaWJhVlJqbzN0ckc5QzcyUFdFQzlQK05FeXBRbzhvZnMwMHpxRkMxeU1SYWxnV2FpbGpiTmVQRHRWSzIrcUZFZVNDcC9BME5sUEhLYTJQaTdXd0JsOC92MGtPUStoWDA1cWcxa3pvQ1NkdWRpMWUrL2NpV3c9PSIsImVuY192ZXIiOiJ2MSJ9.eDM0e8aZwKpHUGwU7-bHkS5XU9nC5itr25NGO-HyMEHZ7tNVghUuI1Ix0CV0DQ512i2-HxJy0TNb8xtx5DlMtg' \
//--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
//--header 'Content-Type: application/json; charset=utf-8' \
//--header 'Accept: */*' \
//--header 'Host: open.feishu.cn' \
//--header 'Connection: keep-alive' \
//--header 'Cookie: passport_web_did=7470875875244851219; passport_trace_id=7470875875253469203; swp_csrf_token=5018fb27-50de-49cb-b512-fc7f13a6fd44; t_beda37=d8903c4c6a155a1d9367157bd90cca60cb574660ed84ef8cc607e0cf4dbd1493; QXV0aHpDb250ZXh0=f0e2506d5e9b45328f52b5986669128e'
