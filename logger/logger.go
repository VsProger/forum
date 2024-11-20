package logger

import (
	"fmt"
)

type Logger interface {
	Info(msg string)
	Error(msg string, err error)
}

type logger struct {
	info string
	e    string
}

func NewLogger() Logger {
	return &logger{
		info: "INFO",
		e:    "ERROR",
	}
}

func (l *logger) Info(msg string) {
	fmt.Println(colorize("INFO:", "\033[0;34m"), msg)
}

func (l *logger) Error(msg string, err error) {
	fmt.Println(colorize("ERROR:", "\033[0;31m"), msg, err)
}

func colorize(prefix, color string) string {
	reset := "\033[0m"
	return fmt.Sprintf("%s%s%s", color, prefix, reset)
}
