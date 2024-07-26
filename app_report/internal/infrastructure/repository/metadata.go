package repository

import (
	"database/sql"
	"github.com/leaf-rain/raindata/app_report/internal/domain/interface_repo"
	"github.com/leaf-rain/raindata/common/rsql"
	"go.uber.org/zap"
	"sync"
)

// todo:实现将字段改为元数据对应关系维护至数据库，比如同一事件fields修改为s1,i1,i2,f1等，可减少表宽度，从而优化存储以及查询。
var _ interface_repo.InterfaceMetadataRepo = (*MetadataWriter)(nil)

func NewMetadataWriter(logger *zap.Logger, db *sql.DB, dbCfg *rsql.SqlConfig) *MetadataWriter {
	return &MetadataWriter{
		logger:        logger,
		db:            db,
		dbCfg:         dbCfg,
		tableEventMap: make(map[int64]*rsql.SinkerTable),
		tableLogMap:   make(map[int64]*rsql.SinkerTable),
		lock:          &sync.Mutex{},
	}
}

type MetadataWriter struct {
	logger        *zap.Logger
	db            *sql.DB
	dbCfg         *rsql.SqlConfig
	tableEventMap map[int64]*rsql.SinkerTable
	tableLogMap   map[int64]*rsql.SinkerTable
	lock          *sync.Mutex
}

func (m MetadataWriter) MetadataPut(appid int64, event string, keys []rsql.FieldInfo) (map[string]string, error) {
	//TODO implement me
	panic("implement me")
}
