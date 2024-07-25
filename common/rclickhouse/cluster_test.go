package rclickhouse

import (
	"context"
	"go.uber.org/zap"
	"testing"
)

var logger *zap.Logger
var cluster *ClickhouseCluster
var ctx = context.Background()

var clusterConfig = ClickhouseConfig{
	Cluster: "",
	DB:      "test",
	Hosts: [][]string{
		{
			"127.0.0.1:9000",
		},
	},
	Protocol: "native",
	Username: "root",
	Password: "yeyangfengqi",
}

func TestMain(m *testing.M) {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	cluster, err = InitClusterConn(&clusterConfig)
	if err != nil {
		panic(err)
	}
	m.Run()
}
