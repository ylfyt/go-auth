package l

import "fmt"

var logChan chan string = make(chan string)

func worker() {
	for msg := range logChan {
		fmt.Println(msg)
	}
}

func init() {
	for i := 0; i < 2; i++ {
		go worker()
	}
}

func I(format string, a ...any) {
	logChan <- "I|" + fmt.Sprintf(format, a...)
}

func E(format string, a ...any) {
	logChan <- "E|" + fmt.Sprintf(format, a...)
}

func W(format string, a ...any) {
	logChan <- "W|" + fmt.Sprintf(format, a...)
}

func D(format string, a ...any) {
	logChan <- "D|" + fmt.Sprintf(format, a...)
}
