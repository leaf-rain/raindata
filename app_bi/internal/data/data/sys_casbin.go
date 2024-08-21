package data

import (
	"context"
	"errors"
	adapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

var _ initDb = (*EntitySysCasbin)(nil)

type EntitySysCasbin struct {
	*Data
	Model *adapter.CasbinRule
}

func NewEntitySysCasbin(data *Data) *EntitySysCasbin {
	return &EntitySysCasbin{
		Data: data,
	}
}

func (i *EntitySysCasbin) MigrateTable(ctx context.Context) error {
	return i.SqlClient.AutoMigrate(&adapter.CasbinRule{})
}

func (i *EntitySysCasbin) TableCreated(context.Context) bool {
	return i.SqlClient.Migrator().HasTable(&adapter.CasbinRule{})
}

func (i *EntitySysCasbin) InitializerName() string {
	var entity adapter.CasbinRule
	return entity.TableName()
}

func (i *EntitySysCasbin) InitializeData(ctx context.Context) error {
	return nil
}

func (i *EntitySysCasbin) DataInserted(ctx context.Context) bool {
	if errors.Is(i.SqlClient.Where(adapter.CasbinRule{Ptype: "p", V0: "9528", V1: "/user/getUserInfo", V2: "GET"}).
		First(&adapter.CasbinRule{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
