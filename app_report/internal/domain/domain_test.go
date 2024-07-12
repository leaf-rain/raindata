package domain

import (
	"context"
	"fmt"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure/config"
	"testing"
	"time"
)

var d *Domain
var ctx = context.Background()

func TestMain(m *testing.M) {
	var path = "../../configs/config.yaml"
	//flag.Set("configFile", path)
	*config.Args.ConfigFile = path
	var err error
	d, err = Initialize()
	if err != nil {
		panic(err)
	}
	now := time.Now()
	m.Run()
	fmt.Println("程序执行耗时:", time.Since(now))
}
