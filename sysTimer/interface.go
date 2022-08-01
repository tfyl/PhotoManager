package sysTimer

import "time"

type Converter interface {
	To(data any) time.Time
}
