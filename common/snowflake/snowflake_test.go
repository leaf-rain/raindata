package snowflake

import (
	"fmt"
	"testing"
)

func TestSnowflakeInt64(t *testing.T) {
	for i := 0; i < 1000; i++ {
		var num = SnowflakeInt32()
		fmt.Println(num)
		if num < 1000000 {
			fmt.Println(num)
		}
	}
}
