package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"log"
	"runtime/debug"
)

const (
	SERVER_RESOURCE_REPORTING  uint16 = 2
	CLIENT_EXEC_COMMAND        uint16 = 3
	SERVER_EXEC_COMMAND        uint16 = 4
	CLIENT_CANCEL_EXEC_COMMAND uint16 = 5
	SERVER_CANCEL_EXEC_COMMAND uint16 = 6
)

func uint16ToBytesLittleEndian(num uint16) ([]byte, int) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, num)
	if err != nil {
		log.Println("binary.Write failed:", err)
		return nil, 0
	}
	return buf.Bytes(), len(buf.Bytes())
}
func bytesToUint16LittleEndian(byteArray []byte) (uint16, error) {
	if len(byteArray) != 2 {
		return 0, fmt.Errorf("byte array must be exactly 2 bytes long")
	}

	var num uint16
	reader := bytes.NewReader(byteArray)
	err := binary.Read(reader, binary.LittleEndian, &num)
	if err != nil {
		return 0, err
	}
	return num, nil
}

const (
	maxNode = 1023
)

func GetSnowflakeId() snowflake.ID {
	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(maxNode)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return node.Generate()
}

func SnowflakeInt64() int64 {
	return GetSnowflakeId().Int64()
}

func GoSafe(f func()) {
	go func() {
		defer Recover()
		f()
	}()
}

func Recover() {
	if err := recover(); err != nil {
		log.Printf("panic recover, err: %v", err)
		debug.PrintStack()
	}
}
