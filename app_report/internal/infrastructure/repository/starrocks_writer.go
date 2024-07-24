package repository

import (
	"context"
	"database/sql"
	"github.com/leaf-rain/raindata/app_report/internal/domain/interface_repo"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure/consts"
	"github.com/leaf-rain/raindata/common/go_sql_driver"
	"go.uber.org/zap"
	"strconv"
	"sync"
)

var _ interface_repo.InterfaceWriterRepo = (*SRWriter)(nil)

func NewSRWriter(logger *zap.Logger, db *sql.DB, dbCfg *go_sql_driver.SqlConfig) *SRWriter {
	return &SRWriter{
		logger:        logger,
		db:            db,
		dbCfg:         dbCfg,
		tableEventMap: make(map[int64]*go_sql_driver.SinkerTable),
		tableLogMap:   make(map[int64]*go_sql_driver.SinkerTable),
		lock:          &sync.Mutex{},
	}
}

type SRWriter struct {
	logger        *zap.Logger
	db            *sql.DB
	dbCfg         *go_sql_driver.SqlConfig
	tableEventMap map[int64]*go_sql_driver.SinkerTable
	tableLogMap   map[int64]*go_sql_driver.SinkerTable
	lock          *sync.Mutex
}

func (w SRWriter) WriterMsg(ctx context.Context, appid int64, event, msg string) error {
	tableInfo, err := w.loadTable(appid, event)
	if err != nil {
		return err
	}
	tableInfo.WriterMsg(go_sql_driver.FetchSingle{
		Data: msg,
	})
	return nil
}

func (w SRWriter) loadTable(appid int64, event string) (*go_sql_driver.SinkerTable, error) {
	var err error
	var m map[int64]*go_sql_driver.SinkerTable
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
			engine = go_sql_driver.TableType_Duplicate
		} else {
			tableName = "Sinker_Event_" + strconv.FormatInt(appid, 10)
			engine = go_sql_driver.TableType_Primary
		}
		sinkerTable, err = go_sql_driver.NewSinkerTable(context.TODO(), w.db, w.logger, &go_sql_driver.SinkerTableConfig{
			TableName:     tableName,
			Database:      w.dbCfg.DB,
			FlushInterval: 1,
			BufferSize:    10000,
			BaseColumn: []go_sql_driver.ColumnWithType{
				{
					Name: consts.KeyIdForMsg,
					Type: go_sql_driver.WhichType(go_sql_driver.Bigint, false),
				},
				{
					Name: consts.KeyEventForMsg,
					Type: go_sql_driver.WhichType(go_sql_driver.String, false),
				},
				{
					Name: consts.KeyCreateTimeForMsg,
					Type: go_sql_driver.WhichType(go_sql_driver.Datetime, false),
				},
				{
					Name: consts.KeyOrder1ForMsg,
					Type: go_sql_driver.WhichType(go_sql_driver.Bigint, true),
				},
				{
					Name: consts.KeyOrder2ForMsg,
					Type: go_sql_driver.WhichType(go_sql_driver.Bigint, true),
				},
				{
					Name: consts.KeyOrder3ForMsg,
					Type: go_sql_driver.WhichType(go_sql_driver.Bigint, true),
				},
			},
			TableType:      engine,
			PrimaryKey:     consts.KeyIdForMsg,
			DistributedKey: consts.KeyCreateTimeForMsg,
			OrderByKey:     consts.KeyEventForMsg + "," + consts.KeyIdForMsg + "," + consts.KeyOrder1ForMsg + "," + consts.KeyOrder2ForMsg + "," + consts.KeyOrder3ForMsg,
			Username:       w.dbCfg.Username,
			Password:       w.dbCfg.Password,
			FeHost:         w.dbCfg.Host,
			FeHttpPort:     w.dbCfg.Port,
		}, "fastjson")
		if err != nil {
			return nil, err
		}
		sinkerTable.Start()
		m[appid] = sinkerTable
	}
	return sinkerTable, nil
}
