package domain

import (
	"encoding/json"
	"testing"
	"time"
)

func TestWriter_WriterMsg(t *testing.T) {
	wm := NewCkWriter(d)
	for num := 0; num < 4; num++ {
		go func() {
			for i := 0; i < 100000; i++ {
				var js, _ = json.Marshal(map[string]interface{}{
					"a":  i,
					"b":  2,
					"c":  3,
					"a1": "1",
					"a2": "2",
					"a3": "3",
				})
				err := wm.WriterMsg(ctx, string(js))
				if err != nil {
					t.Error(err)
				}
			}
		}()
	}
	time.Sleep(time.Hour * 10)
}

func BenchmarkWriterMsg(b *testing.B) {
	wm := NewCkWriter(d)
	for i := 0; i < b.N; i++ {
		var js, _ = json.Marshal(map[string]interface{}{
			"a":  i,
			"b":  2,
			"c":  3,
			"a1": "1",
			"a2": "2",
			"a3": "3",
		})
		err := wm.WriterMsg(ctx, string(js))
		if err != nil {
			b.Error(err)
		}
	}
}
