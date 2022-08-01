package main

import (
	"fmt"
	"github.com/gookit/color"
	"time"
)

func LogError(method, msg string) {
	color.Error.Printf("[%s] ERROR : [%s] %s\n", time.Now().Format("15:04:05.000"), method, msg)
}

func LogInfo(method, msg string) {
	color.Info.Printf("[%s] INFO : [%s] %s\n", time.Now().Format("15:04:05.000"), method, msg)
}

func LogSuccess(method, msg string) {
	color.Success.Printf("[%s] SUCCESS : [%s] %s\n", time.Now().Format("15:04:05.000"), method, msg)
}

func FormatError(method string, err error) error {
	return fmt.Errorf("[%s] : %s", method, err)
}
