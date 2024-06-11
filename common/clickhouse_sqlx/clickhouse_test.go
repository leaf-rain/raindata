package clickhouse_sqlx

import (
	"fmt"
	"testing"
	"time"
)

var ck *Clickhouse

func TestMain(m *testing.M) {
	ckConf := &ClickhouseConfig{
		Host:     []string{"127.0.0.1:9000"},
		UserName: "root",
		Password: "yeyangfengqi",
		Database: "test",
		Debug:    true,
	}
	var err error
	ck, err = NewClickhouse(ckConf)
	if err != nil {
		panic(err)
	}
	var now = time.Now()
	defer fmt.Println("执行耗时:", time.Now().Sub(now))
	m.Run()
}

func TestNewClickhouse(t *testing.T) {
	err := ck.Ping()
	t.Log(err)
}
