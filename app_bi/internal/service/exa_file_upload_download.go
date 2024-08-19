package service

import (
	"errors"
	"github.com/leaf-rain/raindata/app_bi/internal/data/entity"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"github.com/leaf-rain/raindata/app_bi/third_party/utils/upload"
	"mime/multipart"
	"strings"
)

//@function: Upload
//@description: 创建文件上传记录
//@param: file data.ExaFileUploadAndDownload
//@return: error

func (e *FileUploadAndDownloadService) Upload(file entity.ExaFileUploadAndDownload) error {
	return e.data.SqlClient.Create(&file).Error
}

//@function: FindFile
//@description: 查询文件记录
//@param: id uint
//@return: data.ExaFileUploadAndDownload, error

func (e *FileUploadAndDownloadService) FindFile(id uint) (entity.ExaFileUploadAndDownload, error) {
	var file entity.ExaFileUploadAndDownload
	err := e.data.SqlClient.Where("id = ?", id).First(&file).Error
	return file, err
}

//@function: DeleteFile
//@description: 删除文件记录
//@param: file data.ExaFileUploadAndDownload
//@return: err error

func (e *FileUploadAndDownloadService) DeleteFile(file entity.ExaFileUploadAndDownload) (err error) {
	var fileFromDb entity.ExaFileUploadAndDownload
	fileFromDb, err = e.FindFile(file.ID)
	if err != nil {
		return
	}
	oss := upload.NewOss(e.conf)
	if err = oss.DeleteFile(fileFromDb.Key); err != nil {
		return errors.New("文件删除失败")
	}
	err = e.data.SqlClient.Where("id = ?", file.ID).Unscoped().Delete(&file).Error
	return err
}

// EditFileName 编辑文件名或者备注
func (e *FileUploadAndDownloadService) EditFileName(file entity.ExaFileUploadAndDownload) (err error) {
	var fileFromDb entity.ExaFileUploadAndDownload
	return e.data.SqlClient.Where("id = ?", file.ID).First(&fileFromDb).Update("name", file.Name).Error
}

//@function: GetFileRecordInfoList
//@description: 分页获取数据
//@param: info dto.PageInfo
//@return: list interface{}, total int64, err error

func (e *FileUploadAndDownloadService) GetFileRecordInfoList(info rhttp.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	keyword := info.Keyword
	db := e.data.SqlClient.Model(&entity.ExaFileUploadAndDownload{})
	var fileLists []entity.ExaFileUploadAndDownload
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

func (e *FileUploadAndDownloadService) UploadFile(header *multipart.FileHeader, noSave string) (file entity.ExaFileUploadAndDownload, err error) {
	oss := upload.NewOss(e.conf)
	filePath, key, uploadErr := oss.UploadFile(header)
	if uploadErr != nil {
		panic(uploadErr)
	}
	s := strings.Split(header.Filename, ".")
	f := entity.ExaFileUploadAndDownload{
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
