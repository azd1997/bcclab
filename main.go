package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/azd1997/bcclab/cluster"
	"github.com/azd1997/bcclab/server"
)

const (
	ModeCli string = "cli"
	ModeUi string = "ui"
)


var (
	flagMode = flag.String("mode", "cli", "specify running mode (cli or ui)")

	reportChan = make(chan []byte, 100)
	commandChan = make(chan []byte, 100)
)

// main函数逻辑：启动websocket服务器 -> 启动节点管理器
func main()  {

	// 解析命令行程序
	flag.Parse()

	if *flagMode == ModeCli {
		runCliMode()
	} else if *flagMode == ModeUi {
		runUiMode()
	} else {
		printUsage()
	}
}

func runCliMode() {
	// 进行终端界面下的指令输入，由于不同集群输入的参数不尽相同，交给Manager去实现

	// 1. 读取共识协议名称
	var consensus string
	fmt.Print("请指定协议类型：")
	_, err := fmt.Scan(&consensus)
	handleError(err)

	// 2. 创建对应协议的Manager，并启动Manager的参数输入
	err = cluster.StartManager(consensus, ModeCli, reportChan, commandChan)
	handleError(err)
}

func runUiMode() {
	// 1. 指定本机监听地址
	var addr string
	_, err := fmt.Scanf("请输入服务器监听地址（置空则默认为localhost:7777）：%s", &addr)
	handleError(err)

	// 2. 启动websocket服务器
	err = server.StartWebSocketServer(addr, reportChan, commandChan)
	handleError(err)

	// 3. 创建对应协议的ManagerBuilder。mb
	err = cluster.StartManagerBuilder(reportChan, commandChan)
	handleError(err)
}

func printUsage() {
	fmt.Print(`
BccLab —— 区块链共识协议实验平台
Usage: 根据交互提示输入小写协议名称 （pot/raft/...）再根据提示输入实验参数即可
`)
}

// 打印错误、打印用法、退出程序
func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		printUsage()
		os.Exit(1)
	}
}