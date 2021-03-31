package cluster

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var _ Manager = new(PotManager)

func newPotManager(mode string, reportChan, commandChan chan string) (Manager, error) {
	return &PotManager{
		mode:        mode,
		reportChan:  reportChan,
		commandChan: commandChan,
	}, nil
}

type PotManager struct {
	mode                    string
	reportChan, commandChan chan string

	// 参数区
	nPeer, nSeed    int
	shutdownAtTi    int
	shutdownAtTiMap map[int]int
	cheatAtTiMap    map[int]int

	// 集群管理区
	seeds map[string]Node
	peers map[string]Node
}

// 暂时先以fakeNode代替真实的PotNode
func (p *PotManager) startCluster() error {
	p.seeds = map[string]Node{
		"seed01": newFakeNode("seed01"),
	}
	p.peers = map[string]Node{
		"peer01": newFakeNode("peer01"),
	}
	return nil
}

func (p *PotManager) stopCluster() error {
	return nil
}

func (p *PotManager) readParamsFromCli() error {
	var err error
	_, err = fmt.Scanf("请指定参数（nPeer, nSeed, shutdownAtTi）：%d, %d, %d",
		&p.nPeer, &p.nSeed, &p.shutdownAtTi)
	if err != nil {
		return err
	}

	var satmStr string
	_, err = fmt.Scanf("请指定参数（shutdownAtTiMap）（格式为 PeerNo1:Ti1,PeerNo2:Ti2,...）：%s", &satmStr)
	if err != nil {
		return err
	}
	satmSlice := strings.Split(satmStr, ",")
	if p.shutdownAtTiMap == nil {
		p.shutdownAtTiMap = make(map[int]int)
	}
	var pairStr string
	var pair []string
	var peerno, ti, oldTi int
	for i := 0; i < len(satmSlice); i++ {
		pairStr = satmSlice[i]
		pair = strings.Split(pairStr, ":")
		if len(pair) != 2 {
			return errors.New("invalid shutdownAtTiMap params")
		}
		peerno, err = strconv.Atoi(pair[0])
		if err != nil {
			return errors.New("invalid shutdownAtTiMap params")
		}
		ti, err = strconv.Atoi(pair[1])
		if err != nil {
			return errors.New("invalid shutdownAtTiMap params")
		}
		oldTi = p.shutdownAtTiMap[peerno]
		if (oldTi > 0 && ti < oldTi) || (ti > 0) { // 这两种情况下才用ti，也就是说设置值会取较小的那一个正数
			p.shutdownAtTiMap[peerno] = ti
		}
	}

	var catmStr string
	_, err = fmt.Scanf("请指定参数（cheatAtTiMap）（格式为 PeerNo1:Ti1,PeerNo2:Ti2,...）：%s", &catmStr)
	if err != nil {
		return err
	}
	catmSlice := strings.Split(catmStr, ",")
	if p.cheatAtTiMap == nil {
		p.cheatAtTiMap = make(map[int]int)
	}
	for i := 0; i < len(catmSlice); i++ {
		pairStr = catmSlice[i]
		pair = strings.Split(pairStr, ":")
		if len(pair) != 2 {
			return errors.New("invalid cheatAtTiMap params")
		}
		peerno, err = strconv.Atoi(pair[0])
		if err != nil {
			return errors.New("invalid cheatAtTiMap params")
		}
		ti, err = strconv.Atoi(pair[1])
		if err != nil {
			return errors.New("invalid cheatAtTiMap params")
		}
		oldTi = p.shutdownAtTiMap[peerno]
		if (oldTi > 0 && ti < oldTi) || (ti > 0) { // 这两种情况下才用ti，也就是说设置值会取较小的那一个正数
			p.shutdownAtTiMap[peerno] = ti
		}
	}

	return nil
}

func (p *PotManager) Run() error {
	// 从命令行读取参数到p结构体内
	err := p.readParamsFromCli()
	if err != nil {
		return err
	}

	// 启动集群
	err = p.startCluster()
	if err != nil {
		return err
	}

	return nil
}
