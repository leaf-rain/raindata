package clickhouse_sqlx

import (
	"github.com/leaf-rain/raindata/common/ecode"
	"github.com/leaf-rain/raindata/common/parser"
	"sync"
)

var (
	ErrTblNotExist    = ecode.Newf("table doesn't exist")
	selectSQLTemplate = `select name, type, default_kind from system.columns where database = '%s' and table = '%s'`

	// https://github.com/ClickHouse/ClickHouse/issues/24036
	// src/Common/ErrorCodes.cpp
	// src/Storages/MergeTree/ReplicatedMergeTreeBlockOutputStream.cpp
	// ZooKeeper issues(https://issues.apache.org/jira/browse/ZOOKEEPER-4410) can cause ClickHouse exeception: "Code": 999, "Message": "Cannot allocate block number..."
	// CKServer too many parts possibly reason: https://github.com/ClickHouse/ClickHouse/issues/6720#issuecomment-526045768
	// zooKeeper Connection loss issue: https://cwiki.apache.org/confluence/display/ZOOKEEPER/FAQ#:~:text=How%20should%20I%20handle%20the%20CONNECTION_LOSS%20error%3F
	// zooKeeper Session expired issue: https://cwiki.apache.org/confluence/display/ZOOKEEPER/FAQ#:~:text=How%20should%20I%20handle%20SESSION_EXPIRED%3F
	// TOO_MANY_SIMULTANEOUS_QUERIES, NO_ZOOKEEPER, TABLE_IS_READ_ONLY, TOO_MANY_PARTS, UNKNOWN_STATUS_OF_INSERT, KEEPER_EXCEPTION, POCO_EXCEPTION
	replicaSpecificErrorCodes = []int32{202, 225, 242, 252, 319, 999, 1000}
	wrSeriesQuota             = 16384

	SeriesQuotas sync.Map
)

// ClickHouse is an output service consumers from kafka messages
type ClickHouse struct {
	Dims      []*parser.ColumnWithType
	NumDims   int
	IdxSerID  int
	NameKey   string
	cfg       *ClickhouseConfig
	TableName string
	dbName    string

	prepareSQL string
	promSerSQL string
	seriesTbl  string

	distMetricTbls []string
	distSeriesTbls []string
	DimSerID       string
	DimMgmtID      string

	seriesQuota *parser.SeriesQuota

	numFlying int32
	mux       sync.Mutex
	taskDone  *sync.Cond
}

type DistTblInfo struct {
	name    string
	cluster string
}
