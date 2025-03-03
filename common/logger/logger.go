// 日志引擎层
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
	"path/filepath"
)

const DefaultLogPath = "logs" // 默认输出日志文件路径

var logLevel = map[string]zapcore.Level{
	leverDebug: zapcore.DebugLevel,
	leverInfo:  zapcore.InfoLevel,
	leverWarn:  zapcore.WarnLevel,
	leverError: zapcore.ErrorLevel,
}

const (
	leverDebug = "debug"
	leverInfo  = "info"
	leverWarn  = "warn"
	leverError = "error"
)

type LogConfig struct {
	ServerName string `json:"server_name,omitempty" yaml:"server_name"`
	Appid      string `json:"appid,omitempty" yaml:"appid"`

	LogLevel  string `json:"log_level,omitempty" yaml:"log_level"`   // 日志打印级别 debug  info  warning  error
	LogFormat string `json:"log_format,omitempty" yaml:"log_format"` // 输出日志格式	logfmt, json, 默认json

	LogFile           bool   `json:"log_file,omitempty" yaml:"log_file"`                         // 是否输出到文件
	LogPath           string `json:"log_path,omitempty" yaml:"log_path"`                         // 输出日志文件路径
	LogFileMaxSize    int    `json:"log_file_max_size,omitempty" yaml:"log_file_max_size"`       //【日志分割】单个日志文件最多存储量 单位(mb)
	LogFileMaxBackups int    `json:"log_file_max_backups,omitempty" yaml:"log_file_max_backups"` //【日志分割】日志备份文件最多数量
	LogMaxAge         int    `json:"log_max_age,omitempty" yaml:"log_max_age"`                   // 日志保留时间，单位: 天 (day)
	LogCompress       bool   `json:"log_compress,omitempty" yaml:"log_compress"`                 // 是否压缩日志
}

// 初始化日志 logger
func InitLogger(cfg *LogConfig, opts ...zap.Option) (*zap.Logger, error) {
	if cfg == nil {
		log.Println("InitLogger log config is nil")
		cfg = new(LogConfig)
	}
	_, ok := logLevel[cfg.LogLevel]
	if !ok {
		cfg.LogLevel = "info"
	}
	if cfg.LogFile {
		if cfg.LogPath == "" {
			cfg.LogPath = DefaultLogPath
		}
		if cfg.LogFileMaxSize == 0 {
			cfg.LogFileMaxSize = 300
		}
		if cfg.LogMaxAge == 0 {
			cfg.LogFileMaxSize = 5
		}
	}
	encoder := getEncoder(cfg)
	core, err := getCore(cfg, encoder)
	if err != nil {
		return nil, err
	}
	log.Println("InitLogger server log loading success. path:", cfg.LogPath)
	var fields []zap.Field
	if cfg.ServerName != "" {
		fields = append(fields, zap.String("server_name", cfg.ServerName))
	}
	if cfg.Appid != "" {
		fields = append(fields, zap.String("appid", cfg.Appid))
	}
	if len(fields) > 0 {
		opts = append(opts, zap.AddCaller(), zap.Development(), zap.Fields(fields...))
		logger := zap.New(core, opts...)
		return logger, nil
	}
	opts = append(opts, zap.AddCaller(), zap.Development())
	logger := zap.New(core, opts...)
	return logger, nil
}

// getEncoder 编码器(如何写入日志)
func getEncoder(conf *LogConfig) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder     // log 时间格式 例如: 2021-09-11t20:05:54.852+0800
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder // 全小写
	if conf.LogFormat == "logfmt" {
		return zapcore.NewConsoleEncoder(encoderConfig) // 以logfmt格式写入
	}
	return zapcore.NewJSONEncoder(encoderConfig) // 以json格式写入
}

// getLogWriter 获取日志输出方式  日志文件 控制台
func getCore(conf *LogConfig, encoder zapcore.Encoder) (zapcore.Core, error) {
	allLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logLevel[conf.LogLevel]
	})
	if conf.LogFile {
		// 判断日志路径是否存在，如果不存在就创建
		if exist := IsExist(conf.LogPath); !exist {
			if conf.LogPath == "" {
				conf.LogPath = DefaultLogPath
			}
			if err := os.MkdirAll(conf.LogPath, os.ModePerm); err != nil {
				conf.LogPath = DefaultLogPath
				if err := os.MkdirAll(conf.LogPath, os.ModePerm); err != nil {
					return nil, err
				}
			}
		}
		var cores []zapcore.Core
		for k, v := range logLevel {
			if v >= logLevel[conf.LogLevel] {
				f := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
					return lvl == v
				})
				writer := getWriter(conf, k)
				cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(writer), f))
			}
		}
		cores = append(cores, zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), allLevel))
		return zapcore.NewTee(cores...), nil
	} else {
		// 日志只输出到控制台
		return zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), allLevel), nil
	}
}

func getWriter(conf *LogConfig, lever string) io.Writer {
	return &lumberjack.Logger{
		Filename:   filepath.Join(conf.LogPath, lever+".log"), // 日志文件路径
		MaxSize:    conf.LogFileMaxSize,                       // 单个日志文件最大多少 mb
		MaxBackups: conf.LogFileMaxBackups,                    // 日志备份数量
		MaxAge:     conf.LogMaxAge,                            // 日志最长保留时间
		Compress:   conf.LogCompress,                          // 是否压缩日志
	}
}

// IsExist 判断文件或者目录是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
