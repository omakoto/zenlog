package shell

import (
	"strings"
	"bytes"
)

type splitter struct {
	text   []rune

	// Input.
	next   int

	// Output.
	buffer bytes.Buffer
	hasRunes bool

	result []string
}

func newSplitter(text string) splitter {
	return splitter{[]rune(text), 0,  bytes.Buffer{}, false, make([]string, 0, 16)}
}

func (s *splitter) peek() (rune, bool) {
	if s.next < len(s.text) {
		return s.text[s.next], true
	} else {
		return '\x00', false
	}
}

func (s *splitter) read() (rune, bool) {
	r, ok := s.peek()

	if ok {
		s.next += 1
	}
	return r, ok
}

func (s *splitter) hasNext() bool {
	return s.next < len(s.text)
}

func (s *splitter) expect(expectSet string) (rune, bool) {
	r, ok := s.peek()
	if ok {
		if strings.ContainsRune(expectSet, r) {
			s.next += 1
			return r, ok
		}
	}
	return '\x00', false
}

func (s *splitter) pushRune(r rune) {
	s.buffer.WriteRune(r)
	s.hasRunes = true
}

func (s *splitter) pushWord() {
	if s.hasRunes {
		s.hasRunes = false
		s.result = append(s.result, s.buffer.String())
		s.buffer = bytes.Buffer{}
	}
}

func (s *splitter) eatSingleQuote() {
	for {
		r, ok := s.read()
		if r == '\'' || !ok {
			return
		}
		s.pushRune(r)
	}
}

func (s *splitter) eatDoubleQuote() {
	for {
		r, ok := s.read()
		if r == '"' || !ok {
			return
		}
		if r == '\\' {
			r, ok = s.read()
			if !ok {
				return
			}
		}
		s.pushRune(r)
	}
}

func isWhitespace(r rune) bool {
	switch r {
	case ' ', '\t', '\r', '\n', '\v':
		return true
	}
	return false
}

func (s *splitter) tokenize() {
	for {
		r, ok := s.read()
		if !ok {
			break
		}
		if r == '\'' {
			s.eatSingleQuote()
			continue
		}
		if r == '"' {
			s.eatDoubleQuote()
			continue
		}
		if r == '\\' {
			r, ok = s.read()
			if !ok {
				break
			}
			s.pushRune(r)
			continue
		}
		if isWhitespace(r) {
			s.pushWord()
			continue
		}
		s.pushRune(r)
	}
	s.pushWord()
}

func ShellSplit(text string) []string {
	s := newSplitter(text)
	s.tokenize()
	return s.result
}
