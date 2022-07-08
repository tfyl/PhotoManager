package main

import (
	"fmt"
	"github.com/fatih/color"
)

func LogError(method, msg string) {
	color.Red("ERROR : [%s] %s", method, msg)
}

func LogInfo(method, msg string) {
	color.Blue("INFO : [%s] %s", method, msg)
}

func LogSuccess(method, msg string) {
	color.Green("SUCCESS : [%s] %s", method, msg)
}

func FormatError(method string, err error) error {
	return fmt.Errorf("%s : %s", method, err)
}
