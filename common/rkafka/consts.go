package rkafka

const (
	MaxBufferSize                  = 1 << 20 // 1048576
	defaultBufferSize              = 1 << 18 // 262144
	maxFlushInterval               = 600
	defaultFlushInterval           = 10
	defaultTimeZone                = "Local"
	defaultLogLevel                = "info"
	defaultKerberosConfigPath      = "/etc/krb5.conf"
	defaultMaxOpenConns            = 1
	defaultReloadSeriesMapInterval = 3600   // 1 hour
	defaultActiveSeriesRange       = 86400  // 1 day
	defaultHeartbeatInterval       = 3000   // 3 s
	defaultSessionTimeout          = 120000 // 2 min
	defaultRebalanceTimeout        = 120000 // 2 min
	defaultRequestTimeoutOverhead  = 60000  // 1 min
)
