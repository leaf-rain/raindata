package data

import (
	"context"
	"gorm.io/gorm"
)

// SysUserAuthority 是 sysUser 和 sysAuthority 的连接表
type SysUserAuthority struct {
	gorm.Model
	SysUserId      uint `gorm:"index;column:sys_user_id"`
	SysAuthorityId uint `gorm:"index;column:sys_authority_id"`
}

func (s *SysUserAuthority) TableName() string {
	return "sys_user_authority"
}

var _ initDb = (*EntitySysUserAuthority)(nil)

func NewEntitySysUserAuthority(data *Data) *EntitySysUserAuthority {
	return &EntitySysUserAuthority{
		Data: data,
	}
}

type EntitySysUserAuthority struct {
	*Data
	Model *SysUser
}

func (entity *EntitySysUserAuthority) InitializeData(ctx context.Context) error {
	return nil
}

func (entity *EntitySysUserAuthority) MigrateTable(ctx context.Context) error {
	return entity.SqlClient.AutoMigrate(&SysUserAuthority{})
}

func (entity *EntitySysUserAuthority) TableCreated(context.Context) bool {
	return entity.SqlClient.Migrator().HasTable(&SysUserAuthority{})
}

func (entity *EntitySysUserAuthority) InitializerName() string {
	return entity.Model.TableName()
}
