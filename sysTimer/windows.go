//go:build windows

package sysTimer

import (
	"syscall"
	"time"
)

func MakeSysTimer() Converter {
	return &WindowsConverter{}
}

type WindowsConverter struct{}

func (d *WindowsConverter) To(data any) time.Time {
	ts := data.(*syscall.Win32FileAttributeData)
	eTime := earliestTime(fData.CreationTime.Nanoseconds(), fData.LastWriteTime.Nanoseconds())
	return time.Unix(0, eTime)
}

func earliestTime(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
