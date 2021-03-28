package report

type Reporter interface {
	SetReportChan(ch chan *string)
	Report()
}
