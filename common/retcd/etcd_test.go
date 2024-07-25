package retcd

import (
	"context"
	clientv3 "go.retcd.io/retcd/client/v3"
	"go.uber.org/zap"
	"testing"
)

var ctx = context.Background()
var cli *clientv3.Client
var etcdConfig = &Config{
	Endpoints: []string{"127.0.0.1:2379"},
}

func TestMain(m *testing.M) {
	var err error
	cli, err = NewEtcdClient(etcdConfig, zap.NewNop())
	if err != nil {
		panic(err)
	}
	m.Run()
}
