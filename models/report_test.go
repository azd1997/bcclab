package models

import (
	"fmt"
	"testing"
	"time"
)

// channel的for-range用法: for a := range achan {}

// 这种做法可行
func TestReportWithChanString1(t *testing.T) {
	text := "hello"
	ch := make(chan string, 100)

	go func() {
		ch <- text
	}()

	go func() {
		atext := <-ch
		fmt.Println(atext)
	}()
}

func TestReportWithChanBytes(t *testing.T) {
	text := "hello"
	ch := make(chan []byte, 100)

	go func() {
		ch <- []byte(text)
	}()

	go func() {
		atext := <-ch
		fmt.Println(string(atext))
	}()

	time.Sleep(time.Second)
}
