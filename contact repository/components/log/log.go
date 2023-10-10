package log

import "fmt"

type Logger interface {
	Print(value ...string)
	New()
}
type Log struct {
}

func GetLogger() *Log {
	return &Log{}
}

func (l *Log) Print(value ...interface{}) {
	//code to write logs in a file etc...
	fmt.Println(value) //wrapper
	return

}
