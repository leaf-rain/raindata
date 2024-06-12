package clickhouse_sqlx

import (
	"context"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
)

// Sinker object maintains number of task for each partition
type Sinker struct {
	logger      *zap.Logger
	config      *ClickhouseConfig
	clusterConn *ClickhouseCluster
	ctx         context.Context
	cancel      context.CancelFunc
	cfgChangeCh chan struct{}
	exitCh      chan struct{}
	processWg   sync.WaitGroup
	mux         sync.Mutex
	commitDone  *sync.Cond
	state       atomic.Uint32
}

// NewSinker get an instance of sinker with the task list
func NewSinker(config *ClickhouseConfig, logger *zap.Logger) (*Sinker, error) {
	clusterConn, err := InitClusterConn(config)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	s := &Sinker{
		logger:      logger,
		config:      config,
		clusterConn: clusterConn,
		ctx:         ctx,
		cancel:      cancel,
		cfgChangeCh: make(chan struct{}),
		exitCh:      make(chan struct{}),
	}
	s.state.Store(StateStopped)
	s.commitDone = sync.NewCond(&s.mux)
	return s, nil
}
