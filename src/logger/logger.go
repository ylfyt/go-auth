package logger

import (
	"fmt"
	"io"
)

type LoggerLevel int

const (
	LOG_DEBUG LoggerLevel = 0
	LOG_INFO  LoggerLevel = 1
)

var levelMap = map[string]LoggerLevel{
	"D": 0,
	"I": 1,
	"W": 2,
	"E": 3,
}

type Logger struct {
	f        io.Writer
	buff     chan string
	maxLevel LoggerLevel
}

type ILogger interface {
	E(a ...any)
	Ef(format string, a ...any)
	I(a ...any)
	If(format string, a ...any)
	W(a ...any)
	Wf(format string, a ...any)
	D(a ...any)
	Df(format string, a ...any)
}

func (me *Logger) run() {
	for {
		msg := <-me.buff
		fmt.Fprintln(me.f, msg)
	}
}

func (me *Logger) R(id any) ILogger {
	return &ctxLogger{
		id:       id,
		buff:     me.buff,
		maxLevel: me.maxLevel,
	}
}

func (me *Logger) E(a ...any) {
	me.append("E", sliceToString(a...))
}

func (me *Logger) Ef(format string, a ...any) {
	me.append("E", fmt.Sprintf(format, a...))
}
func (me *Logger) I(a ...any) {
	me.append("I", sliceToString(a...))
}

func (me *Logger) If(format string, a ...any) {
	me.append("I", fmt.Sprintf(format, a...))
}

func (me *Logger) W(a ...any) {
	me.append("W", sliceToString(a...))
}
func (me *Logger) Wf(format string, a ...any) {
	me.append("W", fmt.Sprintf(format, a...))
}

func (me *Logger) D(a ...any) {
	me.append("D", sliceToString(a...))
}
func (me *Logger) Df(format string, a ...any) {
	me.append("D", fmt.Sprintf(format, a...))
}

func (me *Logger) append(level, message string) {
	if levelMap[level] < me.maxLevel {
		return
	}
	me.buff <- builder(level, nil, message)
}
