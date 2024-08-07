package rclickhouse

import (
	"fmt"
	"testing"
	"time"
)

func TestNewSinkerTable(t *testing.T) {
	var stConfig = &SinkerTableConfig{
		TableName:     "test1",
		Database:      "test",
		FlushInterval: 10,
		BufferSize:    10000,
		Parse:         "fastjson",
		BaseColumn: []ColumnWithType{
			{
				Name: "id",
				Type: typeInfo["Int64"],
			},
			{
				Name: "contain",
				Type: typeInfo["String"],
			},
		},
		ReplayKey:  "id",
		OrderByKey: "id",
	}
	st, err := NewSinkerTable(ctx, cluster, logger, stConfig, "fastjson")
	if err != nil {
		t.Fatal(err)
	}
	st.Start()
	go func() {
		var now = time.Now()
		for i := 0; i < 1000000; i++ {
			if i%10 == 0 {
				st.fetchCH <- FetchSingle{
					Data: fmt.Sprintf("{\"id\":%d,\"contain\":\"test\",\"t1\":\"test\"}", i),
					Callback: func() {
						t.Log("callback")
					},
				}
			} else {
				st.fetchCH <- FetchSingle{
					Data: fmt.Sprintf("{\"id\":%d,\"contain\":\"test\"}", i),
					Callback: func() {
						t.Log("callback")
					},
				}
			}

		}
		t.Log("-------------------------------->", time.Since(now))
	}()
	time.Sleep(time.Minute * 30)
}
