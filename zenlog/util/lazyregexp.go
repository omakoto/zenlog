package util

import (
	"regexp"
)

type LazyRegexp struct {
	pattern string
	regexp  *regexp.Regexp
}

func NewLazyRegexp(pattern string) LazyRegexp {
	return LazyRegexp{pattern, nil}
}

func (l *LazyRegexp) Pattern() *regexp.Regexp {
	if l.regexp == nil { // No need to take a lock. Worst case, we'll just compile it multiple times.
		l.regexp = regexp.MustCompile(l.pattern)
	}
	return l.regexp
}
