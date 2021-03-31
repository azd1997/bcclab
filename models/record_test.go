package models

import (
	"testing"
	"time"
)

func TestJson(t *testing.T) {
	a := FakeNodeFaultRecord{
		Time:   time.Now().UnixNano(),
		NodeId: "peer01",
	}
	data, err := a.Json()
	if err != nil {
		t.Error(err)
	}
	t.Log(a)
	t.Log(string(data))
}
