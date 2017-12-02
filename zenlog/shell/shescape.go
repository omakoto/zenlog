package shell

import (
	"bytes"
	"strings"

	"github.com/omakoto/zenlog-go/zenlog/util"
)

var (
	reNeedsEscaping = util.NewLazyRegexp(`[^a-zA-Z0-9\-\.\_\/]`)
)

// Escape a string for shell.
func Shescape(s string) string {
	if !reNeedsEscaping.Pattern().MatchString(s) {
		return s
	}
	var buffer bytes.Buffer
	buffer.WriteString("'")
	buffer.WriteString(strings.Replace(s, `'`, `'\''`, -1))
	buffer.WriteString("'")
	return buffer.String()
}
