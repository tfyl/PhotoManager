//go:build darwin

package sysTimer

import (
	"syscall"
	"time"
)

func MakeSysTimer() Converter {
	return &DarwinConverter{}
}

type DarwinConverter struct{}

func (d *DarwinConverter) To(data any) time.Time {
	ts := data.(*syscall.Stat_t)
	return time.Unix(ts.Birthtimespec.Unix())
}
