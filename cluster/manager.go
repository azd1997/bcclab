package cluster

import (
	"fmt"
)

var mInstance Manager

func StartManager(consensus string, mode string, reportChan, commandChan chan []byte) error {
	switch consensus {
	case "pot":
		pm, err := newPotManager(mode, reportChan, commandChan)
		if err != nil {
			return err
		}
		mInstance = pm
		return mInstance.Run()
	default:
		return fmt.Errorf("unknown consensus: %s", consensus)
	}
}

type Manager interface {
	//SetRunMode(mode string)
	//RunMode() string	// 返回 "ui", "cli"

	//SetCommandChan(chan []byte)
	//SetReportChan(chan []byte)

	startCluster() error
	stopCluster() error
	readParamsFromCli() error

	// Run 外部调用，先从命令行获取参数，后执行startCluster
	Run() error

	//// Cli 每种Manager自定义的终端输入参数的方法
	//CliAndStartCluster() error
}



