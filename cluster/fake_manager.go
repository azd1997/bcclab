package cluster

import (
	"fmt"
)

var _ Manager = new(PotManager)

func newFakeManager(mode string, reportChan, commandChan chan string) (Manager, error) {
	return &fakeManager{
		mode:        mode,
		reportChan:  reportChan,
		commandChan: commandChan,
	}, nil
}

// fakeManager 用于测试程序整体的逻辑是否正确
type fakeManager struct {
	mode                    string
	reportChan, commandChan chan string

	// 参数区
	nPeer int

	// 集群管理区
	peers map[string]Node
}

// 暂时先以fakeNode代替真实的PotNode
func (m *fakeManager) startCluster() error {
	m.peers = map[string]Node{
		"peer01": newFakeNode("peer01"),
	}

	if err := m.peers["peer01"].SetReportChan(m.reportChan); err != nil {
		return fmt.Errorf("startCluster: %s", err)
	}
	if err := m.peers["peer01"].Start(); err != nil {
		return err
	}

	// 启动报告线程
	//ticker := time.Tick(time.Second)
	//go func() {
	//	for {
	//		select {
	//		case t := <-ticker:
	//			m.reportChan <- []byte(t.String())
	//		}
	//	}
	//}()

	// 启动一个goroutine去打印节点报告内容
	go func() {
		for data := range m.reportChan {
			fmt.Println("fakeManager: read reportChan: ", string(data))
			//data = data
		}
	}()

	return nil
}

func (m *fakeManager) stopCluster() error {
	return nil
}

func (m *fakeManager) readParamsFromCli() error {
	var err error
	fmt.Print("请指定参数（nPeer）：")
	_, err = fmt.Scanf("%d", &m.nPeer)
	if err != nil {
		return fmt.Errorf("readParamsFromCli:%s", err)
	}
	return nil
}

func (m *fakeManager) Run() error {
	// 从命令行读取参数到p结构体内
	err := m.readParamsFromCli()
	if err != nil {
		return fmt.Errorf("fakeManager.Run: %s", err)
	}

	fmt.Println("参数设定完毕，准备启动集群")

	// 启动集群
	err = m.startCluster()
	if err != nil {
		return fmt.Errorf("fakeManager.Run: %s", err)
	}

	return nil
}
