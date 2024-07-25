package retcd

import (
	clientv3 "go.retcd.io/retcd/client/v3"
	"go.uber.org/zap"
	"time"
)

type Config struct {
	// Endpoints is a list of URLs.
	Endpoints []string `json:"Endpoints,omitempty" yaml:"Endpoints" yml:"Endpoints"`

	// AutoSyncInterval is the interval to update endpoints with its latest members.
	// 0 disables auto-sync. By default auto-sync is disabled.
	AutoSyncInterval int64 `json:"AutoSyncInterval,omitempty" yaml:"AutoSyncInterval" yml:"AutoSyncInterval"`

	// DialTimeout is the timeout for failing to establish a connection.
	DialTimeout int64 `json:"DialTimeout,omitempty" yaml:"DialTimeout" yml:"DialTimeout"`

	// DialKeepAliveTime is the time after which client pings the server to see if
	// transport is alive.
	DialKeepAliveTime int64 `json:"DialKeepAliveTime,omitempty" yaml:"DialKeepAliveTime" yml:"DialKeepAliveTime"`

	// DialKeepAliveTimeout is the time that the client waits for a response for the
	// keep-alive probe. If the response is not received in this time, the connection is closed.
	DialKeepAliveTimeout int64 `json:"DialKeepAliveTimeout,omitempty" yaml:"DialKeepAliveTimeout" yml:"DialKeepAliveTimeout"`

	// MaxCallSendMsgSize is the client-side request send limit in bytes.
	// If 0, it defaults to 2.0 MiB (2 * 1024 * 1024).
	// Make sure that "MaxCallSendMsgSize" < server-side default send/recv limit.
	// ("--max-request-bytes" flag to retcd or "embed.Config.MaxRequestBytes").
	MaxCallSendMsgSize int `json:"MaxCallSendMsgSize,omitempty" yaml:"MaxCallSendMsgSize" yml:"MaxCallSendMsgSize"`

	// MaxCallRecvMsgSize is the client-side response receive limit.
	// If 0, it defaults to "math.MaxInt32", because range response can
	// easily exceed request send limits.
	// Make sure that "MaxCallRecvMsgSize" >= server-side default send/recv limit.
	// ("--max-request-bytes" flag to retcd or "embed.Config.MaxRequestBytes").
	MaxCallRecvMsgSize int `json:"MaxCallRecvMsgSize,omitempty" yaml:"MaxCallRecvMsgSize" yml:"MaxCallRecvMsgSize"`

	// Username is a user name for authentication.
	Username string `json:"Username,omitempty" yaml:"Username" yml:"Username"`

	// Password is a password for authentication.
	Password string `json:"Password,omitempty" yaml:"Password" yml:"Password"`

	// RejectOldCluster when set will refuse to create a client against an outdated cluster.
	RejectOldCluster bool `json:"RejectOldCluster,omitempty" yaml:"RejectOldCluster" yml:"RejectOldCluster"`

	// PermitWithoutStream when set will allow client to send keepalive pings to server without any active streams(RPCs).
	PermitWithoutStream bool `json:"PermitWithoutStream,omitempty" yaml:"PermitWithoutStream" yml:"PermitWithoutStream"`

	// MaxUnaryRetries is the maximum number of retries for unary RPCs.
	MaxUnaryRetries uint `json:"MaxUnaryRetries,omitempty" yaml:"MaxUnaryRetries" yml:"MaxUnaryRetries"`

	// BackoffWaitBetween is the wait time before retrying an RPC.
	BackoffWaitBetween int64 `json:"BackoffWaitBetween,omitempty" yaml:"BackoffWaitBetween" yml:"BackoffWaitBetween"`

	// BackoffJitterFraction is the jitter fraction to randomize backoff wait time.
	BackoffJitterFraction float64 `json:"BackoffJitterFraction,omitempty" yaml:"BackoffJitterFraction" yml:"BackoffJitterFraction"`
}

// NewEtcdClient 创建一个新的EtcdClient实例
func NewEtcdClient(etcdConfig *Config, logger *zap.Logger) (*clientv3.Client, error) {
	if etcdConfig == nil || len(etcdConfig.Endpoints) == 0 {
		return nil, nil
	}
	config := clientv3.Config{
		Endpoints:             etcdConfig.Endpoints,
		AutoSyncInterval:      time.Second * time.Duration(etcdConfig.AutoSyncInterval),
		DialTimeout:           time.Second * time.Duration(etcdConfig.DialTimeout),
		DialKeepAliveTime:     time.Second * time.Duration(etcdConfig.DialKeepAliveTime),
		DialKeepAliveTimeout:  time.Second * time.Duration(etcdConfig.DialKeepAliveTimeout),
		MaxCallSendMsgSize:    etcdConfig.MaxCallSendMsgSize,
		MaxCallRecvMsgSize:    etcdConfig.MaxCallRecvMsgSize,
		Username:              etcdConfig.Username,
		Password:              etcdConfig.Password,
		RejectOldCluster:      etcdConfig.RejectOldCluster,
		Logger:                logger.Named("retcd-client"),
		PermitWithoutStream:   etcdConfig.PermitWithoutStream,
		MaxUnaryRetries:       etcdConfig.MaxUnaryRetries,
		BackoffWaitBetween:    time.Second * time.Duration(etcdConfig.BackoffWaitBetween),
		BackoffJitterFraction: etcdConfig.BackoffJitterFraction,
	}
	c, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}
	return c, nil
}
