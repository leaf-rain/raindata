package main

import (
	"fmt"
	"io"
	"runtime"

	"github.com/rosedblabs/wal"
)

func main() {
	var opt = wal.DefaultOptions
	opt.DirPath = "./"
	wal, _ := wal.Open(opt)
	// write some data
	chunkPosition, _ := wal.Write([]byte("some data 1"))
	// read by the position
	val, _ := wal.Read(chunkPosition)
	fmt.Println(string(val))

	wal.Write([]byte("some data 2"))
	wal.Write([]byte("some data 3"))

	// iterate all data in wal
	reader := wal.NewReader()
	for {
		val, pos, err := reader.Next()
		if err == io.EOF {
			break
		}
		fmt.Println(string(val))
		fmt.Println(pos) // get position of the data for next read
		fmt.Println(reader.CurrentChunkPosition())
		fmt.Println(reader.CurrentSegmentId())
		fmt.Println("-----------------")
	}
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
