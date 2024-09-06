package logger

import (
	"fmt"
	"io"
	"runtime"
	"time"
)

func NewLogger(f io.Writer, level LoggerLevel) *Logger {
	l := &Logger{
		buff:     make(chan string),
		f:        f,
		maxLevel: level,
	}
	go l.run()

	return l
}

func sliceToString[T any](a ...T) string {
	message := ""
	for i, e := range a {
		message += fmt.Sprint(e)
		if i != len(a)-1 {
			message += " "
		}
	}
	return message
}

func builder(level string, id any, message string) string {
	pc, _, line, _ := runtime.Caller(3)
	funcCaller := runtime.FuncForPC(pc)
	funcName := ""
	if funcCaller != nil {
		funcName = funcCaller.Name()
	}
	idStr := "-"
	if id != nil {
		idStr = fmt.Sprint(id)
	}
	return fmt.Sprintf("%s|%s|%s|%s|%s", time.Now().Format("2006-01-02 15:04:05.000"), level, idStr, fmt.Sprintf("%s():%d", funcName, line), message)
}
