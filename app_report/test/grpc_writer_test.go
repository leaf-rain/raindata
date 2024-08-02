package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/leaf-rain/raindata/app_report/api/grpc"
	"github.com/leaf-rain/raindata/common/consts"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestGrpcWriter(t *testing.T) {
	// 1.建立连接 获取client
	conn, err := grpc.Dial("192.168.10.28:8090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewLogServerClient(conn)
	callBidirectionalStream(client)
}
func callBidirectionalStream(client pb.LogServerClient) {
	stream, err := client.StreamReport(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		var err error
		for {
			var msg pb.StreamReportResponse
			err = stream.RecvMsg(&msg)
			if err != nil {
				log.Printf("Receive error: %v", err)
				break
			}
			log.Printf("Received: %s", msg.String())
		}
		close(waitc)
	}()

	for i := 0; i < 10; i++ {
		message := fmt.Sprintf("Message %d", i)
		req := &pb.StreamReportRequest{Message: getRandomJson()}
		if err = stream.Send(req); err != nil {
			log.Fatalf("Send error: %v", err)
		}
		log.Printf("Sent message: %s", message)
		time.Sleep(1 * time.Second)
	}
	stream.CloseSend()
	<-waitc
}

func getRandomJson() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var length = rand.Intn(50)
	var data = make(map[string]interface{}, length)
	data[consts.KeyAppidForMsg] = 1
	data[consts.KeyEventForMsg] = generateRandomString(4)
	for i := 0; i < length; i++ {
		data[generateRandomString(4)] = generateRandomString(4)
	}
	result, _ := json.Marshal(data)
	return string(result)
}

func generateRandomString(length int) string {
	// 设置随机数种子
	rand.New(rand.NewSource(time.Now().UnixNano()))
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result string
	for i := 0; i < length; i++ {
		result += string(charset[rand.Intn(len(charset))])
	}
	return result
}
