package output

import (
	"fmt"

	"github.com/fatih/color"
)

var isQuiet = true

func SetOutput(mute bool) {
	isQuiet = mute
}

func Print(s string, a ...interface{}) {
	if !isQuiet {
		fmt.Printf(s+"\n", a...)
		return
	}
}

func Error(s string, a ...interface{}) {
	if !isQuiet {
		color.Red(s, a...)
		return
	}
}

func Success(s string, a ...interface{}) {
	if !isQuiet {
		color.Green(s, a...)
		return
	}
}

func Warning(s string, a ...interface{}) {
	if !isQuiet {
		color.Yellow(s, a...)
		return
	}
}

func AlwaysPrint(s string, a ...interface{}) {
	fmt.Printf(s+"\n", a...)
}
