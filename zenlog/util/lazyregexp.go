package util

import (
	"regexp"
	"sync"
)

type LazyRegexp struct {
	once    *sync.Once
	pattern string
	regexp  *regexp.Regexp
}

func NewLazyRegexp(pattern string) LazyRegexp {
	return LazyRegexp{&sync.Once{}, pattern, nil}
}

func (l *LazyRegexp) Pattern() *regexp.Regexp {
	l.once.Do(func() {
		l.regexp = regexp.MustCompile(l.pattern)
	})
	return l.regexp
}
