package main

import (
	"context"
	"time"
)

func main() {
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*5)
	ctx2, _ := context.WithTimeout(ctx, time.Second*20)
	go func() {
		for {
			select {
			case <-ctx.Done():
				println("ctx timeout")
			case <-ctx2.Done():
				println("ctx2 timeout")
			}
		}
	}()
	time.Sleep(time.Second * 30)
}
