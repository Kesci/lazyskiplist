package lazyskiplist

import (
	"fmt"
)

var debug bool

// TurnOnDebug turns on debug log for developing
func TurnOnDebug() {
	debug = true
}

func infof(format string, a ...interface{}) {
	if len(format) != 0 && format[len(format)-1] != '\n' {
		format += "\n"
	}
	fmt.Printf(format, a...)
}

func debugf(format string, a ...interface{}) {
	if debug {
		infof(format, a...)
	}
}
