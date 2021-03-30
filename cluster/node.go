package cluster

import (
	"fmt"
	"github.com/azd1997/bcclab/models"
	"time"
)

// Node 对集群管理器管理的Node做要求
type Node interface {
	SetReportChan(reportChan chan []byte) error

	Start() error
	Stop() error
	Fault(faulttype string) error	// 预约最近的一段时间作恶
}
// 后续还可以添加一些细致性的属性设置、动作设定



type fakeNode struct {
	reportChan chan []byte

	done chan struct{}
}

func newFakeNode() *fakeNode {
	return &fakeNode{
		reportChan: nil,
		done:       make(chan struct{}),
	}
}

func (f *fakeNode) Start() error {
	go f.loop()
	return nil
}



func (f *fakeNode) Stop() error {
	f.done <- struct{}{}
	return nil
}

// 定时报告一些内容
func (f *fakeNode) loop() {
	ticker := time.Tick(time.Second)
	for {
		select {
		case <- f.done:
			return
		case t := <- ticker:
			f.report(models.PotNodeStateSwitchRecord{
				Time:     t.UnixNano(),
				NodeId:   "peer01",
				NewState: "xiatian",
			})
		}
	}
}

// report内容本质上是不定的，想要通过统一格式规定不太现实
// 只能利用json存储，由后端、前端同时解释
func (f *fakeNode) report(record models.Record) error {
	data, err := record.Json()
	if err != nil {
		return err
	}
	f.reportChan <- data
	return nil
}

//
//
//// 1. 何时(Time) 某个节点(Who) 状态切换(Act) 成新状态(ActTo)
//// 2. 何时(Time) 某个节点(Who) 发信(Act) 给另一个节点(ActTo)
//// 3.
//type Record struct {
//	ID string	// Record哈希用作标识
//
//	Time int64	// 时间
//	Who string	// 节点id
//	Act string	// 发信 收信
//	ActTo string	// （另一个节点）
//}
//
//type Record map[string]string



// 标准的实现应该用once保证只调用一次
func (f *fakeNode) SetReportChan(reportChan chan []byte) error {
	f.reportChan = reportChan
	return nil
}

func (f *fakeNode) Fault(faulttype string) error {
	fmt.Printf("fault: %s\n", faulttype)
	return f.report(models.FakeNodeFaultRecord{
		Time:   time.Now().UnixNano(),
		NodeId: "peer01",
	})
}

