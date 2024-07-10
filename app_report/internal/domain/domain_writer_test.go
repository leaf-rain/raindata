package domain

import (
	"encoding/json"
	"testing"
	"time"
)

func TestWriter_WriterMsg(t *testing.T) {
	wm := NewCkWriter(d)
	var js, _ = json.Marshal(map[string]interface{}{
		"a":  1,
		"b":  2,
		"c":  3,
		"a1": "1",
		"a2": "2",
		"a3": "3",
	})
	err := wm.WriterMsg(ctx, 1, string(js))
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Hour * 10)
}
