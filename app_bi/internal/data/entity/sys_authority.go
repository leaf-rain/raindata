package entity

import (
	"context"
	"github.com/leaf-rain/raindata/app_bi/third_party/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
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
	data  *Data
	Model *SysIgnoreApi
}

func NewEntitySysAuthority(data *Data) *EntitySysAuthority {
	return &EntitySysAuthority{
		data: data,
	}
}

func (i *EntitySysAuthority) MigrateTable(ctx context.Context) error {
	return i.data.SqlClient.AutoMigrate(&SysAuthority{})
}

func (i *EntitySysAuthority) TableCreated(context.Context) bool {
	return i.data.SqlClient.Migrator().HasTable(&SysAuthority{})
}

func (i *EntitySysAuthority) InitializeData(ctx context.Context) (context.Context, error) {
	entities := []SysAuthority{
		{AuthorityId: 888, AuthorityName: "普通用户", ParentId: utils.Pointer[uint](0), DefaultRouter: "dashboard"},
		{AuthorityId: 9528, AuthorityName: "测试角色", ParentId: utils.Pointer[uint](0), DefaultRouter: "dashboard"},
		{AuthorityId: 8881, AuthorityName: "普通用户子角色", ParentId: utils.Pointer[uint](888), DefaultRouter: "dashboard"},
	}

	if err := i.data.SqlClient.Create(&entities).Error; err != nil {
		return ctx, errors.Wrapf(err, "%s表数据初始化失败!", SysAuthority{}.TableName())
	}
	// data authority
	if err := i.data.SqlClient.Model(&entities[0]).Association("DataAuthorityId").Replace(
		[]*SysAuthority{
			{AuthorityId: 888},
			{AuthorityId: 9528},
			{AuthorityId: 8881},
		}); err != nil {
		return ctx, errors.Wrapf(err, "%s表数据初始化失败!",
			i.data.SqlClient.Model(&entities[0]).Association("DataAuthorityId").Relationship.JoinTable.Name)
	}
	if err := i.data.SqlClient.Model(&entities[1]).Association("DataAuthorityId").Replace(
		[]*SysAuthority{
			{AuthorityId: 9528},
			{AuthorityId: 8881},
		}); err != nil {
		return ctx, errors.Wrapf(err, "%s表数据初始化失败!",
			i.data.SqlClient.Model(&entities[1]).Association("DataAuthorityId").Relationship.JoinTable.Name)
	}

	next := context.WithValue(ctx, i.Model.TableName(), entities)
	return next, nil
}

func (i *EntitySysAuthority) DataInserted(ctx context.Context) bool {
	if errors.Is(i.data.SqlClient.Where("authority_id = ?", "8881").
		First(&SysAuthority{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
