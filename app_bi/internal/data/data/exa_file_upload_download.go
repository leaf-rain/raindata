package data

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ initDb = (*EntityExaFileUploadAndDownload)(nil)

type ExaFileUploadAndDownload struct {
	gorm.Model
	Name string `json:"name" gorm:"comment:文件名"` // 文件名
	Url  string `json:"url" gorm:"comment:文件地址"` // 文件地址
	Tag  string `json:"tag" gorm:"comment:文件标签"` // 文件标签
	Key  string `json:"key" gorm:"comment:编号"`   // 编号
}

func (ExaFileUploadAndDownload) TableName() string {
	return "exa_file_upload_and_downloads"
}

type EntityExaFileUploadAndDownload struct {
	*Data
	Model *ExaFileUploadAndDownload
}

func NewEntityExaFileUploadAndDownload(data *Data) *EntityExaFileUploadAndDownload {
	return &EntityExaFileUploadAndDownload{
		Data: data,
	}
}

func (i *EntityExaFileUploadAndDownload) MigrateTable(ctx context.Context) error {
	return i.SqlClient.AutoMigrate(&ExaFileUploadAndDownload{})
}

func (i *EntityExaFileUploadAndDownload) TableCreated(context.Context) bool {
	return i.SqlClient.Migrator().HasTable(&ExaFileUploadAndDownload{})
}

func (i *EntityExaFileUploadAndDownload) InitializeData(ctx context.Context) error {
	entities := []ExaFileUploadAndDownload{
		{Name: "10.png", Url: "https://qmplusimg.henrongyi.top/gvalogo.png", Tag: "png", Key: "158787308910.png"},
		{Name: "logo.png", Url: "https://qmplusimg.henrongyi.top/1576554439myAvatar.png", Tag: "png", Key: "1587973709logo.png"},
	}
	if err := i.SqlClient.Create(&entities).Error; err != nil {
		return errors.Wrap(err, ExaFileUploadAndDownload{}.TableName()+"表数据初始化失败!")
	}
	return nil
}

func (i *EntityExaFileUploadAndDownload) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	lookup := ExaFileUploadAndDownload{Name: "logo.png", Key: "1587973709logo.png"}
	if errors.Is(db.First(&lookup, &lookup).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
