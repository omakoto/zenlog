package util

import "time"

// Clock is a mockable clock interface.
type Clock interface {
	Now() time.Time
}

type clock struct {
}

// Return the current time.
func (clock) Now() time.Time {
	return time.Now()
}

// Create a new (real) Clock.
func NewClock() Clock {
	return clock{}
}

// InjectedClock is a mock clock.
type InjectedClock struct {
	time time.Time
}

func (c InjectedClock) Now() time.Time {
	return c.time
}

// NewInjectedClock creates a new mock clock.
func NewInjectedClock(time time.Time) Clock {
	return InjectedClock{time}
}
