package zenlog

//import (
//	"strings"
//	"bytes"
//)
//
//type splitter struct {
//	text   []rune
//
//	// Input.
//	next   int
//
//	// Output.
//	buffer bytes.Buffer
//	hasRunes bool
//
//	result []string
//}
//
//func newSplitter(text string) splitter {
//	return splitter{[]rune(text), 0,  bytes.Buffer{}, false, make([]string, 0, 16)}
//}
//
//func (s *splitter) peek() (rune, bool) {
//	if s.next < len(s.text) {
//		return s.text[s.next], true
//	} else {
//		return '\x00', false
//	}
//}
//
//func (s *splitter) read() (rune, bool) {
//	r, ok := s.peek()
//
//	if ok {
//		s.next += 1
//	}
//	return r, ok
//}
//
//func (s *splitter) expect(expectSet string) (rune, bool) {
//	r, ok := s.peek()
//	if ok {
//		if strings.ContainsRune(expectSet, r) {
//			s.next += 1
//			return r, ok
//		}
//	}
//	return '\x00', false
//}
//
//func (s *splitter) push(r rune) {
//	s.buffer.WriteRune(r)
//	s.hasRunes = true
//}
//
//func (s *splitter) pushWord() {
//	if s.hasRunes {
//		s.hasRunes = false
//		s.result = append(s.result, s.buffer.String())
//		s.buffer = bytes.Buffer{}
//	}
//}
//
//func ShellSplit(text string) []string {
//
//	s := newSplitter(text)
//	for {
//		r, ok := s.read()
//		if !ok {
//			break
//		}
//		switch r {
//		case '\'':
//			for {
//				r, ok = s.read()
//				if r == '\'' || !ok {
//					break
//				}
//				s.push(r)
//			}
//		case '"':
//			for {
//				r, ok = s.read()
//				if r == '"' || !ok {
//					break
//				}
//				if r == '\\' {
//				}
//				s.push(r)
//			}
//		}
//
//	}
//	s.pushWord()
//}
