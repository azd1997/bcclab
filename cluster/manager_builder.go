package cluster

import "log"

var mbInstance *ManagerBuilder

// BasicManager 基础Manager，不负责任何协议集群的工作
// 职责是启动消息处理循环，之后再根据前端传来的指令生成具体的Manager
type ManagerBuilder struct {
	reportChan  chan string
	commandChan chan string

	done chan struct{}
}

func (mb *ManagerBuilder) instantiateManager(consensus string) {

}

func (mb *ManagerBuilder) handleCommand() error {
	// 解析command

	// 处理command

	return nil
}

func (mb *ManagerBuilder) handleCommandLoop() {
	for {
		select {
		case <-mb.done:
			return
		case command := <-mb.commandChan:
			if err := mb.handleCommand(); err != nil {
				log.Printf("handle command fail. err=%s, command=%s\n", err, string(command))
			}
		}
	}
}

func StartManagerBuilder(reportChan, commandChan chan string) error {
	if mbInstance == nil {
		mbInstance = &ManagerBuilder{reportChan: reportChan, commandChan: commandChan, done: make(chan struct{})}
	}
	go mbInstance.handleCommandLoop()
	return nil
}
