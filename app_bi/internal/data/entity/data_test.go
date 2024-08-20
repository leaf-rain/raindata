package entity

import (
	"flag"
	"fmt"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/common/logger"
	"go.uber.org/zap"
	"testing"
)

var zapLogger *zap.Logger
var bc conf.Bootstrap
var data *Data

func TestMain(m *testing.M) {
	flag.Parse()
	c := config.New(
		config.WithSource(
			file.NewSource("../../../configs"),
		),
	)
	if err := c.Load(); err != nil {
		panic(err)
	}
	err := c.Scan(&bc)
	if err != nil {
		panic(err)
	}
	zapLogger, err = logger.InitLogger(&logger.LogConfig{
		ServerName:        bc.Log.ServerName,
		Appid:             bc.Log.Appid,
		LogLevel:          bc.Log.LogLevel,
		LogFormat:         bc.Log.LogFormat,
		LogFile:           bc.Log.LogFile,
		LogPath:           bc.Log.LogPath,
		LogFileMaxSize:    int(bc.Log.LogFileMaxSize),
		LogFileMaxBackups: int(bc.Log.LogFileMaxBackups),
		LogMaxAge:         int(bc.Log.LogMaxAge),
		LogCompress:       bc.Log.LogCompress,
	})
	if err != nil {
		panic(err)
	}
	var cancel func()
	data, cancel, err = NewData(&bc, zapLogger)
	if err != nil {
		panic(err)
	}
	defer cancel()
	m.Run()
}

func TestName(t *testing.T) {
	fmt.Print("--------------------------------")
}
