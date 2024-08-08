package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/leaf-rain/raindata/common/ecode"
	"github.com/leaf-rain/raindata/common/logger"
	"github.com/leaf-rain/raindata/common/rclickhouse"
	"github.com/leaf-rain/raindata/common/rsql"
	"github.com/spf13/viper"
	"log"
	"net"
)

type Config struct {
	HttpAddr     string                        `json:"HttpAddr" yaml:"HttpAddr"`
	GrpcAddr     string                        `json:"GrpcAddr" yaml:"GrpcAddr"`
	Version      string                        `json:"Version" yaml:"Version"`
	Mode         string                        `json:"Mode" yaml:"Mode"`
	SecretKey    string                        `json:"SecretKey" yaml:"SecretKey"`
	LogConfig    *logger.LogConfig             `json:"LogConfig" yaml:"LogConfig"`
	CKConfig     *rclickhouse.ClickhouseConfig `json:"CKConfig" yaml:"CKConfig"`
	SqlConfig    *rsql.SqlConfig               `json:"SqlConfig" yaml:"SqlConfig"`
	HttpListener net.Listener                  `json:"-" yaml:"-"`
	GrpcListener net.Listener                  `json:"-" yaml:"-"`
}

func GetLogCfgByConfig(cfg *Config) *logger.LogConfig {
	return cfg.LogConfig
}

func GetCKCfgByConfig(cfg *Config) *rclickhouse.ClickhouseConfig {
	return cfg.CKConfig
}

func GetSqlCfgByConfig(cfg *Config) *rsql.SqlConfig {
	return cfg.SqlConfig
}

func InitConfig(opt *CmdArgs) (cfg *Config, err error) {
	cfg = new(Config)
	var cfPath = opt.ConfigFile
	if *cfPath == "" {
		log.Println("config file path is  empty")
		return nil, ecode.ERR_CONFIG_PATH
	}
	vip := viper.New()
	vip.SetConfigFile(*cfPath)
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
	cfg.Mode = *opt.Mode
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
