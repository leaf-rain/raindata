package upload

import (
	"context"
	"errors"
	"fmt"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"mime/multipart"
	"time"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"go.uber.org/zap"
)

type Qiniu struct {
	conf   *conf.Bootstrap
	logger *zap.Logger
}

//@object: u *Qiniu
//@function: UploadFile
//@description: 上传文件
//@param: file *multipart.FileHeader
//@return: string, string, error

func (u *Qiniu) UploadFile(file *multipart.FileHeader) (string, string, error) {
	putPolicy := storage.PutPolicy{Scope: u.conf.Oss.Qiniu.Bucket}
	mac := qbox.NewMac(u.conf.Oss.Qiniu.AccessKey, u.conf.Oss.Qiniu.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := u.qiniuConfig()
	formUploader := storage.NewFormUploader(cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{Params: map[string]string{"x:name": "github logo"}}

	f, openError := file.Open()
	if openError != nil {
		u.logger.Error("function file.Open() failed", zap.Any("err", openError.Error()))

		return "", "", errors.New("function file.Open() failed, err:" + openError.Error())
	}
	defer f.Close()                                                  // 创建文件 defer 关闭
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename) // 文件名格式 自己可以改 建议保证唯一性
	putErr := formUploader.Put(context.Background(), &ret, upToken, fileKey, f, file.Size, &putExtra)
	if putErr != nil {
		u.logger.Error("function formUploader.Put() failed", zap.Any("err", putErr.Error()))
		return "", "", errors.New("function formUploader.Put() failed, err:" + putErr.Error())
	}
	return u.conf.Oss.Qiniu.ImgPath + "/" + ret.Key, ret.Key, nil
}

//@object: u *Qiniu
//@function: DeleteFile
//@description: 删除文件
//@param: key string
//@return: error

func (u *Qiniu) DeleteFile(key string) error {
	mac := qbox.NewMac(u.conf.Oss.Qiniu.AccessKey, u.conf.Oss.Qiniu.SecretKey)
	cfg := u.qiniuConfig()
	bucketManager := storage.NewBucketManager(mac, cfg)
	if err := bucketManager.Delete(u.conf.Oss.Qiniu.Bucket, key); err != nil {
		u.logger.Error("function bucketManager.Delete() failed", zap.Any("err", err.Error()))
		return errors.New("function bucketManager.Delete() failed, err:" + err.Error())
	}
	return nil
}

//@object: u *Qiniu
//@function: qiniuConfig
//@description: 根据配置文件进行返回七牛云的配置
//@return: *storage.Config

func (u *Qiniu) qiniuConfig() *storage.Config {
	cfg := storage.Config{
		UseHTTPS:      u.conf.Oss.Qiniu.UseHTTPS,
		UseCdnDomains: u.conf.Oss.Qiniu.UseCdnDomains,
	}
	switch u.conf.Oss.Qiniu.Zone { // 根据配置文件进行初始化空间对应的机房
	case "ZoneHuadong":
		cfg.Zone = &storage.ZoneHuadong
	case "ZoneHuabei":
		cfg.Zone = &storage.ZoneHuabei
	case "ZoneHuanan":
		cfg.Zone = &storage.ZoneHuanan
	case "ZoneBeimei":
		cfg.Zone = &storage.ZoneBeimei
	case "ZoneXinjiapo":
		cfg.Zone = &storage.ZoneXinjiapo
	}
	return &cfg
}
