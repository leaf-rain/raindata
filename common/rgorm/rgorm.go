package rgorm

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"log"
	"os"
	"time"
)

type DtGromConfig struct {
	DriverName   string           `yaml:"driverName" json:"driverName"`
	DbSource     string           `yaml:"dbSource" json:"dbSource"`
	MaxOpenConns int              `yaml:"maxOpenConns" json:"maxOpenConns"`
	MaxIdleConns int              `yaml:"maxIdleConns" json:"maxIdleConns"`
	IdleTimeOut  int              `yaml:"idleTimeOut" json:"idleTimeOut"`
	Debug        bool             `yaml:"debug" json:"debug"`
	Logger       logger.Interface `yaml:"-" json:"-"`
}

var _ logger.Interface = (*GormZapLogger)(nil)

func NewGormZapLogger(l *zap.Logger) *GormZapLogger {
	return &GormZapLogger{logger: l.Sugar()}
}

type GormZapLogger struct {
	logger *zap.SugaredLogger
}

func (g GormZapLogger) LogMode(level logger.LogLevel) logger.Interface {
	return g
}

func (g GormZapLogger) Info(ctx context.Context, s string, i ...interface{}) {
	g.logger.Infof(s, i...)
}

func (g GormZapLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	g.logger.Warnf(s, i...)
}

func (g GormZapLogger) Error(ctx context.Context, s string, i ...interface{}) {
	g.logger.Errorf(s, i...)
}

const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

var (
	slowThreshold = 200 * time.Millisecond

	infoStr      = Green + "%s\n" + Reset + Green + "[info] " + Reset
	warnStr      = BlueBold + "%s\n" + Reset + Magenta + "[warn] " + Reset
	errStr       = Magenta + "%s\n" + Reset + Red + "[error] " + Reset
	traceStr     = Green + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
	traceWarnStr = Green + "%s " + Yellow + "%s\n" + Reset + RedBold + "[%.3fms] " + Yellow + "[rows:%v]" + Magenta + " %s" + Reset
	traceErrStr  = RedBold + "%s " + MagentaBold + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
)

func (g GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil && g.logger.Level() >= zapcore.ErrorLevel && (!errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		if rows == -1 {
			g.logger.Errorf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.logger.Errorf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > slowThreshold && g.logger.Level() >= zapcore.WarnLevel:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", slowThreshold)
		if rows == -1 {
			g.logger.Errorf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.logger.Errorf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case g.logger.Level() == zapcore.InfoLevel:
		sql, rows := fc()
		if rows == -1 {
			g.logger.Errorf(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.logger.Errorf(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func NewRGrom(conf DtGromConfig) *gorm.DB {
	if conf.Logger == nil {
		LogLevel := logger.Silent
		if conf.Debug {
			LogLevel = logger.Info
		} else {
			LogLevel = logger.Error
		}
		conf.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      LogLevel,
				Colorful:      true,
			},
		)
	}
	var config = &gorm.Config{
		Logger:               conf.Logger,
		DisableAutomaticPing: true,
	}
	db, err := gorm.Open(mysql.Open(conf.DbSource), config)
	if err != nil {
		panic(err)
	}
	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDb.SetConnMaxLifetime(time.Second * time.Duration(conf.IdleTimeOut))
	sqlDb.SetMaxOpenConns(conf.MaxOpenConns)
	sqlDb.SetMaxIdleConns(conf.MaxIdleConns)
	return db
}

// WithDefaultTimeout 是一个封装了超时逻辑的写入函数
func WithDefaultTimeout(ctx context.Context, db *gorm.DB, outTime int) *gorm.DB {
	if ctx == nil {
		ctx = context.TODO()
	}
	// 定义一个默认超时时间
	if outTime <= 0 {
		outTime = 5
	}
	defaultTimeout := time.Duration(outTime) * time.Second
	// 使用context.WithTimeout创建一个带有超时的上下文
	ctx, _ = context.WithTimeout(ctx, defaultTimeout)
	// 将带有超时的上下文传递给写入操作
	return db.WithContext(ctx)
}
