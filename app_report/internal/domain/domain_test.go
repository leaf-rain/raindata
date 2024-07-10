package domain

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var d *Domain
var ctx = context.Background()

func TestMain(m *testing.M) {
	var err error
	d, err = Initialize()
	if err != nil {
		panic(err)
	}
	now := time.Now()
	m.Run()
	fmt.Println("程序执行耗时:", time.Since(now))
}
