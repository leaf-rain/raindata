package service

import (
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data"
	"github.com/leaf-rain/raindata/app_bi/internal/data/dto"
	"go.uber.org/zap"
)

type AutoCodeService struct {
	data *data.Data
	log  *zap.Logger
	conf *conf.Bootstrap
}

type Database interface {
	GetDB(businessDB string) (data []dto.Db, err error)
	GetTables(businessDB string, dbName string) (data []dto.Table, err error)
	GetColumn(businessDB string, tableName string, dbName string) (data []dto.Column, err error)
}
