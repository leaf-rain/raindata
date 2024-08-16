package upload

import (
	"errors"
	"fmt"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/zap"
)

type AwsS3 struct {
	conf   *conf.Bootstrap
	logger *zap.Logger
}

//@object: u *AwsS3
//@function: UploadFile
//@description: Upload file to Aws S3 using aws-sdk-go. See https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-basic-bucket-operations.html#s3-examples-bucket-ops-upload-file-to-bucket
//@param: file *multipart.FileHeader
//@return: string, string, error

func (u *AwsS3) UploadFile(file *multipart.FileHeader) (string, string, error) {
	session := u.newSession()
	uploader := s3manager.NewUploader(session)

	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename)
	filename := u.conf.Oss.AwsS3.PathPrefix + "/" + fileKey
	f, openError := file.Open()
	if openError != nil {
		u.logger.Error("function file.Open() failed", zap.Any("err", openError.Error()))
		return "", "", errors.New("function file.Open() failed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(u.conf.Oss.AwsS3.Bucket),
		Key:    aws.String(filename),
		Body:   f,
	})
	if err != nil {
		u.logger.Error("function uploader.Upload() failed", zap.Any("err", err.Error()))
		return "", "", err
	}

	return u.conf.Oss.AwsS3.BaseURL + "/" + filename, fileKey, nil
}

//@object: u *AwsS3
//@function: DeleteFile
//@description: Delete file from Aws S3 using aws-sdk-go. See https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-basic-bucket-operations.html#s3-examples-bucket-ops-delete-bucket-item
//@param: file *multipart.FileHeader
//@return: string, string, error

func (u *AwsS3) DeleteFile(key string) error {
	session := u.newSession()
	svc := s3.New(session)
	filename := u.conf.Oss.AwsS3.PathPrefix + "/" + key
	bucket := u.conf.Oss.AwsS3.Bucket

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		u.logger.Error("function svc.DeleteObject() failed", zap.Any("err", err.Error()))
		return errors.New("function svc.DeleteObject() failed, err:" + err.Error())
	}

	_ = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	return nil
}

// newSession Create S3 session
func (u *AwsS3) newSession() *session.Session {
	sess, _ := session.NewSession(&aws.Config{
		Region:           aws.String(u.conf.Oss.AwsS3.Region),
		Endpoint:         aws.String(u.conf.Oss.AwsS3.Endpoint), //minio在这里设置地址,可以兼容
		S3ForcePathStyle: aws.Bool(u.conf.Oss.AwsS3.S3ForcePathStyle),
		DisableSSL:       aws.Bool(u.conf.Oss.AwsS3.DisableSSL),
		Credentials: credentials.NewStaticCredentials(
			u.conf.Oss.AwsS3.SecretID,
			u.conf.Oss.AwsS3.SecretKey,
			"",
		),
	})
	return sess
}
