package logger

import (
	"fmt"
)

type ctxLogger struct {
	buff chan string
	id   any
}

func (me *ctxLogger) E(a ...any) {
	me.buff <- builder("E", me.id, sliceToString(a...))
}
func (me *ctxLogger) Ef(format string, a ...any) {
	me.buff <- builder("E", me.id, fmt.Sprintf(format, a...))
}

func (me *ctxLogger) I(a ...any) {
	me.buff <- builder("I", me.id, sliceToString(a...))
}
func (me *ctxLogger) If(format string, a ...any) {
	me.buff <- builder("I", me.id, fmt.Sprintf(format, a...))
}

func (me *ctxLogger) W(a ...any) {
	me.buff <- builder("W", me.id, sliceToString(a...))
}
func (me *ctxLogger) Wf(format string, a ...any) {
	me.buff <- builder("W", me.id, fmt.Sprintf(format, a...))
}

func (me *ctxLogger) D(a ...any) {
	me.buff <- builder("W", me.id, sliceToString(a...))
}
func (me *ctxLogger) Df(format string, a ...any) {
	me.buff <- builder("W", me.id, fmt.Sprintf(format, a...))
}
