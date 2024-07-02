package main

import (
	"fmt"
	"runtime"
)

func main() {
	var a []int
	for i := 0; i < 1000000; i++ {
		a = append(a, i)
		if len(a)%10 == 0 {
			printAlloc()
			go t(a)
			a = a[:0]
		}
	}
}

func t(list []int) {
	fmt.Println(list)
}

func printAlloc() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB\n", m.Alloc/1024/1024)           // 近期分配的内存
	fmt.Printf("TotalAlloc = %v MiB\n", m.TotalAlloc/1024/1024) // 自上次调用以来分配的总内存
	fmt.Printf("Sys = %v MiB\n", m.Sys/1024/1024)               // 系统为Go分配的物理内存
	fmt.Printf("NumGC = %v\n", m.NumGC)
	fmt.Printf("--------------------------------------------------------------------\n")
}
