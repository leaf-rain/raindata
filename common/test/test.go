package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	// 打开文件
	file, err := os.Open("./test/FirstChargeDialog.prefab")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 创建一个带缓冲的 Scanner，用于按行读取
	scanner := bufio.NewScanner(file)

	// 设置扫描器的分隔符为换行符
	scanner.Split(bufio.ScanLines)

	// 遍历文件的每一行
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	// 检查读取过程中是否出现了错误
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
