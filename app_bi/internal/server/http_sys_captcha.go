package server

import (
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

// 当开启多服务器部署时，替换下面的配置，使用redis共享存储验证码
// var store = captcha.NewDefaultRedisStore()
var store = base64Captcha.DefaultMemStore

type BaseApi struct {
	data *data.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

// Captcha
// @Tags      Base
// @Summary   生成验证码
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  rhttp.Response{data=systemRes.SysCaptchaResponse,msg=string}  "生成验证码,返回包括随机数id,base64,验证码长度,是否开启验证码"
// @Router    /base/captcha [post]
func (b *BaseApi) Captcha(c *gin.Context) {
	// 判断验证码是否开启
	openCaptcha := svc.conf.Captcha.OpenCaptcha               // 是否开启防爆次数
	openCaptchaTimeOut := svc.conf.Captcha.OpenCaptchaTimeOut // 缓存超时时间
	key := c.ClientIP()
	v, ok := global.BlackCache.Get(key)
	if !ok {
		global.BlackCache.Set(key, 1, time.Second*time.Duration(openCaptchaTimeOut))
	}

	var oc bool
	if openCaptcha == 0 || openCaptcha < interfaceToInt(v) {
		oc = true
	}
	// 字符,公式,验证码配置
	// 生成默认数字的driver
	driver := base64Captcha.NewDriverDigit(svc.conf.Captcha.ImgHeight, svc.conf.Captcha.ImgWidth, svc.conf.Captcha.KeyLong, 0.7, 80)
	// cp := base64Captcha.NewCaptcha(driver, store.UseWithCtx(c))   // v8下使用redis
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err := cp.Generate()
	if err != nil {
		svr.log.Error("验证码获取失败!", zap.Error(err))
		rhttp.FailWithMessage("验证码获取失败", c)
		return
	}
	rhttp.OkWithDetailed(systemRes.SysCaptchaResponse{
		CaptchaId:     id,
		PicPath:       b64s,
		CaptchaLength: svc.conf.Captcha.KeyLong,
		OpenCaptcha:   oc,
	}, "验证码获取成功", c)
}

// 类型转换
func interfaceToInt(v interface{}) (i int) {
	switch v := v.(type) {
	case int:
		i = v
	default:
		i = 0
	}
	return
}
