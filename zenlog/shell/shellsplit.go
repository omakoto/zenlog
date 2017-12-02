package shell

import (
	"bytes"
)

// Shell tokenizer.
// TODO Handle operators such as |.

type splitter struct {
	text []rune

	// Input.
	next int

	wasSpecial bool

	// Output.
	buffer   bytes.Buffer
	hasRunes bool

	result []string
}

func newSplitter(text string) splitter {
	return splitter{[]rune(text), 0, false, bytes.Buffer{}, false, make([]string, 0, 16)}
}

func (s *splitter) peek() (rune, bool) {
	if s.next < len(s.text) {
		return s.text[s.next], true
	}
	return '\x00', false
}

func (s *splitter) read() (rune, bool) {
	r, ok := s.peek()

	if ok {
		s.next++
	}
	return r, ok
}

func (s *splitter) pushRuneNoSpecial(r rune) {
	s.buffer.WriteRune(r)
	s.hasRunes = true
}

func (s *splitter) pushRune(r rune) {
	s.pushRuneNoSpecial(r)
	s.wasSpecial = isSpecialChar(r)
}

func (s *splitter) pushWord() {
	if s.hasRunes {
		s.hasRunes = false
		s.result = append(s.result, s.buffer.String())
		s.buffer = bytes.Buffer{}
		s.wasSpecial = false
	}
}

func (s *splitter) eatSingleQuote() {
	for {
		r, ok := s.read()
		if !ok {
			return
		}
		if r == '\'' {
			s.pushRune(r)
			return
		}
		s.pushRune(r)
	}
}

func (s *splitter) eatDoubleQuote(end rune) {
	for {
		r, ok := s.read()
		if !ok {
			return
		}
		if r == end {
			s.pushRune(r)
			return
		}
		if s.maybeEatDollar(r) {
			continue
		}
		if r == '\\' {
			s.pushRuneNoSpecial(r)
			r, ok = s.read()
			if !ok {
				return
			}
		}
		s.pushRune(r)
	}
}

func (s *splitter) maybeEatDollar(r rune) bool {
	if r == '$' {
		s.pushRune(r)

		next, ok := s.peek()
		if !ok || isWhitespace(next) {
			return true
		}
		if next == '(' {
			s.read()
			s.pushRuneNoSpecial(next)
			s.tokenize(')')
			return true
		}
		if next == '{' {
			s.read()
			s.pushRuneNoSpecial(next)
			s.tokenize('}')
			return true
		}
		if next == '\'' {
			s.read()
			s.pushRuneNoSpecial(next)
			s.eatSingleQuote()
			return true
		}
		if next == '"' {
			s.read()
			s.pushRuneNoSpecial(next)
			s.eatDoubleQuote('"')
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

func isSpecialChar(r rune) bool {
	switch r {
	case ';', '!', '<', '>', '(', ')', '|', '&':
		return true
	}
	return false
}

func isCommandSeparatorChar(r rune) bool {
	switch r {
	case ';', '(', ')', '|', '&':
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
		if end >= 0 && end == int(r) {
			s.pushRuneNoSpecial(r)
			break
		}
		if end < 0 && isWhitespace(r) {
			s.pushWord()
			continue
		}
		if s.wasSpecial != isSpecialChar(r) {
			s.pushWord()
		}
		if r == '\\' {
			s.pushRune(r)
			r, ok = s.read()
			if !ok {
				break
			}
			s.pushRune(r)
			continue
		}
		if r == '\'' {
			s.pushRune(r)
			s.eatSingleQuote()
			continue
		}
		if r == '"' {
			s.pushRune(r)
			s.eatDoubleQuote('"')
			continue
		}
		if r == '`' {
			s.pushRune(r)
			s.eatDoubleQuote('`')
			continue
		}
		if s.maybeEatDollar(r) {
			continue
		}
		if !s.hasRunes && r == '#' {
			for {
				s.pushRuneNoSpecial(r)
				r, ok = s.read()
				if !ok {
					s.pushWord()
					return
				}
			}
		}
		s.pushRune(r)
	}
	if end < 0 {
		s.pushWord()
	}
}

// Split splits a whole command line into tokens.
// Example: "cat fi\ le.txt|grep -V ^# >'out$$.txt' # Find non-comment lines."
// -> output: cat, fi\ le.txt, |, grep, -V, ^#, >, out.txt ,# Find non-comment lines.
func Split(text string) []string {
	s := newSplitter(text)
	s.tokenize(-1)
	return s.result
}

func IsCommandSeparator(text string) bool {
	for _, r := range text {
		if !isCommandSeparatorChar(r) {
			return false
		}
	}
	return true
}
