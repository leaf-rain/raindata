package go_sql_driver

import (
	"context"
	"database/sql"
	"testing"
)

var defalutSqlConfig = &SqlConfig{
	DB:           "test",
	Host:         "127.0.0.1",
	MaxIdleConns: 10,
	MaxOpenConns: 100,
	Password:     "",
	Port:         "9030",
	Username:     "root",
}

var db *sql.DB

var ctx = context.Background()

func TestMain(m *testing.M) {
	var err error
	db, err = NewSql(defalutSqlConfig)
	if err != nil {
		panic(err)
	}
	m.Run()
}
