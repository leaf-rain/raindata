package config

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/leaf-rain/raindata/app_report/pkg/logger"
	"github.com/leaf-rain/raindata/common/clickhouse_sqlx"
	"github.com/leaf-rain/raindata/common/ecode"
	"github.com/leaf-rain/raindata/common/etcd"
	"github.com/spf13/viper"
	"log"
	"net"
)

type Config struct {
	Ctx              context.Context                   `json:"-" yaml:"-"`
	HttpAddr         string                            `json:"HttpAddr" yaml:"HttpAddr"`
	GrpcAddr         string                            `json:"GrpcAddr" yaml:"GrpcAddr"`
	Version          string                            `json:"Version" yaml:"Version"`
	Mode             string                            `json:"Mode" yaml:"Mode"`
	SecretKey        string                            `json:"SecretKey" yaml:"SecretKey"`
	LogConfig        *logger.LogConfig                 `json:"LogConfig" yaml:"LogConfig"`
	MetadataPath     string                            `json:"MetadataPath" yaml:"MetadataPath"`
	EtcdConfig       *etcd.Config                      `json:"EtcdConfig" yaml:"EtcdConfig"`
	ClickhouseConfig *clickhouse_sqlx.ClickhouseConfig `json:"ClickhouseConfig" yaml:"ClickhouseConfig"`
	HttpListener     net.Listener                      `json:"-" yaml:"-"`
	GrpcListener     net.Listener                      `json:"-" yaml:"-"`
}

func GetLogCfgByConfig(cfg *Config) *logger.LogConfig {
	return cfg.LogConfig
}

func GetCtx(cfg *Config) context.Context {
	return cfg.Ctx
}

func GetEtcdConfig(cfg *Config) *etcd.Config {
	return cfg.EtcdConfig
}

func GetClickhouseConfig(cfg *Config) *clickhouse_sqlx.ClickhouseConfig {
	return cfg.ClickhouseConfig
}

func InitConfig(opt *CmdArgs) (cfg *Config, err error) {
	cfg = new(Config)
	cfg.Ctx = context.Background()
	var cfPath = opt.ConfigFile
	if cfPath == "" {
		log.Println("config file path is  empty")
		return nil, ecode.ERR_CONFIG_PATH
	}
	vip := viper.New()
	vip.SetConfigFile(cfPath)
	err = vip.ReadInConfig()
	if err != nil {
		log.Println("config file vip.ReadInConfig err:", err)
		return nil, ecode.ERR_CONFIG_UNMARSHAL
	}
	err = vip.Unmarshal(cfg)
	if err != nil {
		log.Println("config file vip.Unmarshal err:", err)
		return nil, ecode.ERR_CONFIG_UNMARSHAL
	}
	cfg.Mode = opt.Mode
	log.Println("config file loading success.", cfg)
	// 加载端口
	if cfg.HttpAddr != "" {
		cfg.HttpListener, err = net.Listen("tcp", cfg.HttpAddr)
		if err != nil {
			log.Println("config file http net.Listen err:", err)
			return nil, ecode.ERR_HTTP_CONFIG
		}
	}
	if cfg.GrpcAddr != "" {
		cfg.GrpcListener, err = net.Listen("tcp", cfg.GrpcAddr)
		if err != nil {
			log.Println("config file grpc net.Listen err:", err)
			return nil, ecode.ERR_GRPC_CONFIG
		}
	}
	go dynamicConfig(vip, cfg)
	return cfg, err
}

func dynamicConfig(vc *viper.Viper, cfg *Config) {
	vc.WatchConfig()
	vc.OnConfigChange(func(event fsnotify.Event) {
		fmt.Printf("Detect config change: %s \n", event.String())
		if err := vc.Unmarshal(cfg); err != nil {
			log.Println("config Unmarshal err", err)
		}
	})
}
