package rsql

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"testing"
	"time"
)

var st *SinkerTable
var logger = zap.NewExample()
var defaultConfig = &SinkerTableConfig{
	BufferSize:    100,
	Database:      "test",
	FeHost:        "127.0.0.1",
	FeHttpPort:    "8030",
	FlushInterval: 3,
	OrderByKey:    "_id",
	PrimaryKey:    "",
	TableName:     "test",
	TableType:     TableType_Duplicate,
	Username:      "root",
	Password:      "",
	BaseColumn: []ColumnWithType{
		{
			Name: "_id",
			Type: &TypeInfo{
				Type:     INT,
				Nullable: true,
			},
		},
		{
			Name: "name",
			Type: &TypeInfo{
				Type:     STRING,
				Nullable: true,
			},
		},
	},
}

func TestNewSinkerTable(t *testing.T) {
	var err error
	st, err = NewSinkerTable(context.Background(), db, logger, defaultConfig, "fastjson")
	if err != nil {
		t.Fatal(err)
	}
	st.Start()
	go func() {
		for i := 0; i < 200; i++ {
			if i%10 == 0 {
				st.fetchCH <- FetchSingle{
					Data: fmt.Sprintf("{\"_id\":%d,\"contain\":\"test\",\"t1\":\"test\"}", i),
					Callback: func() {
						t.Log("callback")
					},
				}
			} else {
				st.fetchCH <- FetchSingle{
					Data: fmt.Sprintf("{\"_id\":%d,\"contain\":\"test\"}", i),
					Callback: func() {
						t.Log("callback")
					},
				}
			}

		}
	}()
	time.Sleep(time.Minute * 30)
}

func TestSinkerTable_sendStarRocks(t *testing.T) {
	var err error
	st, err = NewSinkerTable(context.Background(), db, logger, defaultConfig, "fastjson")
	if err != nil {
		t.Fatal(err)
	}
	st.sendStarRocks("Sinker_Log_0.1722231973241517155.wal.pending", "10")
}
