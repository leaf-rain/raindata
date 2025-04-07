package main

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var (
	ErrMsgType = errors.New("不支持的消息类型")
)

type Handel func(ws ws, code uint16, msg string)

type connMap struct {
	sync.Map
	handel Handel
}

func newCommMap(handel Handel) *connMap {
	return &connMap{
		Map:    sync.Map{},
		handel: handel,
	}
}

func (c *connMap) getAllConn() []ws {
	var result []ws
	c.Range(func(key, value interface{}) bool {
		if item, ok := value.(ws); ok {
			result = append(result, item)
		}
		return true
	})
	return result
}

type ws struct {
	Id      int64
	Conn    *websocket.Conn
	MsgLock *sync.Mutex
}

func (ws ws) SendMsg(code uint16, msg interface{}) error {
	var body []byte
	switch msg.(type) {
	case string:
		body = []byte(msg.(string))
	case []byte:
		body = msg.([]byte)
	default:
		return ErrMsgType
	}
	var length = len(body)
	head, headLength := uint16ToBytesLittleEndian(code)
	var sendBody = make([]byte, length+headLength)
	copy(sendBody[:headLength], head)
	copy(sendBody[headLength:], body)
	ws.MsgLock.Lock()
	defer ws.MsgLock.Unlock()
	err := ws.Conn.WriteMessage(websocket.BinaryMessage, sendBody)
	if err != nil {
		log.Printf("msg send failed. err:%v\n", err)
	}
	return err
}

var (
	ErrUserOffline = errors.New("用户离线")
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// todo:鉴权
		return true
	},
}

func (c *connMap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		http.Error(w, "WebSocket upgrade error", http.StatusInternalServerError)
		return
	}
	// 添加到map
	od, err := c.shakeHands(conn)
	if err != nil {
		log.Println("Handshake error:", err)
		conn.Close()
		return
	}
	go c.receiveMsg(od)
	w.WriteHeader(http.StatusOK)
}

func (c *connMap) shakeHands(conn *websocket.Conn) (ws, error) {
	var id = SnowflakeInt64()
	od := ws{
		Id:      id,
		Conn:    conn,
		MsgLock: new(sync.Mutex),
	}
	c.Store(id, od)
	log.Printf("New connection established with ID: %d, conn:%s\n", id, conn.RemoteAddr().String())
	return od, nil
}

func (c *connMap) receiveMsg(od ws) {
	defer func() {
		c.Delete(od.Id)
		od.Conn.Close()
		log.Printf("Connection closed for ID: %d, addr:%s\n", od.Id, od.Conn.RemoteAddr().String())
	}()
	for {
		_, message, err := od.Conn.ReadMessage()
		if err != nil || len(message) < 2 {
			log.Printf("Read error for ID %d: %v\n", od.Id, err)
			return
		}
		log.Printf("Received message from ID %d: %s\n", od.Id, message)
		head, err := bytesToUint16LittleEndian(message[:2])
		if err != nil || head == 0 {
			log.Printf("Read error for ID %d: %v\n", od.Id, err)
			continue
		}
		body := string(message[2:])
		c.handel(od, head, body)
	}
}
