package clickhouse_sqlx

import (
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/jmoiron/sqlx"
	"runtime"
	"strings"
)

type ClickhouseConfig struct {
	Host         []string `json:"host" yaml:"host"`
	UserName     string   `json:"userName" yaml:"userName"`
	Password     string   `json:"password" yaml:"password"`
	Database     string   `json:"database" yaml:"database"`
	MaxOpenConns int      `json:"maxOpenConns" yaml:"maxOpenConns"`
	MaxIdleConns int      `json:"maxIdleConns" yaml:"maxIdleConns"`
	Debug        bool     `json:"debug" yaml:"debug"`
}

func (config ClickhouseConfig) ToTcpAddr() string {
	if len(config.Host) == 0 {
		return ""
	}
	var tcpInfo = "tcp://" + config.Host[0]
	if config.Database != "" {
		tcpInfo += "/" + config.Database
	}
	tcpInfo += "?"
	if config.UserName != "" && config.Password != "" {
		tcpInfo += "username=" + config.UserName + "&password=" + config.Password
	}
	if config.Debug {
		tcpInfo += "&debug=true"
	} else {
		tcpInfo += "&debug=false"
	}
	if len(config.Host) > 1 {
		tmpHosts := strings.Join(config.Host[1:], ",")
		tcpInfo += "&alt_hosts=" + tmpHosts
	}
	return tcpInfo
}

type Clickhouse struct {
	*sqlx.DB
}

func NewClickhouse(config *ClickhouseConfig) (*Clickhouse, error) {
	sqlxDb, err := sqlx.Open("clickhouse", config.ToTcpAddr())
	if err != nil {
		return nil, err
	}
	numCPU := runtime.NumCPU()
	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = numCPU
	}
	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = numCPU
	}
	sqlxDb.SetMaxOpenConns(config.MaxOpenConns)
	sqlxDb.SetMaxIdleConns(config.MaxIdleConns)
	err = sqlxDb.Ping()
	if err != nil {
		return nil, err
	}
	return &Clickhouse{
		sqlxDb,
	}, nil
}

func (ck *Clickhouse) Close() error {
	err := ck.DB.Close()
	return err
}
