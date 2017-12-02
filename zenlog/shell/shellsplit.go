package shell

import (
	"bytes"
	"strings"
)

// Shell tokenizer.
// TODO Handle operators such as |.

type splitter struct {
	text []rune

	// Input.
	next int

	// Output.
	buffer   bytes.Buffer
	hasRunes bool

	result []string
}

func newSplitter(text string) splitter {
	return splitter{[]rune(text), 0, bytes.Buffer{}, false, make([]string, 0, 16)}
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

func (s *splitter) eatDoubleQuote(end rune) {
	for {
		r, ok := s.read()
		if r == end || !ok {
			return
		}
		if s.maybeEatDoller(r) {
			continue
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

func (s *splitter) maybeEatDoller(r rune) bool {
	if r == '$' {
		s.pushRune(r)
		next, ok := s.peek()
		if !ok {
			return true
		}
		if next == '(' {
			s.tokenize(')')
			return true
		}
		if next == '{' {
			s.tokenize('}')
			return true
		}
		return true
	}
	return false
}

func isWhitespace(r rune) bool {
	switch r {
	case ' ', '\t', '\r', '\n', '\v':
		return true
	}
	return false
}

func (s *splitter) tokenize(end int) {
	for {
		r, ok := s.read()
		if !ok {
			break
		}
		if r == '\\' {
			r, ok = s.read()
			if !ok {
				break
			}
			s.pushRune(r)
			continue
		}
		if r == '\'' {
			s.eatSingleQuote()
			continue
		}
		if r == '"' {
			s.eatDoubleQuote('"')
			continue
		}
		if r == '`' {
			s.eatDoubleQuote('`')
			continue
		}
		if s.maybeEatDoller(r) {
			continue
		}
		if  end < 0 && isWhitespace(r) {
			s.pushWord()
			continue
		}
		if end >= 0 && end == int(r) {
			s.pushRune(r)
			break
		}
		s.pushRune(r)
	}
	if end < 0 {
		s.pushWord()
	}
}

func ShellSplit(text string) []string {
	s := newSplitter(text)
	s.tokenize(-1)
	return s.result
}
