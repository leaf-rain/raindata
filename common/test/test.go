package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	fileSize := 10 * 1024 // 10MB
	fileName := "myfile.text"
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	if err := os.Truncate(f.Name(), int64(fileSize)); err != nil {
		panic(err)
	}
	fmt.Println(getFileSize(fileName))
	// 可选：写入一些数据到文件开头
	if _, err := f.Write([]byte("Hello, World!\n")); err != nil {
		panic(err)
	}
	if _, err := f.Write([]byte("Hello, World!\n")); err != nil {
		panic(err)
	}
	if _, err := f.Write([]byte("Hello, World!\n")); err != nil {
		panic(err)
	}
	f.Close()
	fmt.Println(getFileSize(fileName))

	f, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	if err := os.Truncate(f.Name(), int64(fileSize)); err != nil {
		panic(err)
	}
	fmt.Println(getFileSize(fileName))
	// 可选：写入一些数据到文件开头
	if _, err := f.Write([]byte("Hello, World!\n")); err != nil {
		panic(err)
	}
	if _, err := f.Write([]byte("Hello, World!\n")); err != nil {
		panic(err)
	}
	if _, err := f.Write([]byte("Hello, World!\n")); err != nil {
		panic(err)
	}
	f.Close()
	fmt.Println(getFileSize(fileName))
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

func getFileSize(filePath string) (string, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}
	fileSize := fileInfo.Size()
	return humanReadableByteCount(fileSize, false), nil
}

func humanReadableByteCount(b int64, si bool) string {
	var unit string = "iB"
	if si {
		unit = "B"
	}
	units := []string{"K" + unit, "M" + unit, "G" + unit, "T" + unit, "P" + unit, "E" + unit, "Z" + unit}
	if b < 1024 {
		return fmt.Sprintf("%d B", b)
	}
	magnitude := int64(0)
	for ; b >= 1024; b /= 1024 {
		magnitude++
	}
	return fmt.Sprintf("%.1f %s", float64(b)/1024.0, units[magnitude-1])
}
