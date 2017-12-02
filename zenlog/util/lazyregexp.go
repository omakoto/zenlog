package util

import (
	"regexp"
)

// LazyRegexp is a sharable lazily compiled regexp.
type LazyRegexp struct {
	pattern string
	regexp  *regexp.Regexp
}

// NewLazyRegexp returns a sharable lazily compiled regexp.
func NewLazyRegexp(pattern string) LazyRegexp {
	return LazyRegexp{pattern, nil}
}

// Pattern returns a shared regexp.
func (l *LazyRegexp) Pattern() *regexp.Regexp {
	if l.regexp == nil { // No need to take a lock. Worst case, we'll just compile it multiple times.
		l.regexp = regexp.MustCompile(l.pattern)
	}
	return l.regexp
}
