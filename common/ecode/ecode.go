package ecode

import "errors"

var (
	ERR_CONFIG_PATH      = errors.New("配置路径错误")
	ERR_CONFIG_UNMARSHAL = errors.New("配置解析错误")
	ERR_HTTP_CONFIG      = errors.New("http配置解析错误")
	ERR_GRPC_CONFIG      = errors.New("grpc配置解析错误")
	ERR_EMAIL_TYPE       = errors.New("邮件类型不支持")
	ERR_EMAIL_Id         = errors.New("邮件ID错误")
	ERR_APP_ROUTER       = errors.New("方法未找到")
)
