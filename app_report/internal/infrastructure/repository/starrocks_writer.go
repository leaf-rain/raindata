package repository

import (
	"context"
	"database/sql"
	"github.com/leaf-rain/raindata/app_report/internal/domain/interface_repo"
	"github.com/leaf-rain/raindata/common/consts"
	"github.com/leaf-rain/raindata/common/rsql"
	"go.uber.org/zap"
	"strconv"
	"sync"
)

var _ interface_repo.InterfaceWriterRepo = (*SRWriter)(nil)

func NewSRWriter(logger *zap.Logger, db *sql.DB, dbCfg *rsql.SqlConfig) *SRWriter {
	return &SRWriter{
		logger:        logger,
		db:            db,
		dbCfg:         dbCfg,
		tableEventMap: make(map[int64]*rsql.SinkerTable),
		tableLogMap:   make(map[int64]*rsql.SinkerTable),
		lock:          &sync.Mutex{},
	}
}

type SRWriter struct {
	logger        *zap.Logger
	db            *sql.DB
	dbCfg         *rsql.SqlConfig
	tableEventMap map[int64]*rsql.SinkerTable
	tableLogMap   map[int64]*rsql.SinkerTable
	lock          *sync.Mutex
}

func (w SRWriter) WriterMsg(ctx context.Context, appid int64, event, msg string) error {
	tableInfo, err := w.loadTable(appid, event)
	if err != nil {
		return err
	}
	tableInfo.WriterMsg(rsql.FetchSingle{
		Data: msg,
	})
	return nil
}

func (w SRWriter) loadTable(appid int64, event string) (*rsql.SinkerTable, error) {
	var err error
	var m map[int64]*rsql.SinkerTable
	if event == "log" {
		m = w.tableLogMap
	} else {
		m = w.tableEventMap
	}
	w.lock.Lock()
	defer w.lock.Unlock()
	sinkerTable, ok := m[appid]
	if !ok {
		tableName := ""
		var engine int
		if event == "log" {
			tableName = "Sinker_Log_" + strconv.FormatInt(appid, 10)
			engine = rsql.TableType_Duplicate
		} else {
			tableName = "Sinker_Event_" + strconv.FormatInt(appid, 10)
			engine = rsql.TableType_Primary
		}
		sinkerTable, err = rsql.NewSinkerTable(context.TODO(), w.db, w.logger, &rsql.SinkerTableConfig{
			TableName:     tableName,
			Database:      w.dbCfg.DB,
			FlushInterval: 1,
			BufferSize:    10000,
			BaseColumn: []rsql.ColumnWithType{
				{
					Name: consts.KeyIdForMsg,
					Type: rsql.WhichType(rsql.Bigint, false),
				},
				{
					Name: consts.KeyEventForMsg,
					Type: rsql.WhichType(rsql.String, false),
				},
				{
					Name: consts.KeyCreateTimeForMsg,
					Type: rsql.WhichType(rsql.Datetime, false),
				},
			},
			TableType:      engine,
			PrimaryKey:     consts.KeyEventForMsg + "," + consts.KeyIdForMsg,
			DistributedKey: consts.KeyCreateTimeForMsg,
			OrderByKey:     consts.KeyEventForMsg + "," + consts.KeyIdForMsg,
			Username:       w.dbCfg.Username,
			Password:       w.dbCfg.Password,
			FeHost:         w.dbCfg.Host,
			FeHttpPort:     w.dbCfg.HttpPort,
		}, "fastjson")
		if err != nil {
			return nil, err
		}
		sinkerTable.Start()
		m[appid] = sinkerTable
	}
	return sinkerTable, nil
}
