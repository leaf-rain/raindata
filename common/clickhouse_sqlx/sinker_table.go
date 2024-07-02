package clickhouse_sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/leaf-rain/raindata/common/parser"
	"go.uber.org/zap"
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
}

// SinkerTable object maintains number of task for each partition
type SinkerTable struct {
	logger      *zap.Logger
	config      *ClickhouseConfig
	clusterConn *ClickhouseCluster
	ctx         context.Context
	cancel      context.CancelFunc
	fetchCH     chan Fetch
	exitCh      chan struct{}
	processWg   sync.WaitGroup
	mux         sync.Mutex
	commitDone  *sync.Cond
	state       atomic.Uint32
	parser      parser.Parser
	fieldMap    *sync.Map
	// table相关
	sinkerTableConfig *SinkerTableConfig
	columns           Columns
	// todo:wal功能
}

// NewSinkerTable get an instance of sinker with the task list
func NewSinkerTable(ctx context.Context, config *ClickhouseConfig, cc *ClickhouseCluster, logger *zap.Logger, sinkerTableConfig *SinkerTableConfig, parser parser.Parser) (*SinkerTable, error) {
	s := &SinkerTable{
		logger:            logger,
		config:            config,
		clusterConn:       cc,
		ctx:               ctx,
		fetchCH:           make(chan Fetch, 1000),
		exitCh:            make(chan struct{}),
		sinkerTableConfig: sinkerTableConfig,
		parser:            parser,
		fieldMap:          new(sync.Map),
	}
	s.state.Store(StateStopped)
	s.commitDone = sync.NewCond(&s.mux)
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
	// todo:查询表已经存在的字段信息存储到field map中
	conn := c.clusterConn.GetShardConn(time.Now().Unix())
	ck, _, err := conn.NextGoodReplica(0)
	if err != nil {
		return err
	}
	var query = fmt.Sprintf(`select name, type, default_kind from system.columns where database = '%s' and table = '%s'`,
		c.sinkerTableConfig.Database, c.sinkerTableConfig.TableName)
	// todo:如果表不存在，则根据初始化字段or配置创建表。
	var rs *sql.Rows
	rs, err = ck.db.Query(query)
	if err != nil {
		// todo:处理没有这张表的情况，创建表&添加基础字段。
		return err
	}
	defer rs.Close()
	var name, typ, defaultKind string
	for rs.Next() {
		if err = rs.Scan(&name, &typ, &defaultKind); err != nil {
			return err
		}
		c.columns.Put(name, &ColumnWithType{Name: name, Type: WhichType(typ)})
	}
	return nil
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

func (c *SinkerTable) processFetch() {
	c.processWg.Add(1)
	defer c.processWg.Done()
	var bufLength int
	var data []string
	flushFn := func(data []string) {
		bufLength = len(data)
		if bufLength > 0 {
			c.logger.Info(c.sinkerTableConfig.TableName+" flush msg.", zap.Int("bufLength", bufLength))
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
			// todo: 校验是否有没有添加的字段
			for i := range tmpData {
				parse, err = c.parser.Parse([]byte(tmpData[i]))
				if err != nil {
					c.logger.Error(c.sinkerTableConfig.TableName+" parse error", zap.Error(err), zap.String("data", tmpData[i]))
					continue
				}
				newKeys = parse.GetNewKeys(c.fieldMap)
				if len(newKeys) > 0 {
					// todo:添加表字段
				}
			}
			data = append(data, tmpData...)
			bufLength = len(data)
			if bufLength > c.sinkerTableConfig.BufferSize {
				tmp := make([]string, len(data))
				copy(tmp, data)
				data = data[:0]
				go flushFn(tmp)
				ticker.Reset(time.Duration(c.sinkerTableConfig.FlushInterval) * time.Second)
			}
		case <-ticker.C:
			tmp := make([]string, len(data))
			copy(tmp, data)
			data = data[:0]
			go flushFn(tmp)
		case <-c.ctx.Done():
			c.logger.Info("stopped processing loop", zap.String("table", c.sinkerTableConfig.Database+"."+c.sinkerTableConfig.TableName))
			return
		}
	}
}
