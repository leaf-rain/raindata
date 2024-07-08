package main

import (
	"fmt"
	"runtime"
)

func main() {
	a := "中国"
	println(len(a))
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
