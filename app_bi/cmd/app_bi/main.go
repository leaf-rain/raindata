package main

import (
	"context"
	"flag"
	"github.com/leaf-rain/raindata/app_bi/third_party/klogs"
	"os"

	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/common/logger"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
}

func newApp(lg log.Logger, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(lg),
		kratos.Server(
			hs,
		),
	)
}

func main() {
	flag.Parse()
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()
	if err := c.Load(); err != nil {
		panic(err)
	}
	var ctx = context.Background()
	var bc *conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	zapLogger, err := logger.InitLogger(&logger.LogConfig{
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
	lg := klogs.NewLogger(zapLogger)

	app, cleanup, err := wireApp(ctx, bc.Server, bc, zapLogger, lg)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
