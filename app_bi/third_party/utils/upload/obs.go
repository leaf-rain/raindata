package upload

import (
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"go.uber.org/zap"
	"mime/multipart"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/pkg/errors"
)

type Obs struct {
	conf   *conf.Bootstrap
	logger *zap.Logger
}

func (u *Obs) NewHuaWeiObsClient() (client *obs.ObsClient, err error) {
	return obs.New(u.conf.Oss.Huaweiyun.AccessKey, u.conf.Oss.Huaweiyun.SecretKey, u.conf.Oss.Huaweiyun.Endpoint)
}

func (u *Obs) UploadFile(file *multipart.FileHeader) (string, string, error) {
	// var open multipart.File
	open, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer open.Close()
	filename := file.Filename
	input := &obs.PutObjectInput{
		PutObjectBasicInput: obs.PutObjectBasicInput{
			ObjectOperationInput: obs.ObjectOperationInput{
				Bucket: u.conf.Oss.Huaweiyun.Bucket,
				Key:    filename,
			},
			HttpHeader: obs.HttpHeader{
				ContentType: file.Header.Get("content-type"),
			},
		},
		Body: open,
	}

	var client *obs.ObsClient
	client, err = u.NewHuaWeiObsClient()
	if err != nil {
		return "", "", errors.Wrap(err, "获取华为对象存储对象失败!")
	}

	_, err = client.PutObject(input)
	if err != nil {
		return "", "", errors.Wrap(err, "文件上传失败!")
	}
	filepath := u.conf.Oss.Huaweiyun.Path + "/" + filename
	return filepath, filename, err
}

func (u *Obs) DeleteFile(key string) error {
	client, err := u.NewHuaWeiObsClient()
	if err != nil {
		return errors.Wrap(err, "获取华为对象存储对象失败!")
	}
	input := &obs.DeleteObjectInput{
		Bucket: u.conf.Oss.Huaweiyun.Bucket,
		Key:    key,
	}
	var output *obs.DeleteObjectOutput
	output, err = client.DeleteObject(input)
	if err != nil {
		return errors.Wrapf(err, "删除对象(%s)失败!, output: %v", key, output)
	}
	return nil
}
