package output

import "github.com/fatih/color"

func ErrorWithColor(s string, a ...interface{}) {
	color.Red(s, a...)
}

func SuccessWithColor(s string, a ...interface{}) {
	color.Green(s, a...)
}

func WarningWithColor(s string, a ...interface{}) {
	color.Yellow(s, a...)
}
