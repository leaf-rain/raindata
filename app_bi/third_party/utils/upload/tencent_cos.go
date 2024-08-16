package upload

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/tencentyun/cos-go-sdk-v5"
	"go.uber.org/zap"
)

type TencentCOS struct {
	conf   *conf.Bootstrap
	logger *zap.Logger
}

// UploadFile upload file to COS
func (u *TencentCOS) UploadFile(file *multipart.FileHeader) (string, string, error) {
	client := u.NewClient()
	f, openError := file.Open()
	if openError != nil {
		u.logger.Error("function file.Open() failed", zap.Any("err", openError.Error()))
		return "", "", errors.New("function file.Open() failed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename)

	_, err := client.Object.Put(context.Background(), u.conf.Oss.TencentCOS.PathPrefix+"/"+fileKey, f, nil)
	if err != nil {
		panic(err)
	}
	return u.conf.Oss.TencentCOS.BaseURL + "/" + u.conf.Oss.TencentCOS.PathPrefix + "/" + fileKey, fileKey, nil
}

// DeleteFile delete file form COS
func (u *TencentCOS) DeleteFile(key string) error {
	client := u.NewClient()
	name := u.conf.Oss.TencentCOS.PathPrefix + "/" + key
	_, err := client.Object.Delete(context.Background(), name)
	if err != nil {
		u.logger.Error("function bucketManager.Delete() failed", zap.Any("err", err.Error()))
		return errors.New("function bucketManager.Delete() failed, err:" + err.Error())
	}
	return nil
}

// NewClient init COS client
func (u *TencentCOS) NewClient() *cos.Client {
	urlStr, _ := url.Parse("https://" + u.conf.Oss.TencentCOS.Bucket + ".cos." + u.conf.Oss.TencentCOS.Region + ".myqcloud.com")
	baseURL := &cos.BaseURL{BucketURL: urlStr}
	client := cos.NewClient(baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  u.conf.Oss.TencentCOS.SecretID,
			SecretKey: u.conf.Oss.TencentCOS.SecretKey,
		},
	})
	return client
}
