package upload

import (
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"go.uber.org/zap"
	"mime/multipart"
)

// OSS 对象存储接口
type OSS interface {
	UploadFile(file *multipart.FileHeader) (string, string, error)
	DeleteFile(key string) error
}

// NewOss OSS的实例化方法
func NewOss(conf *conf.Bootstrap, logger *zap.Logger) OSS {
	switch conf.Oss.Type {
	case "local":
		return &Local{
			conf:   conf,
			logger: logger,
		}
	case "qiniu":
		return &Qiniu{
			conf:   conf,
			logger: logger,
		}
	case "tencent-cos":
		return &TencentCOS{
			conf:   conf,
			logger: logger,
		}
	case "aliyun-oss":
		return &AliyunOSS{
			conf:   conf,
			logger: logger,
		}
	case "huawei-obs":
		return &Obs{
			conf:   conf,
			logger: logger,
		}
	case "aws-s3":
		return &AwsS3{
			conf:   conf,
			logger: logger,
		}
	case "cloudflare-r2":
		return &CloudflareR2{
			conf:   conf,
			logger: logger,
		}
	default:
		return &Local{
			conf:   conf,
			logger: logger,
		}
	}
}
