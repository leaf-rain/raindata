package clickhouse_sqlx

const (
	MaxBufferSize                  = 1 << 20 // 1048576
	defaultBufferSize              = 1 << 18 // 262144
	maxFlushInterval               = 600
	defaultFlushInterval           = 10
	defaultTimeZone                = "Local"
	defaultLogLevel                = "info"
	defaultMaxOpenConns            = 1
	defaultReloadSeriesMapInterval = 3600   // 1 hour
	defaultActiveSeriesRange       = 86400  // 1 day
	defaultHeartbeatInterval       = 3000   // 3 s
	defaultSessionTimeout          = 120000 // 2 min
	defaultRebalanceTimeout        = 120000 // 2 min
	defaultRequestTimeoutOverhead  = 60000  // 1 min
)

const (
	EngineMergeTree          = 1
	EngineReplacingMergeTree = 2
)
