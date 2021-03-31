package server

import (
	"sync"

	"github.com/gorilla/websocket"
)

// 由于程序中可能并发的goroutine向conn进行读写，所以封装一下
// 读和写可同时，但读之间互斥、写之间互斥
type wsConn struct {
	conn *websocket.Conn
	rlock sync.Mutex
	wlock sync.Mutex
}

func (c *wsConn) ReadMessage() (messageType int, p []byte, err error) {
	c.rlock.Lock()
	defer c.rlock.Unlock()
	return c.conn.ReadMessage()
}

func (c *wsConn) WriteMessage(messageType int, data []byte) error {
	c.wlock.Lock()
	defer c.wlock.Unlock()
	return c.conn.WriteMessage(messageType, data)
}

