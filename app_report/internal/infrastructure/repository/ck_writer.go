package repository

import (
	"context"
	"github.com/leaf-rain/raindata/app_report/internal/domain/interface_repo"
	"github.com/leaf-rain/raindata/common/consts"
	"github.com/leaf-rain/raindata/common/rclickhouse"
	"go.uber.org/zap"
	"strconv"
	"sync"
)

var _ interface_repo.InterfaceWriterRepo = (*CKWriter)(nil)

func NewCkWriter(logger *zap.Logger, ckCluster *rclickhouse.ClickhouseCluster) *CKWriter {
	return &CKWriter{
		logger:        logger,
		ckCluster:     ckCluster,
		tableEventMap: make(map[int64]*rclickhouse.SinkerTable),
		tableLogMap:   make(map[int64]*rclickhouse.SinkerTable),
		lock:          &sync.Mutex{},
	}
}

type CKWriter struct {
	logger        *zap.Logger
	ckCluster     *rclickhouse.ClickhouseCluster
	tableEventMap map[int64]*rclickhouse.SinkerTable
	tableLogMap   map[int64]*rclickhouse.SinkerTable
	lock          *sync.Mutex
}

func (w CKWriter) WriterMsg(ctx context.Context, appid int64, event, msg string) error {
	tableInfo, err := w.loadTable(appid, event)
	if err != nil {
		return err
	}
	tableInfo.WriterMsg(rclickhouse.FetchSingle{
		Data: msg,
	})
	return nil
}

func (w CKWriter) loadTable(appid int64, event string) (*rclickhouse.SinkerTable, error) {
	var err error
	var m map[int64]*rclickhouse.SinkerTable
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
			engine = rclickhouse.EngineMergeTree
		} else {
			tableName = "Sinker_Event_" + strconv.FormatInt(appid, 10)
			engine = rclickhouse.EngineReplacingMergeTree
		}
		sinkerTable, err = rclickhouse.NewSinkerTable(context.TODO(), w.ckCluster, w.logger, &rclickhouse.SinkerTableConfig{
			Database:      w.ckCluster.GetDb(),
			TableName:     tableName,
			BufferSize:    10000,
			FlushInterval: 1,
			Parse:         "tcp",
			OrderByKey:    consts.KeyEventForMsg + "," + consts.KeyIdForMsg,
			ReplayKey:     consts.KeyVersionForMsg,
			BaseColumn: []rclickhouse.ColumnWithType{
				{
					Name: consts.KeyIdForMsg,
					Type: rclickhouse.WhichType("Int64"),
				},
				{
					Name: consts.KeyEventForMsg,
					Type: rclickhouse.WhichType("String"),
				},
			},
			Engine: engine,
		}, "fastjson")
		if err != nil {
			return nil, err
		}
		sinkerTable.Start()
		m[appid] = sinkerTable
	}
	return sinkerTable, nil
}
