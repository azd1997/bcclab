// package server websocket服务器
package server

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	DefaultServeAddr   = "localhost:7777"
	DefaultServeRouter = "/ws"
)

var (
	upgrader = websocket.Upgrader{} // 使用默认参数
)

// websocketServer websocket服务器
type websocketServer struct {
	serveMux http.ServeMux
	server   *http.Server
	router   string // 监听路由

	wsconn *websocket.Conn // 需要注意的是该Conn只支持单读单写。不能并发写或者并发读
	wlock  sync.Mutex      // 因为有两个地方需要写消息，所以整个写锁，读锁就不整了

	reportChan  chan string // node的报告都塞到这个chan
	commandChan chan string // 前端UI给的指令都写到这个chan

	done chan struct{} // 关闭通知
}

// 报告循环主动关闭；处理循环被动关闭（前端点击关闭后，前端先关闭连接，服务端被动关闭）
func (ws *websocketServer) stop() {
	ws.done <- struct{}{}
}

// 将请求中携带的文字内容传出
func (ws *websocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, http.Header{
		"protocol": []string{"bcclab-json"},
	})
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	ws.wsconn = c

	for {
		// 虽然有ReadJSON方法，但是不建议在这里解析，server只做通信的事
		mt, message, err := c.ReadMessage() // messgae为json文本消息
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		// 要求消息类型必须是Text
		if mt == websocket.TextMessage {
			// 将消息写入到commandChan中
			ws.commandChan <- string(message)
		}

		// 利用pong回复已收到消息，但命令消息的执行结果没法直接回复
		// 1是因为丢给chan
		// 2是因为有些命令的执行结果需要一段时间，不可能在这里阻塞它
		err = ws.write(websocket.PongMessage, nil)
		if err != nil {
			log.Println("pong:", err)
			break
		}
	}
}

func (ws *websocketServer) write(mt int, msg []byte) error {
	ws.wlock.Lock()
	ws.wlock.Unlock()
	return ws.wsconn.WriteMessage(mt, msg)
}

func (ws *websocketServer) reportLoop() {

	for {
		select {
		case <-ws.done:
			log.Println("report loop stopped")
			return
		case msg := <-ws.reportChan:
			err := ws.write(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}

func StartWebSocketServer(addr string, reportChan chan string, commandChan chan string) error {
	// 首先初始化好websocketserver实例
	if strings.TrimSpace(addr) == "" {
		addr = DefaultServeAddr
	}
	if reportChan == nil || commandChan == nil {
		return errors.New("require initialized reportChan and commandChan")
	}
	ws := &websocketServer{
		router:      DefaultServeRouter,
		reportChan:  reportChan,
		commandChan: commandChan,
		done:        make(chan struct{}),
	}
	ws.serveMux.Handle(ws.router, ws) // 之所以使用"/ws"是怕以后有其他需求
	ws.server = &http.Server{Addr: addr, Handler: &ws.serveMux}

	// 2. 启动websocket监听
	if err := ws.server.ListenAndServe(); err != nil {
		return err
	}

	// 3. 启动websocketserver自身处理reportChan循环
	go ws.reportLoop()
	return nil
}
