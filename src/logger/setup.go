package logger

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

func NewLogger(directory string, filename string, level LoggerLevel, stdout bool) *Logger {

	writer := os.Stdout
	if !stdout {
		f, err := getLogFile(directory, filename)
		if err != nil {
			panic(err)
		}
		writer = f
	}

	l := &Logger{
		buff:      make(chan string),
		f:         writer,
		maxLevel:  level,
		ref:       time.Date(2024, 9, 18, 0, 0, 0, 0, time.UTC),
		directory: directory,
		filename:  filename,
		stdout:    stdout,
	}
	go l.run()

	return l
}

func getLogFile(dir string, currName string) (*os.File, error) {
	file, err := os.OpenFile(fmt.Sprintf("%s/%s.log", dir, currName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		if scanner.Err() == nil {
			return file, nil
		}
		file.Close()
		return nil, err
	}
	firstLine := scanner.Text()

	insertedAt, err := time.Parse("2006-01-02 15:04:05.000", strings.Split(firstLine, "|")[0])
	if err != nil {
		file.Close()
		return nil, err
	}
	now := time.Now()
	t1 := time.Date(insertedAt.Year(), insertedAt.Month(), insertedAt.Day(), 0, 0, 0, 0, insertedAt.Location())
	t2 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if t1.Sub(t2) >= 0 {
		return file, nil
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}
	err = os.Rename(fmt.Sprintf("%s/%s.log", dir, currName), fmt.Sprintf("%s/%s_%s.log", dir, currName, t1.Format("2006-01-02")))
	if err != nil {
		return nil, err
	}
	return getLogFile(dir, currName)

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
