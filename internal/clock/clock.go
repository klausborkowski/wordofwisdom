package clock

import "time"

type clockKey string

const ClockCtxKey clockKey = "clock"

type SystemClock struct {
}

func (s SystemClock) Now() time.Time {
	return time.Now()
}
