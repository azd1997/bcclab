package server

import (
	"fmt"
	"github.com/azd1997/bcclab/models"
	"github.com/gorilla/websocket"
	"net/url"
	"testing"
	"time"
)

func TestStartWebSocketServer(t *testing.T) {
	reportChan := make(chan string, 100)
	commandChan := make(chan string, 100)
	addr := "localhost:7777"

	// 启动一个goroutine作报告
	go func() {
		fmt.Println("manager写report")
		ticker := time.Tick(5 * time.Second)
		for {
			select {
			case t := <-ticker:
				record := &models.FakeNodeFaultRecord{
					Time:   t.UnixNano(),
					NodeId: "peer01",
				}
				data, err := record.Json()
				if err != nil {
					panic(err)
				}
				reportChan <- string(data)
				fmt.Println("manager send report: ", string(data))
			}
		}
	}()

	// 启动一个goroutine运行websocketserver
	go func() {
		if err := StartWebSocketServer(addr, reportChan, commandChan); err != nil {
			panic(err)
		}

		// 内部将reportChan内容推给ui

		fmt.Println("server读command")
		// 从commandChan读取命令
		for command := range commandChan {
			fmt.Println("server recv command: ", command)
		}
	}()

	time.Sleep(5 * time.Second) // 等一段时间，保证server已启动，client启动时能直接连上

	// 启动一个goroutine作websocket客户端
	go func() {
		conn, _, err := websocket.DefaultDialer.Dial((&url.URL{
			Scheme: "ws",
			Host:   addr,
			Path:   DefaultServeRouter,
		}).String(), nil)
		if err != nil {
			panic(err)
		}
		//defer conn.Close() 不能让它关闭

		// conn写
		go func() {
			fmt.Println("ui写command")
			// 隔10秒发一个命令
			ticker := time.Tick(10 * time.Second)
			for {
				select {
				case t := <-ticker:
					if err := conn.WriteMessage(websocket.TextMessage, []byte("command "+t.String())); err != nil {
						panic(err)
					} else {
						fmt.Println("ui send command: ", "command "+t.String())
					}

				}
			}
		}()

		// conn读
		go func() {
			fmt.Println("ui读report")
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					panic(err)
				}
				fmt.Println("ui recv report: ", string(msg))
			}
		}()

	}()

	time.Sleep(30 * time.Second)
}
