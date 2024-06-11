package main

func iterateKeys(jsonStr string, prefix string) {

}

func main() {
	jsonStr := `{
		"name": "John",
		"age": 30,
		"city": "New York",
	}`

	iterateKeys(jsonStr, "*")
}
