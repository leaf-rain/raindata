package clickhouse_sqlx

import (
	"testing"
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
	st.start()
}
