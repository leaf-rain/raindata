package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type WheelWheelPools struct {
	rd redis.Cmdable
}

func (r *WheelWheelPools) BaseDataGet(ctx context.Context, cName, key, field string) (string, error) {
	return r.rd.HGet(ctx, key, field).Result()
}

func (r *WheelWheelPools) BaseDataSet(ctx context.Context, cName, key, field, value string) error {
	return r.rd.HSet(ctx, key, field, value).Err()
}

func (r *WheelWheelPools) BaseDataGetAll() <-chan ScanEvent {
	var ch = make(chan ScanEvent)
	go func() {
		defer close(ch)
		var ctx = context.Background()
		// 要扫描的哈希表键名
		hashKey := "wheel_WheelPools"
		// 初始化游标
		cursor := uint64(0)
		// 每次拉取 500 条数据
		count := int64(500)
		for {
			// 执行 HSCAN 命令
			vals, nextCursor, err := r.rd.HScan(ctx, hashKey, cursor, "", count).Result()
			if err != nil {
				fmt.Println("HSCAN error:", err)
				break
			}
			// 处理当前批次的数据
			for i := 0; i < len(vals); i += 2 {
				field := vals[i]
				value := vals[i+1]
				//fmt.Printf("Field: %s, Value: %s\n", field, value)
				ch <- ScanEvent{
					CName: field,
					Key:   hashKey,
					Field: field,
					Value: value,
				}
			}
			// 更新游标
			cursor = nextCursor
			// 如果游标为 0，表示扫描完成
			if cursor == 0 {
				break
			}
			// 测试
			break
		}
	}()
	return ch
}
