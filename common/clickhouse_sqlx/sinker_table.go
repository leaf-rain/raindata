package clickhouse_sqlx

import (
	"context"
	"go.uber.org/zap"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type SinkerTableConfig struct {
	TableName string
	Database  string

	// AutoSchema will auto fetch the schema from clickhouse
	AutoSchema     bool
	ExcludeColumns []string
	Dims           []struct {
		Name       string
		Type       string
		SourceName string
	} `json:"dims"`
	// DynamicSchema will add columns present in message to clickhouse. Requires AutoSchema be true.
	DynamicSchema struct {
		Enable      bool
		NotNullable bool
		MaxDims     int // the upper limit of dynamic columns number, <=0 means math.MaxInt16. protecting dirty data attack
		// A column is added for new key K if all following conditions are true:
		// - K isn't in ExcludeColumns
		// - number of existing columns doesn't reach MaxDims-1
		// - WhiteList is empty, or K matchs WhiteList
		// - BlackList is empty, or K doesn't match BlackList
		WhiteList string // the regexp of white list
		BlackList string // the regexp of black list
	}

	FlushInterval int
	BufferSize    int
	TimeZone      string  `json:"timeZone"`
	TimeUnit      float64 `json:"timeUnit"`
}

// SinkerTable object maintains number of task for each partition
type SinkerTable struct {
	logger      *zap.Logger
	config      *ClickhouseConfig
	clusterConn *ClickhouseCluster
	ctx         context.Context
	cancel      context.CancelFunc
	dataCh      chan interface{}
	exitCh      chan struct{}
	processWg   sync.WaitGroup
	mux         sync.Mutex
	commitDone  *sync.Cond
	state       atomic.Uint32

	data []interface{}

	// table相关
	sinkerTableConfig *SinkerTableConfig
}

// NewSinkerTable get an instance of sinker with the task list
func NewSinkerTable(ctx context.Context, config *ClickhouseConfig, cc *ClickhouseCluster, logger *zap.Logger, sinkerTableConfig *SinkerTableConfig) (*SinkerTable, error) {
	s := &SinkerTable{
		logger:            logger,
		config:            config,
		clusterConn:       cc,
		ctx:               ctx,
		dataCh:            make(chan interface{}, 1000),
		exitCh:            make(chan struct{}),
		sinkerTableConfig: sinkerTableConfig,
	}
	s.state.Store(StateStopped)
	s.commitDone = sync.NewCond(&s.mux)
	return s, nil
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
		var wg sync.WaitGroup
		// todo:将数据写进数据库
	}

	ticker := time.NewTicker(time.Duration(c.sinkerTableConfig.FlushInterval) * time.Second)
	defer ticker.Stop()
	traceId := "NO_RECORDS_FETCHED"
	wait := false
	for {
		select {
		case data := <-c.dataCh:
			if c.state.Load() == StateStopped {
				continue
			}
			fetch := fetches.Fetch.Records()
			if wait {
				c.logger.Info(c.sinkerTableConfig.TableName+" flush msg.", zap.String("traceId", traceId),
					zap.String("message", "bufThreshold not reached, use old traceId"),
					zap.String("trace_id", traceId),
					zap.Int("records", len(fetch)),
					zap.Int("totalLength", bufLength))
			} else {
				traceId = fetches.TraceId
				util.LogTrace(traceId, util.TraceKindProcessStart, zap.Int("records", len(fetch)))
			}
			items, done := int64(len(fetch)), int64(-1)
			var concurrency int
			if concurrency = int(items/1000) + 1; concurrency > MaxParallelism {
				concurrency = MaxParallelism
			}

			var wg sync.WaitGroup
			var err error
			wg.Add(concurrency)
			for i := 0; i < concurrency; i++ {
				go func() {
					for {
						index := atomic.AddInt64(&done, 1)
						if index >= items || c.state.Load() == util.StateStopped {
							wg.Done()
							break
						}

						rec := fetch[index]
						msg := &model.InputMessage{
							Topic:     rec.Topic,
							Partition: int(rec.Partition),
							Key:       rec.Key,
							Value:     rec.Value,
							Offset:    rec.Offset,
							Timestamp: &rec.Timestamp,
						}
						tablename := ""
						for _, it := range rec.Headers {
							if it.Key == "__table_name" {
								tablename = string(it.Value)
								break
							}
						}

						c.tasks.Range(func(key, value any) bool {
							tsk := value.(*Service)
							if (tablename != "" && tsk.clickhouse.TableName == tablename) || tsk.taskCfg.Topic == rec.Topic {
								//bufLength++
								atomic.AddInt64(&bufLength, 1)
								if e := tsk.Put(msg, traceId, flushFn); e != nil {
									atomic.StoreInt64(&done, items)
									err = e
									// decrise the error record
									util.Rs.Dec(1)
									return false
								}
							}
							return true
						})
					}
				}()
			}
			wg.Wait()

			// record the latest offset in order
			// assume the c.state was reset to stopped when facing error, so that further fetch won't get processed
			if err == nil {
				for _, f := range *fetches.Fetch {
					for i := range f.Topics {
						ft := &f.Topics[i]
						if recMap[ft.Topic] == nil {
							recMap[ft.Topic] = make(map[int32]*model.BatchRange)
						}
						for j := range ft.Partitions {
							fpr := ft.Partitions[j].Records
							if len(fpr) == 0 {
								continue
							}
							lastOff := fpr[len(fpr)-1].Offset
							firstOff := fpr[0].Offset

							or, ok := recMap[ft.Topic][ft.Partitions[j].Partition]
							if !ok {
								or = &model.BatchRange{Begin: math.MaxInt64, End: -1}
								recMap[ft.Topic][ft.Partitions[j].Partition] = or
							}
							if or.End < lastOff {
								or.End = lastOff
							}
							if or.Begin > firstOff {
								or.Begin = firstOff
							}
						}
					}
				}
			}

			if bufLength > int64(bufThreshold) {
				flushFn(traceId, "bufLength reached")
				ticker.Reset(time.Duration(c.grpConfig.FlushInterval) * time.Second)
				wait = false
			} else {
				wait = true
			}
		case <-ticker.C:
			flushFn(traceId, "ticker.C triggered")
		case <-c.ctx.Done():
			util.Logger.Info("stopped processing loop", zap.String("group", c.grpConfig.Name))
			return
		}
	}
}
