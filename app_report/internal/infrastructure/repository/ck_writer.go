package repository

import (
	"context"
	"github.com/leaf-rain/raindata/app_report/internal/domain/interface_repo"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure/consts"
	"github.com/leaf-rain/raindata/common/clickhouse_sqlx"
	"go.uber.org/zap"
	"strconv"
	"sync"
)

var _ interface_repo.InterfaceWriterRepo = (*CKWriter)(nil)

func NewCkWriter(logger *zap.Logger, ckCluster *clickhouse_sqlx.ClickhouseCluster) *CKWriter {
	return &CKWriter{
		logger:        logger,
		ckCluster:     ckCluster,
		tableEventMap: make(map[int64]*clickhouse_sqlx.SinkerTable),
		tableLogMap:   make(map[int64]*clickhouse_sqlx.SinkerTable),
		lock:          &sync.Mutex{},
	}
}

type CKWriter struct {
	logger        *zap.Logger
	ckCluster     *clickhouse_sqlx.ClickhouseCluster
	tableEventMap map[int64]*clickhouse_sqlx.SinkerTable
	tableLogMap   map[int64]*clickhouse_sqlx.SinkerTable
	lock          *sync.Mutex
}

func (w CKWriter) WriterMsg(ctx context.Context, appid int64, event, msg string) error {
	tableInfo, err := w.loadTable(appid, event)
	if err != nil {
		return err
	}
	tableInfo.WriterMsg(clickhouse_sqlx.FetchSingle{
		Data: msg,
	})
	return nil
}

func (w CKWriter) loadTable(appid int64, event string) (*clickhouse_sqlx.SinkerTable, error) {
	var err error
	var m map[int64]*clickhouse_sqlx.SinkerTable
	if event == "log" {
		m = w.tableLogMap
	} else {
		m = w.tableEventMap
	}
	sinkerTable, ok := m[appid]
	if !ok {
		w.lock.Lock()
		defer w.lock.Unlock()
		tableName := ""
		if event == "log" {
			tableName = "Sinker_Log_" + strconv.FormatInt(appid, 10)
		} else {
			tableName = "Sinker_Event_" + strconv.FormatInt(appid, 10)
		}
		sinkerTable, err = clickhouse_sqlx.NewSinkerTable(context.TODO(), w.ckCluster, w.logger, &clickhouse_sqlx.SinkerTableConfig{
			Database:      w.ckCluster.GetDb(),
			TableName:     tableName,
			BufferSize:    10000,
			FlushInterval: 1,
			Parse:         "tcp",
			OrderByKey:    consts.KeyEventForMsg,
			ReplayKey:     consts.KeyIdForMsg,
			BaseColumn: []clickhouse_sqlx.ColumnWithType{
				{
					Name: consts.KeyIdForMsg,
					Type: clickhouse_sqlx.WhichType("Int64"),
				},
				{
					Name: consts.KeyEventForMsg,
					Type: clickhouse_sqlx.WhichType("String"),
				},
			},
		}, "fastjson")
		if err != nil {
			return nil, err
		}
		sinkerTable.Start()
	}
	return sinkerTable, nil
}
