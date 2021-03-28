package cluster

import "github.com/azd1997/bcclab/report"

type Manager interface {
	SetCommandChan(chan string)
	SetReporter(r report.Reporter)

	StartCluster()
	StopCluster()


}

func NewManager(consensus string) Manager {

}
