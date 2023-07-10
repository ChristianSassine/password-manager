package output

import (
	"fmt"

	"github.com/fatih/color"
)

var withColor = true

func SetOutput(color bool) {
	withColor = color
}

func NormalLn(s string, a ...interface{}) {
	fmt.Printf(s+"\n", a...)
}

func Error(s string, a ...interface{}) {
	if withColor {
		color.Red(s, a...)
		return
	}
	fmt.Printf(s, a...)
}

func Success(s string, a ...interface{}) {
	if withColor {
		color.Green(s, a...)
		return
	}
	fmt.Printf(s, a...)
}

func Warning(s string, a ...interface{}) {
	if withColor {
		color.Yellow(s, a...)
		return
	}
	fmt.Printf(s, a...)
}
