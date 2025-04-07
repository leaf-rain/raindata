package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"testing"
	"time"
)

func Test_Main(t *testing.T) {
	serverURL := "ws://localhost:8080/ws"
	// 解析WebSocket URL
	u, err := url.Parse(serverURL)
	if err != nil {
		log.Fatal("无法解析URL:", err)
	}

	// 设置WebSocket连接的Dialer
	dialer := websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// 连接到WebSocket服务器
	conn, resp, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("无法连接到WebSocket服务器: %v, 响应: %v", err, resp)
	}
	defer conn.Close()

	log.Println("已连接到WebSocket服务器")

	// 创建一个通道来接收系统信号（如Ctrl+C）
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// 启动一个goroutine来发送消息
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				id := "test"
				od := `
for((i=0;i<10;i++))
do
    sleep 1
    echo $(date +"%Y-%m-%d %H:%M:%S")
done
`
				req := map[string]interface{}{
					"id":   id,
					"data": od,
				}
				message, _ := json.Marshal(req)
				body, _ := uint16ToBytesLittleEndian(CLIENT_EXEC_COMMAND)
				body = append(body, message...)
				log.Printf("发送消息: %s\n", string(message))
				err := conn.WriteMessage(websocket.TextMessage, body)
				if err != nil {
					log.Printf("发送消息失败: %v\n", err)
					return
				}
				time.Sleep(time.Second * 3)
				body, _ = uint16ToBytesLittleEndian(CLIENT_CANCEL_EXEC_COMMAND)
				body = append(body, []byte(id)...)
				log.Printf("发送消息: %s\n", string(body))
				err = conn.WriteMessage(websocket.TextMessage, body)
				if err != nil {
					log.Printf("发送消息失败: %v\n", err)
					return
				}
			case <-interrupt:
				log.Println("接收到中断信号，关闭连接")
				err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Printf("关闭连接失败: %v\n", err)
				}
				return
			}
		}
	}()

	// 启动一个goroutine来接收消息
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("接收消息失败: %v\n", err)
				return
			}
			log.Printf("接收消息: %s\n", message)
		}
	}()

	// 阻塞主goroutine，直到接收到中断信号
	<-interrupt
	log.Println("程序退出")
}
