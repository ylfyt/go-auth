package logger

import (
	"fmt"
	"io"
)

type Logger struct {
	f    io.Writer
	buff chan string
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
		id:   id,
		buff: me.buff,
	}
}

func (me *Logger) E(a ...any) {
	me.buff <- builder("E", nil, sliceToString(a...))
}

func (me *Logger) Ef(format string, a ...any) {
	me.buff <- builder("E", nil, fmt.Sprintf(format, a...))
}
func (me *Logger) I(a ...any) {
	me.buff <- builder("I", nil, sliceToString(a...))
}

func (me *Logger) If(format string, a ...any) {
	me.buff <- builder("I", nil, fmt.Sprintf(format, a...))
}

func (me *Logger) W(a ...any) {
	me.buff <- builder("W", nil, sliceToString(a...))
}
func (me *Logger) Wf(format string, a ...any) {
	me.buff <- builder("W", nil, fmt.Sprintf(format, a...))
}

func (me *Logger) D(a ...any) {
	me.buff <- builder("D", nil, sliceToString(a...))
}
func (me *Logger) Df(format string, a ...any) {
	me.buff <- builder("D", nil, fmt.Sprintf(format, a...))
}
