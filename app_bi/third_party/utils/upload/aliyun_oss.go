package upload

import (
	"errors"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"mime/multipart"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.uber.org/zap"
)

type AliyunOSS struct {
	conf   *conf.Bootstrap
	logger *zap.Logger
}

func (u *AliyunOSS) UploadFile(file *multipart.FileHeader) (string, string, error) {
	bucket, err := u.NewBucket()
	if err != nil {
		u.logger.Error("function AliyunOSS.NewBucket() Failed", zap.Any("err", err.Error()))
		return "", "", errors.New("function AliyunOSS.NewBucket() Failed, err:" + err.Error())
	}

	// 读取本地文件。
	f, openError := file.Open()
	if openError != nil {
		u.logger.Error("function file.Open() Failed", zap.Any("err", openError.Error()))
		return "", "", errors.New("function file.Open() Failed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭
	// 上传阿里云路径 文件名格式 自己可以改 建议保证唯一性
	// yunFileTmpPath := filepath.Join("uploads", time.Now().Format("2006-01-02")) + "/" + file.Filename
	yunFileTmpPath := u.conf.Oss.Aliyun.BasePath + "/" + "uploads" + "/" + time.Now().Format("2006-01-02") + "/" + file.Filename

	// 上传文件流。
	err = bucket.PutObject(yunFileTmpPath, f)
	if err != nil {
		u.logger.Error("function formUploader.Put() Failed", zap.Any("err", err.Error()))
		return "", "", errors.New("function formUploader.Put() Failed, err:" + err.Error())
	}

	return u.conf.Oss.Aliyun.BucketUrl + "/" + yunFileTmpPath, yunFileTmpPath, nil
}

func (u *AliyunOSS) DeleteFile(key string) error {
	bucket, err := u.NewBucket()
	if err != nil {
		u.logger.Error("function AliyunOSS.NewBucket() Failed", zap.Any("err", err.Error()))
		return errors.New("function AliyunOSS.NewBucket() Failed, err:" + err.Error())
	}

	// 删除单个文件。objectName表示删除OSS文件时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	// 如需删除文件夹，请将objectName设置为对应的文件夹名称。如果文件夹非空，则需要将文件夹下的所有object删除后才能删除该文件夹。
	err = bucket.DeleteObject(key)
	if err != nil {
		u.logger.Error("function bucketManager.Delete() failed", zap.Any("err", err.Error()))
		return errors.New("function bucketManager.Delete() failed, err:" + err.Error())
	}

	return nil
}

func (u *AliyunOSS) NewBucket() (*oss.Bucket, error) {
	// 创建OSSClient实例。
	client, err := oss.New(u.conf.Oss.Aliyun.Endpoint, u.conf.Oss.Aliyun.AccessKeyId, u.conf.Oss.Aliyun.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	// 获取存储空间。
	bucket, err := client.Bucket(u.conf.Oss.Aliyun.BucketName)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}
