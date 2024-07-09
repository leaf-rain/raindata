package main

import "github.com/leaf-rain/fastjson"

func main() {
	var p, _ = fastjson.Parse(`{"t":1}`)
	p2, _ := fastjson.Parse("{}")
	p.Range(func(key []byte, v *fastjson.Value) {
		p2.Set(string(key), v)
	})
	println(p.Type())
	println(p2.String())
}
