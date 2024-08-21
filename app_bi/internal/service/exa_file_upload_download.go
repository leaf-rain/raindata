package service

import (
	"github.com/leaf-rain/raindata/app_bi/internal/biz"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"github.com/leaf-rain/raindata/app_bi/third_party/rhttp"
	"mime/multipart"
)

//@function: Upload
//@description: 创建文件上传记录
//@param: file data.ExaFileUploadAndDownload
//@return: error

func (e *FileUploadAndDownloadService) Upload(file data.ExaFileUploadAndDownload) error {
	b := biz.NewFileUploadAndDownload(e.biz)
	return b.Upload(file)
}

//@function: FindFile
//@description: 查询文件记录
//@param: id uint
//@return: data.ExaFileUploadAndDownload, error

func (e *FileUploadAndDownloadService) FindFile(id uint) (data.ExaFileUploadAndDownload, error) {
	b := biz.NewFileUploadAndDownload(e.biz)
	return b.FindFile(id)
}

//@function: DeleteFile
//@description: 删除文件记录
//@param: file data.ExaFileUploadAndDownload
//@return: err error

func (e *FileUploadAndDownloadService) DeleteFile(file data.ExaFileUploadAndDownload) (err error) {
	b := biz.NewFileUploadAndDownload(e.biz)
	return b.DeleteFile(file)
}

// EditFileName 编辑文件名或者备注
func (e *FileUploadAndDownloadService) EditFileName(file data.ExaFileUploadAndDownload) (err error) {
	b := biz.NewFileUploadAndDownload(e.biz)
	return b.EditFileName(file)
}

//@function: GetFileRecordInfoList
//@description: 分页获取数据
//@param: info dto.PageInfo
//@return: list interface{}, total int64, err error

func (e *FileUploadAndDownloadService) GetFileRecordInfoList(info rhttp.PageInfo) (list interface{}, total int64, err error) {
	b := biz.NewFileUploadAndDownload(e.biz)
	return b.GetFileRecordInfoList(info)
}

//@function: UploadFile
//@description: 根据配置文件判断是文件上传到本地或者七牛云
//@param: header *multipart.FileHeader, noSave string
//@return: file data.ExaFileUploadAndDownload, err error

func (e *FileUploadAndDownloadService) UploadFile(header *multipart.FileHeader, noSave string) (file data.ExaFileUploadAndDownload, err error) {
	b := biz.NewFileUploadAndDownload(e.biz)
	return b.UploadFile(header, noSave)
}
