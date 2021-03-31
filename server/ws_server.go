// package server websocket服务器
package server

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
)

const (
	DefaultServeAddr   = "localhost:7777"
	DefaultServeRouter = "/ws"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {	// 对访问的主机一概放行
			return true
		},
	} // 使用默认参数
)

// websocketServer websocket服务器
// 限制连接数量为1
type websocketServer struct {
	serveMux http.ServeMux
	server   *http.Server
	router   string // 监听路由

	wsconn *wsConn

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
	// 只允许一个连接
	if ws.wsconn != nil {
		w.WriteHeader(404)
		w.Write([]byte("BccLab server只允许一个客户端连接"))
	}

	c, err := upgrader.Upgrade(w, r, http.Header{
		"protocol": []string{"bcclab-json"},
	})
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	log.Println("ui conn: ", c.RemoteAddr())
	defer c.Close()
	if ws.wsconn == nil {
		ws.wsconn = &wsConn{conn: c}
	}

	for {
		// 虽然有ReadJSON方法，但是不建议在这里解析，server只做通信的事
		mt, message, err := ws.wsconn.ReadMessage() // messgae为json文本消息
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
		err = ws.wsconn.WriteMessage(websocket.PongMessage, nil)
		if err != nil {
			log.Println("pong:", err)
			break
		}
	}
}

func (ws *websocketServer) reportLoop() {

	for {
		select {
		case <-ws.done:
			log.Println("report loop stopped")
			return
		case msg := <-ws.reportChan:
			if ws.wsconn == nil {continue}	// 还未建立连接则丢弃报告
			err := ws.wsconn.WriteMessage(websocket.TextMessage, []byte(msg))
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
	go func() {
		if err := ws.server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	// 3. 启动websocketserver自身处理reportChan循环
	go ws.reportLoop()
	return nil
}
