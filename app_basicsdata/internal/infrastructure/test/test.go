package main

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
)

func remainDataDecode() {
	type Person struct {
		Name  string
		Age   int
		Other map[string]interface{} `mapstructure:",remain"`
	}

	input := map[string]interface{}{
		"name":   "Tim",
		"age":    31,
		"email":  "one@gmail.com",
		"gender": "male",
	}

	var result Person
	err := mapstructure.Decode(input, &result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", result)
}

func main() {
	remainDataDecode()
}
