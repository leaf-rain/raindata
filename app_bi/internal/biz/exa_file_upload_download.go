package biz

import (
	"errors"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"github.com/leaf-rain/raindata/app_bi/third_party/utils/upload"
	"mime/multipart"
	"strings"
)

//@function: Upload
//@description: 创建文件上传记录
//@param: file data.ExaFileUploadAndDownload
//@return: error

func (e *FileUploadAndDownload) Upload(file data.ExaFileUploadAndDownload) error {
	return e.data.SqlClient.Create(&file).Error
}

//@function: FindFile
//@description: 查询文件记录
//@param: id uint
//@return: data.ExaFileUploadAndDownload, error

func (e *FileUploadAndDownload) FindFile(id uint) (data.ExaFileUploadAndDownload, error) {
	var file data.ExaFileUploadAndDownload
	err := e.data.SqlClient.Where("id = ?", id).First(&file).Error
	return file, err
}

//@function: DeleteFile
//@description: 删除文件记录
//@param: file data.ExaFileUploadAndDownload
//@return: err error

func (e *FileUploadAndDownload) DeleteFile(file data.ExaFileUploadAndDownload) (err error) {
	var fileFromDb data.ExaFileUploadAndDownload
	fileFromDb, err = e.FindFile(file.ID)
	if err != nil {
		return
	}
	oss := upload.NewOss(e.data.Config, e.logger)
	if err = oss.DeleteFile(fileFromDb.Key); err != nil {
		return errors.New("文件删除失败")
	}
	err = e.data.SqlClient.Where("id = ?", file.ID).Unscoped().Delete(&file).Error
	return err
}

// EditFileName 编辑文件名或者备注
func (e *FileUploadAndDownload) EditFileName(file data.ExaFileUploadAndDownload) (err error) {
	var fileFromDb data.ExaFileUploadAndDownload
	return e.data.SqlClient.Where("id = ?", file.ID).First(&fileFromDb).Update("name", file.Name).Error
}

//@function: GetFileRecordInfoList
//@description: 分页获取数据
//@param: info rhttp.PageInfo
//@return: list interface{}, total int64, err error

func (e *FileUploadAndDownload) GetFileRecordInfoList(info rhttp.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	keyword := info.Keyword
	db := e.data.SqlClient.Model(&data.ExaFileUploadAndDownload{})
	var fileLists []data.ExaFileUploadAndDownload
	if len(keyword) > 0 {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&fileLists).Error
	return fileLists, total, err
}

//@function: UploadFile
//@description: 根据配置文件判断是文件上传到本地或者七牛云
//@param: header *multipart.FileHeader, noSave string
//@return: file data.ExaFileUploadAndDownload, err error

func (e *FileUploadAndDownload) UploadFile(header *multipart.FileHeader, noSave string) (file data.ExaFileUploadAndDownload, err error) {
	oss := upload.NewOss(e.data.Config, e.logger)
	filePath, key, uploadErr := oss.UploadFile(header)
	if uploadErr != nil {
		panic(uploadErr)
	}
	s := strings.Split(header.Filename, ".")
	f := data.ExaFileUploadAndDownload{
		Url:  filePath,
		Name: header.Filename,
		Tag:  s[len(s)-1],
		Key:  key,
	}
	if noSave == "0" {
		return f, e.Upload(f)
	}
	return f, nil
}
