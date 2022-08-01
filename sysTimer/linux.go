//go:build linux

package sysTimer

import (
	"syscall"
	"time"
)

func MakeSysTimer() Converter {
	return &LinuxConverter{}
}

type LinuxConverter struct{}

func (d *LinuxConverter) To(data any) time.Time {
	ts := data.(*syscall.Stat_t)
	mtime := time.Unix(ts.Atim.Unix())
	ctime := time.Unix(ts.Ctim.Unix())
	atime := time.Unix(ts.Atim.Unix())

	return lowestTime(mtime, ctime, atime)
}

func lowestTime(times ...time.Time) time.Time {
	lowestTs := times[0]

	for _, t := range times {
		if t.Before(lowestTs) {
			lowestTs = t
		}
	}

	return lowestTs
}
