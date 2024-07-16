package go_sql_driver

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const dsn = "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true"

// SqlConfig configuration parameters
type SqlConfig struct {
	DB           string `json:"db,omitempty" yaml:"db"`
	Host         string `json:"host,omitempty" yaml:"host"`
	Port         string `json:"port,omitempty" yaml:"port"`
	Username     string `json:"username,omitempty" yaml:"username"`
	Password     string `json:"password,omitempty" yaml:"password"`
	MaxOpenConns int    `json:"maxOpenConns,omitempty" yaml:"maxOpenConns"`
	MaxIdleConns int    `json:"maxIdleConns,omitempty" yaml:"maxIdleConns"`
}

func NewSql(cfg *SqlConfig) (*sql.DB, error) {
	dbDsn := fmt.Sprintf(dsn,
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB)
	db, _ := sql.Open("mysql", dbDsn)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Hour) // mysql default conn timeout=8h, should < mysql_timeout
	err := db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
