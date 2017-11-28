package util

import (
	"regexp"
)

type Sanitizer struct {
	reSanitizer *regexp.Regexp
	reCrLf      *regexp.Regexp
	reCr        *regexp.Regexp
	reBs *regexp.Regexp

	empty     []byte
	nl        []byte
	bs []byte
}

func NewSanitizer() *Sanitizer {
	s := Sanitizer{}

	/*
	   str.gsub! %r!(
	         \a                         # Bell
	         | \e \x5B .*? [\x40-\x7E]  # CSI
	         | \e \x5D .*? \x07         # Set terminal title
	         | \e \( .                  # 3 byte sequence
	         | \e [\x40-\x5A\x5C\x5F]   # 2 byte sequence
	         )!x, ""

	   # Also clean up CR/LFs.
	   str.gsub!(%r! \s* \x0d* \x0a !x, "\n") # Remove end-of-line CRs.
	   str.gsub!(%r! \s* \x0d !x, "\n")       # Replace orphan CRs with LFs.

	   # Also replace ^H's.
	   str.gsub!(%r! \x08 !x, '^H');
	*/
	s.reSanitizer = regexp.MustCompile(CleanUpRegexp(`
	         \a                           # Bell
	         | \x1b \x5B .*? [\x40-\x7E]  # CSI
	         | \x1b \x5D .*? \x07         # Set terminal title
	         | \x1b \( .                  # 3 byte sequence
	         | \x1b [\x40-\x5A\x5C\x5F]   # 2 byte sequence
		`))

	s.reCrLf = regexp.MustCompile(`\s*\r*\n`)
	s.reCr = regexp.MustCompile(`\s*\r`)
	s.reBs = regexp.MustCompile(`\x08`)

	s.empty = make([]byte, 0)
	s.nl = []byte("\n")
	s.bs = []byte("^H")

	return &s
}

func (s *Sanitizer) Sanitize(data []byte) []byte {

	data = s.reSanitizer.ReplaceAllLiteral(data, s.empty)
	data = s.reCrLf.ReplaceAllLiteral(data, s.nl)
	data = s.reCr.ReplaceAllLiteral(data, s.nl)
	data = s.reBs.ReplaceAllLiteral(data, s.bs)

	return data
}
