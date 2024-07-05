package clickhouse_sqlx

import (
	"context"
	"github.com/leaf-rain/raindata/common/parser"
	"go.uber.org/zap"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Decoder decodes the contents of b into v.
// It's primarily used for decoding contents of a file into a map[string]any.
type Decoder interface {
	Decode(b []byte, v map[string]any) error
}

type SinkerTableConfig struct {
	TableName     string           `json:"table_name,omitempty"`
	Database      string           `json:"database,omitempty"`
	FlushInterval int              `json:"flush_interval,omitempty"`
	BufferSize    int              `json:"buffer_size,omitempty"`
	TimeZone      string           `json:"time_zone,omitempty"`
	TimeUnit      float64          `json:"time_unit,omitempty"`
	Parse         string           `json:"parse,omitempty"` // 选择解析方式
	BaseColumn    []ColumnWithType `json:"base_column,omitempty"`
	ReplayKey     string           `json:"replay_key"`
	OrderByKey    string           `json:"order_by_key"`
}

// SinkerTable object maintains number of task for each partition
type SinkerTable struct {
	logger         *zap.Logger
	clusterConn    *ClickhouseCluster
	ctx            context.Context
	cancel         context.CancelFunc
	fetchCH        chan Fetch
	exitCh         chan struct{}
	processWg      sync.WaitGroup
	mux            sync.Mutex
	commitDone     *sync.Cond
	state          atomic.Uint32
	parser         parser.Parser
	prepareSQLHttp string
	prepareSQLTcp  string
	fieldMap       *sync.Map // key:columnName, value:*ColumnWithType
	// table相关
	sinkerTableConfig *SinkerTableConfig
	// todo:wal功能
}

// NewSinkerTable get an instance of sinker with the task list
func NewSinkerTable(ctx context.Context, cc *ClickhouseCluster, logger *zap.Logger, sinkerTableConfig *SinkerTableConfig, parserTy string) (*SinkerTable, error) {
	s := &SinkerTable{
		logger:            logger,
		clusterConn:       cc,
		ctx:               ctx,
		fetchCH:           make(chan Fetch, 1000),
		exitCh:            make(chan struct{}),
		sinkerTableConfig: sinkerTableConfig,
		fieldMap:          new(sync.Map),
	}
	s.state.Store(StateStopped)
	s.commitDone = sync.NewCond(&s.mux)
	s.parser, _ = parser.NewParse(parser.ParserConfig{
		Ty:     parserTy,
		Logger: logger,
	})
	// 存储默认字段信息
	for i := range sinkerTableConfig.BaseColumn {
		s.fieldMap.Store(sinkerTableConfig.BaseColumn[i].Name, &sinkerTableConfig.BaseColumn[i])
	}
	err := s.flushFields()
	if err != nil {
		return nil, err
	}
	err = s.wal()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (c *SinkerTable) wal() error {
	// todo:wal恢复未处理的数据
	return nil
}

func (c *SinkerTable) flushFields() error {
	conn := c.clusterConn.GetShardConn(time.Now().Unix())
	ck, _, err := conn.NextGoodReplica(0)
	if err != nil {
		return err
	}
	var query = getTableColumns(c.sinkerTableConfig.Database, c.sinkerTableConfig.TableName)
	var rs *Rows
	rs, err = ck.Query(query)
	if err == nil {
		var columns = make(map[string]*ColumnWithType)
		defer rs.Close()
		var name, typ, defaultKind string
		for rs.Next() {
			if err = rs.Scan(&name, &typ, &defaultKind); err != nil {
				return err
			}
			ct := &ColumnWithType{Name: name, Type: WhichType(typ)}
			c.fieldMap.Store(name, ct)
			columns[name] = ct
		}
		if len(columns) == 0 {
			query = createTable(c.sinkerTableConfig.Database, c.sinkerTableConfig.TableName, c.sinkerTableConfig.ReplayKey, c.sinkerTableConfig.OrderByKey, c.fieldMap)
			err = ck.Exec(query)
			if err != nil {
				return err
			}
			return c.flushFields()
		}
		var addColumns []string
		var allColumns []string
		c.fieldMap.Range(func(key, value any) bool {
			if _, ok := columns[key.(string)]; !ok {
				addColumns = append(addColumns, addTableColumns(c.sinkerTableConfig.Database, c.sinkerTableConfig.TableName, key.(string), value.(*ColumnWithType).Type.ToString()))
			}
			allColumns = append(allColumns, key.(string))
			return true
		})
		sort.Slice(allColumns, func(i, j int) bool {
			return allColumns[i] < allColumns[j]
		})
		if len(addColumns) > 0 {
			for i := range addColumns {
				err = ck.Exec(addColumns[i])
				if err != nil {
					c.logger.Error(c.sinkerTableConfig.TableName+" add table columns error", zap.Error(err), zap.String("sql", addColumns[i]))
					return err
				}
			}
		}
		// 拼接入库sql
		columnsStr := strings.Join(allColumns, ",")
		var params = make([]string, len(allColumns))
		for i := range params {
			params[i] = "?"
		}
		c.prepareSQLHttp = "INSERT INTO " + c.sinkerTableConfig.Database + "." + c.sinkerTableConfig.TableName + " (" + columnsStr + ") " +
			"VALUES (" + strings.Join(params, ",") + ")"
		c.prepareSQLTcp = "INSERT INTO " + c.sinkerTableConfig.Database + "." + c.sinkerTableConfig.TableName + " (" + columnsStr + ") "
	}
	return err
}

func (c *SinkerTable) start() {
	if c.state.Load() == StateRunning {
		return
	}
	c.state.Store(StateRunning)
	go c.processFetch()
}

func (c *SinkerTable) stop() {
	if c.state.Load() == StateStopped {
		return
	}
	c.state.Store(StateStopped)
	c.exitCh <- struct{}{}
	c.processWg.Wait()
}

func (c *SinkerTable) restart() {
	c.stop()
	c.start()
}

func (c *SinkerTable) Fetch2Row(fetch Fetch) ([][]any, error) {
	var result [][]any
	for _, item := range fetch.GetData() {
		val, err := c.metric2Row([]byte(item))
		if err != nil {
			c.logger.Error(c.sinkerTableConfig.TableName+" parse msg error", zap.Error(err), zap.String("msg", item))
			return nil, err
		}
		result = append(result, val)
	}
	return result, nil
}

func (c *SinkerTable) metric2Row(msg []byte) ([]any, error) {
	metric, err := c.parser.Parse(msg)
	if err != nil {
		return nil, err
	}
	var rows []any
	var allColumns []*ColumnWithType
	c.fieldMap.Range(func(key, value any) bool {
		allColumns = append(allColumns, value.(*ColumnWithType))
		return true
	})
	sort.Slice(allColumns, func(i, j int) bool {
		return allColumns[i].Name < allColumns[j].Name
	})
	for _, item := range allColumns {
		val := GetValueByType(metric, item)
		rows = append(rows, val)
	}
	return rows, err
}

func (c *SinkerTable) processFetch() {
	c.processWg.Add(1)
	defer c.processWg.Done()
	var bufLength int
	var data FetchArray
	flushFn := func(data Fetch) {
		bufLength = len(data.GetData())
		c.logger.Info(c.sinkerTableConfig.TableName+" flush msg.", zap.Int("bufLength", bufLength))
		if bufLength == 0 {
			return
		}
		conn := c.clusterConn.GetShardConn(time.Now().Unix())
		ck, _, err := conn.NextGoodReplica(0)
		if err != nil {
			c.logger.Error(c.sinkerTableConfig.TableName+" conn.NextGoodReplica error", zap.Error(err))
			return
		}
		var parserSql string
		if c.sinkerTableConfig.Parse == "http" {
			parserSql = c.prepareSQLHttp
		} else {
			parserSql = c.prepareSQLTcp
		}
		var ckData [][]any
		ckData, err = c.Fetch2Row(data)
		if err != nil {
			c.logger.Error(c.sinkerTableConfig.TableName+" Fetch2Row error", zap.Error(err))
			return
		}
		err = ck.Write(parserSql, ckData)
		if err != nil {
			c.logger.Error(c.sinkerTableConfig.TableName+" conn.Write error", zap.Error(err), zap.Int("bufLength", bufLength))
		} else {
			c.logger.Info(c.sinkerTableConfig.TableName+" flush msg success.", zap.Int("bufLength", bufLength))
			// todo: wal标记删除状态
		}
	}

	ticker := time.NewTicker(time.Duration(c.sinkerTableConfig.FlushInterval) * time.Second)
	defer ticker.Stop()
	var err error
	var parse parser.Metric
	var newKeys map[string]string
	for {
		select {
		case fetch := <-c.fetchCH:
			if c.state.Load() == StateStopped {
				continue
			}
			// todo:wal
			tmpData := fetch.GetData()
			for i := range tmpData {
				parse, err = c.parser.Parse([]byte(tmpData[i]))
				if err != nil {
					c.logger.Error(c.sinkerTableConfig.TableName+" parse error", zap.Error(err), zap.String("data", tmpData[i]))
					continue
				}
				newKeys = parse.GetNewKeys(c.fieldMap)
				if len(newKeys) > 0 {
					for k, v := range newKeys {
						c.fieldMap.Store(k, &ColumnWithType{Name: k, Type: WhichType(v)})
					}
					err = c.flushFields()
					if err != nil {
						c.logger.Error(c.sinkerTableConfig.TableName+" flush fields error", zap.Error(err))
						continue
					}
				}
			}
			data.Data = append(data.Data, fetch.GetData()...)
			data.Callback = append(data.Callback, fetch.GetCallback()...)
			bufLength = len(data.Data)
			if bufLength > c.sinkerTableConfig.BufferSize {
				tmp := data.Copy()
				go flushFn(tmp)
				data = FetchArray{}
				ticker.Reset(time.Duration(c.sinkerTableConfig.FlushInterval) * time.Second)
			}
		case <-ticker.C:
			tmp := data.Copy()
			go flushFn(tmp)
			data = FetchArray{}
		case <-c.ctx.Done():
			c.logger.Info("stopped processing loop", zap.String("table", c.sinkerTableConfig.Database+"."+c.sinkerTableConfig.TableName))
			return
		}
	}
}
