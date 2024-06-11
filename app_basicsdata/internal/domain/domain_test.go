package domain

import (
	"context"
	"github.com/leaf-rain/raindata/common/str"
	"os"
	"testing"
)

var domain *Domain
var ctx context.Context

func TestMain(m *testing.M) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = os.Chdir(str.RemoveSuffix(dir, "app_basicsdata") + "/app_basicsdata")
	if err != nil {
		panic(err)
	}
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	domain, err = Initialize()
	if err != nil {
		panic(err)
	}
	m.Run()
}
