package main

import "golang.org/x/oauth2"

func main() {
	println(len(oauth2.GenerateVerifier()))
	println(len("GxB2VAGUYx93eNSHIyZOzfyxP3hlrfBT"))
}
