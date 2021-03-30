// package cluster
//
// 运行模式mode: cli/ui
//
// cli模式下不需要reportChan和commandChan，但为了程序调试目的，使能reportChan，并且将report写入一个文本
// 1) 通过命令行交互输入完整所需参数后启动集群管理器，并用集群管理器启动集群
// 2) 集群启动后所有节点的数据报告写到reportChan，由集群管理器负责将报告写到文件中
//
// ui模式下需要这两个chan。
// 1) 通过命令行交互输入参数，启动websocket服务器
// 2) 程序内同时启动ManagerBuilder，用来处理命令
// 3) 客户端UI连接后端服务器 （“连接”）
// 4) 客户端填好集群参数后，全部选项都打包成json传给后端 （“启动集群”）
// 5) 集群运行期间数据报告写道reportChan传给客户端UI
//
// command概念
// command分为两大类：
// 由于使用UI过程中，处理UI命令的后端模块先是ManagerBuilder，后是具体的，比如说PotManager
// ManagerBuilder需要处理
//		ping消息，回复pong消息（用于“连接”）
//      带集群参数的启动消息 （用于”集群启动“）
// 而PotManager用于处理后续的所有指令
// 运行时的commad无非就是关机、作恶两种，关掉某个节点就立刻关闭，作恶就预约接下来最近的时段进行
package cluster
