package main

import (
	"fmt"
	"github.com/leaf-rain/fastjson"
	"github.com/leaf-rain/raindata/common/snowflake"
	"sync"
)

func main() {
	var t = new(sync.Map)
	t.Store("a", nil)
	t.Store("a1", nil)
	t.Store("a2", nil)
	t.Store("a3", nil)
	t.Store("b", nil)
	t.Store("c", nil)
	t.Store("_id", nil)
	t.Store("event", nil)
	var wg sync.WaitGroup
	var pp fastjson.ParserPool
	for n := 0; n < 8; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//for i := 0; i < 100000000; i++ {
			//	var id = snowflake.SnowflakeInt64()
			//	var js = []byte(fmt.Sprintf(`{"a":%d,"a1":"1","a2":"2","a3":"3","b":2,"c":3,"_id":%d,"event":"log"}`, i, id))
			//	m, _ := fastjson.ParseBytes(js)
			//	obj, _ := m.Object()
			//	obj.Visit(func(key []byte, v *fastjson.Value) {
			//		_, ok := t.Load(string(key))
			//		if !ok {
			//			panic("-----------")
			//		}
			//	})
			//}
			for i := 0; i < 100000000; i++ {
				var id = snowflake.SnowflakeInt64()
				var js = []byte(fmt.Sprintf(`{"a":%d,"a1":"1","a2":"2","a3":"3","b":2,"c":3,"_id":%d,"event":"log"}`, i, id))
				var p = pp.Get()
				m, _ := p.ParseBytes(js)
				obj, _ := m.Object()
				obj.Visit(func(key []byte, v *fastjson.Value) {
					_, ok := t.Load(string(key))
					if !ok {
						panic("-----------")
					}
				})
				pp.Put(p)
			}
		}()
	}
	wg.Wait()
	//for i := 0; i < 100000000; i++ {
	//	var id = snowflake.SnowflakeInt64()
	//	var js = []byte(fmt.Sprintf(`{"a":%d,"a1":"1","a2":"2","a3":"3","b":2,"c":3,"_id":%d,"event":"log"}`, i, id))
	//	m, _ := ps.Parse(js)
	//	p := m.GetNewKeys(t)
	//	if len(p) > 0 {
	//		panic("-----------")
	//	}
	//}
}
