package service

import (
	"github.com/leaf-rain/raindata/app_bi/internal/biz"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
)

type FileUploadAndDownloadService struct {
	*Service
}

func NewFileUploadAndDownloadService(service *Service) *FileUploadAndDownloadService {
	return &FileUploadAndDownloadService{
		service,
	}
}

//@function: FindOrCreateFile
//@description: 上传文件时检测当前文件属性，如果没有文件则创建，有则返回文件的当前切片
//@param: fileMd5 string, fileName string, chunkTotal int
//@return: file data.ExaFile, err error

func (e *FileUploadAndDownloadService) FindOrCreateFile(fileMd5 string, fileName string, chunkTotal int) (file data.ExaFile, err error) {
	b := biz.NewFileUploadAndDownload(e.biz)
	return b.FindOrCreateFile(fileMd5, fileName, chunkTotal)
}

//@function: CreateFileChunk
//@description: 创建文件切片记录
//@param: id uint, fileChunkPath string, fileChunkNumber int
//@return: error

func (e *FileUploadAndDownloadService) CreateFileChunk(id uint, fileChunkPath string, fileChunkNumber int) error {
	b := biz.NewFileUploadAndDownload(e.biz)
	return b.CreateFileChunk(id, fileChunkPath, fileChunkNumber)
}

//@function: DeleteFileChunk
//@description: 删除文件切片记录
//@param: fileMd5 string, fileName string, filePath string
//@return: error

func (e *FileUploadAndDownloadService) DeleteFileChunk(fileMd5 string, filePath string) error {
	b := biz.NewFileUploadAndDownload(e.biz)
	return b.DeleteFileChunk(fileMd5, filePath)
}
