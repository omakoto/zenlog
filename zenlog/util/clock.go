package util

import "time"

type Clock interface {
	Now() time.Time
}

type clock struct {
}

func (clock) Now() time.Time {
	return time.Now()
}

func NewClock() Clock {
	return clock{}
}

type InjectedClock struct {
	time time.Time
}

func (c InjectedClock) Now() time.Time {
	return c.time
}

func NewInjectedClock(time time.Time) Clock {
	return InjectedClock{time}
}