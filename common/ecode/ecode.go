package ecode

import "errors"

var (
	ERR_USER_AUTH     = errors.New("用户密码错误")
	ERR_USER_NOTFOUND = errors.New("用户不存在")

	ERR_CONFIG_PATH      = errors.New("配置路径错误")
	ERR_CONFIG_UNMARSHAL = errors.New("配置解析错误")
	ERR_HTTP_CONFIG      = errors.New("http配置解析错误")
	ERR_GRPC_CONFIG      = errors.New("grpc配置解析错误")
	ERR_EMAIL_TYPE       = errors.New("邮件类型不支持")
	ERR_EMAIL_Id         = errors.New("邮件ID错误")
	ERR_APP_ROUTER       = errors.New("方法未找到")

	ERR_MSG_NOT_APPID = errors.New("消息中未找到appid")
)

var ( // db
	ERR_DB_MISSING_DB_CONTEXT        = errors.New("missing db in context")
	ERR_DB_MISSING_DEPENDENT_CONTEXT = errors.New("missing dependent value in context")
	ERR_DB_TYPE_MISMATCH             = errors.New("db type mismatch")
)

var ( // token
	ERR_TOKEN_EXPIRED       = errors.New("Token is expired")
	ERR_TOKEN_NOT_VALID_YET = errors.New("Token not active yet")
	ERR_TOKEN_MALFORMED     = errors.New("That's not even a token")
	ERR_TOKEN_INVALID       = errors.New("Couldn't handle this token:")
)
