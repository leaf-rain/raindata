package clickhouse_sqlx

import (
	"context"
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
	TableName     string
	Database      string
	FlushInterval int
	BufferSize    int
	TimeZone      string
	TimeUnit      float64
	Parse         string
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
	fieldMap    sync.Map

	data []string

	// table相关
	sinkerTableConfig *SinkerTableConfig
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
	}
	s.state.Store(StateStopped)
	s.commitDone = sync.NewCond(&s.mux)
	s.flushFields()
	return s, nil
}

func (c *SinkerTable) flushFields() {
	// todo:查询表已经存在的字段信息存储到field map中
	conn := c.clusterConn.GetShardConn(time.Now().Unix())
	conn.NextGoodReplica()

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
	flushFn := func(traceId, with string) {
		c.mux.Lock()
		defer c.mux.Unlock()
		bufLength = len(c.data)
		if bufLength > 0 {
			c.logger.Info(c.sinkerTableConfig.TableName+" flush msg.", zap.String("traceId", traceId), zap.String("with", with), zap.Int("bufLength", bufLength))
		}
		// todo: 校验是否有没有添加的字段
		var err error
		var parse parser.Metric
		var tmp = make([][]byte, len(c.data))
		for i := range c.data {
			parse, err = c.parser.Parse([]byte(c.data[i]))
			if err != nil {
				c.logger.Error("[processFetch] parser.Parse failed.", zap.String("traceId", traceId), zap.String("with", with), zap.Error(err))
				continue
			}
			tmp[i] = []byte(c.data[i])
		}
		c.data = c.data[:]
		// todo:将数据写进数据库

	}

	ticker := time.NewTicker(time.Duration(c.sinkerTableConfig.FlushInterval) * time.Second)
	defer ticker.Stop()
	traceId := "NO_RECORDS_FETCHED"
	wait := false
	for {
		select {
		case fetch := <-c.fetchCH:
			if c.state.Load() == StateStopped {
				continue
			}
			if wait {
				c.logger.Info(c.sinkerTableConfig.TableName+" flush msg.", zap.String("traceId", traceId),
					zap.String("message", "bufThreshold not reached, use old traceId"),
					zap.String("trace_id", traceId),
					zap.Int("records", len(fetch.Data)),
					zap.Int("totalLength", bufLength))
			} else {
				traceId = fetch.TraceId
				c.logger.Info("process fetch.", zap.String("traceId", fetch.TraceId), zap.Int("records", len(fetch.Data)))
			}
			// 将数据添加到本地
			c.data = append(c.data, fetch.Data...)
			bufLength = len(c.data)
			if bufLength > c.sinkerTableConfig.BufferSize {
				flushFn(traceId, "bufLength reached")
				ticker.Reset(time.Duration(c.sinkerTableConfig.FlushInterval) * time.Second)
				wait = false
			} else {
				wait = true
			}
		case <-ticker.C:
			flushFn(traceId, "ticker.C triggered")
		case <-c.ctx.Done():
			c.logger.Info("stopped processing loop", zap.String("table", c.sinkerTableConfig.Database+"."+c.sinkerTableConfig.TableName))
			return
		}
	}
}
