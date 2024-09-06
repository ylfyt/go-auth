package logger

import (
	"fmt"
)

type ctxLogger struct {
	buff     chan string
	id       any
	maxLevel LoggerLevel
}

func (me *ctxLogger) E(a ...any) {
	me.append("E", sliceToString(a...))
}
func (me *ctxLogger) Ef(format string, a ...any) {
	me.append("E", fmt.Sprintf(format, a...))
}

func (me *ctxLogger) I(a ...any) {
	me.append("I", sliceToString(a...))
}
func (me *ctxLogger) If(format string, a ...any) {
	me.append("I", fmt.Sprintf(format, a...))
}

func (me *ctxLogger) W(a ...any) {
	me.append("W", sliceToString(a...))
}
func (me *ctxLogger) Wf(format string, a ...any) {
	me.append("W", fmt.Sprintf(format, a...))
}

func (me *ctxLogger) D(a ...any) {
	me.append("W", sliceToString(a...))
}
func (me *ctxLogger) Df(format string, a ...any) {
	me.append("W", fmt.Sprintf(format, a...))
}

func (me *ctxLogger) append(level, message string) {
	if levelMap[level] > me.maxLevel {
		return
	}
	me.buff <- builder(level, me.id, message)
}
