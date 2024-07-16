package go_sql_driver

import (
	"context"
	"database/sql"
	"github.com/leaf-rain/raindata/common/parser"
	"go.uber.org/zap"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

const (
	StateRunning uint32 = 0
	StateStopped uint32 = 1
)

// Decoder decodes the contents of b into v.
// It's primarily used for decoding contents of a file into a map[string]any.
type Decoder interface {
	Decode(b []byte, v map[string]any) error
}

type SinkerTableConfig struct {
	TableName      string           `json:"table_name,omitempty"`
	Database       string           `json:"database,omitempty"`
	FlushInterval  int              `json:"flush_interval,omitempty"`
	BufferSize     int              `json:"buffer_size,omitempty"`
	BaseColumn     []ColumnWithType `json:"base_column,omitempty"`
	TableType      int              `json:"table_type"`
	PrimaryKey     string           `json:"primary_key"`
	DistributedKey string           `json:"distributed_key"`
	OrderByKey     string           `json:"order_by_key"`
}

// SinkerTable object maintains number of task for each partition
type SinkerTable struct {
	logger     *zap.Logger
	db         *sql.DB
	ctx        context.Context
	cancel     context.CancelFunc
	fetchCH    chan Fetch
	exitCh     chan struct{}
	processWg  sync.WaitGroup
	mux        sync.Mutex
	commitDone *sync.Cond
	state      atomic.Uint32
	parser     parser.Parser
	fieldMap   *sync.Map // key:columnName, value:*ColumnWithType
	// table相关
	sinkerTableConfig *SinkerTableConfig
	// todo:wal功能
}

// NewSinkerTable get an instance of sinker with the task list
func NewSinkerTable(ctx context.Context, db *sql.DB, logger *zap.Logger, sinkerTableConfig *SinkerTableConfig, parserTy string) (*SinkerTable, error) {
	s := &SinkerTable{
		logger:            logger,
		db:                db,
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
	if s.sinkerTableConfig.FlushInterval <= 0 {
		s.sinkerTableConfig.FlushInterval = 1
	}
	return s, nil
}

func (c *SinkerTable) wal() error {
	// todo:wal恢复未处理的数据
	return nil
}

func (c *SinkerTable) flushFields() error {
	var query = getTableColumns(c.sinkerTableConfig.Database, c.sinkerTableConfig.TableName)
	var err error
	var rs *sql.Rows
	rs, err = c.db.Query(query)
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
			query = createTable(c.sinkerTableConfig.TableType, c.sinkerTableConfig.Database, c.sinkerTableConfig.TableName, c.sinkerTableConfig.PrimaryKey, c.sinkerTableConfig.DistributedKey, c.sinkerTableConfig.OrderByKey, c.fieldMap)
			_, err = c.db.Exec(query)
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
				_, err = c.db.Exec(addColumns[i])
				if err != nil {
					c.logger.Error(c.sinkerTableConfig.TableName+" add table columns error", zap.Error(err), zap.String("sql", addColumns[i]))
					return err
				}
			}
		}
	}
	return err
}

func (c *SinkerTable) Start() {
	if c.state.Load() == StateRunning {
		return
	}
	c.state.Store(StateRunning)
	go c.processFetch()
}

func (c *SinkerTable) Stop() {
	if c.state.Load() == StateStopped {
		return
	}
	c.state.Store(StateStopped)
	c.exitCh <- struct{}{}
	c.processWg.Wait()
}

func (c *SinkerTable) Restart() {
	c.Stop()
	c.Start()
}

func (c *SinkerTable) WriterMsg(msg FetchSingle) {
	c.fetchCH <- msg
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
		if item.Name == "_create_time" {
			rows = append(rows, time.Now().Unix())
		} else {
			val := GetValueByType(metric, item)
			rows = append(rows, val)
		}
	}
	return rows, err
}

func (c *SinkerTable) processFetch() {
	c.processWg.Add(1)
	defer c.processWg.Done()
	var data FetchArray
	var err error
	flushFn := func(data Fetch) {
		bufLength := len(data.GetData())
		c.logger.Info(c.sinkerTableConfig.TableName+" flush msg.", zap.Int("bufLength", bufLength))
		if bufLength == 0 {
			return
		}
		var ckData [][]any
		ckData, err = c.Fetch2Row(data)
		if err != nil {
			c.logger.Error(c.sinkerTableConfig.TableName+" Fetch2Row error", zap.Error(err))
			return
		}
		// todo:写文件
		err = c.db.Write(c.prepareSQL, ckData)
		if err != nil {
			c.logger.Error(c.sinkerTableConfig.TableName+" conn.Write error", zap.Error(err), zap.Int("bufLength", bufLength))
		} else {
			c.logger.Info(c.sinkerTableConfig.TableName+" flush msg success.", zap.Int("bufLength", bufLength))
			// todo: wal标记删除状态
		}
	}

	ticker := time.NewTicker(time.Duration(c.sinkerTableConfig.FlushInterval) * time.Second)
	defer ticker.Stop()
	var newKeys map[string]string
	for {
		select {
		case fetch := <-c.fetchCH:
			if c.state.Load() == StateStopped {
				continue
			}
			tmpData := fetch.GetData()
			for i := range tmpData {
				var parse parser.Metric
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
				// todo:写入文件
			}
			data.Data = append(data.Data, tmpData...)
			data.Callback = append(data.Callback, fetch.GetCallback()...)
			bufLength := len(data.Data)
			if bufLength > c.sinkerTableConfig.BufferSize {
				tmp := data.Copy()
				go flushFn(tmp)
				data = FetchArray{}
				ticker.Reset(time.Duration(c.sinkerTableConfig.FlushInterval) * time.Second)
			}
		case <-ticker.C:
			c.logger.Info("--------------------------------------------")
			tmp := data.Copy()
			go flushFn(tmp)
			data = FetchArray{}
		case <-c.ctx.Done():
			c.logger.Info("stopped processing loop", zap.String("table", c.sinkerTableConfig.Database+"."+c.sinkerTableConfig.TableName))
			return
		}
	}
}
