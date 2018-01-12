package lazyskiplist

import (
	"fmt"
)

var debug bool

func TurnOnDebug() {
	debug = true
}

func Infof(format string, a ...interface{}) {
	if len(format) != 0 && format[len(format)-1] != '\n' {
		format += "\n"
	}
	fmt.Printf(format, a...)
}

func Debugf(format string, a ...interface{}) {
	if !debug {
		Infof(format, a...)
	}
}
