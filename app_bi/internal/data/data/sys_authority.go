package data

import (
	"context"
	"github.com/leaf-rain/raindata/app_bi/third_party/utils"
	"github.com/pkg/errors"
	"time"
)

type SysAuthority struct {
	CreatedAt       time.Time       // 创建时间
	UpdatedAt       time.Time       // 更新时间
	DeletedAt       *time.Time      `sql:"index"`
	AuthorityId     uint            `json:"authorityId" gorm:"not null;unique;primary_key;comment:角色ID;size:90"` // 角色ID
	AuthorityName   string          `json:"authorityName" gorm:"comment:角色名"`                                    // 角色名
	ParentId        *uint           `json:"parentId" gorm:"comment:父角色ID"`                                       // 父角色ID
	DefaultRouter   string          `json:"defaultRouter" gorm:"comment:默认菜单;default:dashboard"`                 // 默认菜单(默认dashboard)
	DataAuthorityId []*SysAuthority `json:"dataAuthorityId" gorm:"-"`
	Children        []SysAuthority  `json:"children" gorm:"-"`
	Users           []SysUser       `json:"-" gorm:"-"`
}

func (SysAuthority) TableName() string {
	return "sys_authorities"
}

var _ initDb = (*EntitySysAuthority)(nil)

type EntitySysAuthority struct {
	*Data
	Model *SysAuthority
}

func NewEntitySysAuthority(data *Data) *EntitySysAuthority {
	return &EntitySysAuthority{
		Data: data,
	}
}

func (i *EntitySysAuthority) MigrateTable(ctx context.Context) error {
	return i.SqlClient.AutoMigrate(&SysAuthority{})
}

func (i *EntitySysAuthority) TableCreated(context.Context) bool {
	return i.SqlClient.Migrator().HasTable(&SysAuthority{})
}

func (i *EntitySysAuthority) InitializeData(ctx context.Context) (err error) {
	entities := []SysAuthority{
		{AuthorityId: 1, AuthorityName: "系统管理员", ParentId: utils.Pointer[uint](0), DefaultRouter: "dashboard"},
	}
	if err = i.SqlClient.Create(&entities).Error; err != nil {
		return errors.Wrapf(err, "%s表数据初始化失败!", SysAuthority{}.TableName())
	}
	return nil
}
