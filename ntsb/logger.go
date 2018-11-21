package ntsb

import "fmt"

type Logger struct{}

func NewLogger() Logger {
	return Logger{}
}

func (Logger) Println(a ...interface{}) {
	fmt.Println(a...)
}
