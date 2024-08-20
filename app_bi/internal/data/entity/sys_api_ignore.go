package entity

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type SysIgnoreApi struct {
	gorm.Model
	Path   string `json:"path" gorm:"comment:api路径"`             // api路径
	Method string `json:"method" gorm:"default:POST;comment:方法"` // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE
	Flag   bool   `json:"flag" gorm:"-"`                         // 是否忽略
}

func (SysIgnoreApi) TableName() string {
	return "sys_ignore_apis"
}

var _ initDb = (*EntitySysIgnoreApi)(nil)

type EntitySysIgnoreApi struct {
	data  *Data
	Model *SysIgnoreApi
}

func NewEntitySysIgnoreApi(data *Data) *EntitySysIgnoreApi {
	return &EntitySysIgnoreApi{
		data: data,
	}
}

func (i *EntitySysIgnoreApi) MigrateTable(ctx context.Context) error {
	return i.data.SqlClient.AutoMigrate(&SysIgnoreApi{})
}

func (i *EntitySysIgnoreApi) TableCreated(context.Context) bool {
	return i.data.SqlClient.Migrator().HasTable(&SysIgnoreApi{})
}

func (i *EntitySysIgnoreApi) InitializeData(ctx context.Context) (context.Context, error) {
	entities := []SysIgnoreApi{
		{Method: "GET", Path: "/swagger/*any"},
		{Method: "GET", Path: "/api/freshCasbin"},
		{Method: "GET", Path: "/uploads/file/*filepath"},
		{Method: "GET", Path: "/health"},
		{Method: "HEAD", Path: "/uploads/file/*filepath"},
		{Method: "POST", Path: "/autoCode/llmAuto"},
		{Method: "POST", Path: "/system/reloadSystem"},
		{Method: "POST", Path: "/base/login"},
		{Method: "POST", Path: "/base/captcha"},
		{Method: "POST", Path: "/init/initdb"},
		{Method: "POST", Path: "/init/checkdb"},
	}
	if err := i.data.SqlClient.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, SysIgnoreApi{}.TableName()+"表数据初始化失败!")
	}
	next := context.WithValue(ctx, i.Model.TableName(), entities)
	return next, nil
}

func (i *EntitySysIgnoreApi) DataInserted(ctx context.Context) bool {
	if errors.Is(i.data.SqlClient.Where("path = ? AND method = ?", "/swagger/*any", "GET").
		First(&SysIgnoreApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
