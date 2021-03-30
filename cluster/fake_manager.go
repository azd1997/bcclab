package cluster

import (
	"fmt"
)

var _ Manager = new(PotManager)

func newFakeManager(mode string, reportChan, commandChan chan []byte) (Manager, error) {
	return &PotManager{
		mode:            mode,
		reportChan:      reportChan,
		commandChan:     commandChan,
	}, nil
}

// fakeManager 用于测试程序整体的逻辑是否正确
type fakeManager struct {
	mode string
	reportChan, commandChan chan []byte

	// 参数区
	nPeer int

	// 集群管理区
	peers map[string]Node
}

// 暂时先以fakeNode代替真实的PotNode
func (m *fakeManager) startCluster() error {
	m.peers = map[string]Node{
		"peer01":newFakeNode(),
	}
	return nil
}

func (m *fakeManager) stopCluster() error {
	return nil
}

func (m *fakeManager) readParamsFromCli() error {
	var err error
	fmt.Print("请指定参数（nPeer）：")
	_, err = fmt.Scan("%d", &m.nPeer)
	if err != nil {
		return err
	}
	return nil
}

func (m *fakeManager) Run() error {
	// 从命令行读取参数到p结构体内
	err := m.readParamsFromCli()
	if err != nil {
		return err
	}

	// 启动集群
	err = m.startCluster()
	if err != nil {
		return err
	}

	return nil
}
