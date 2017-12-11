package shell

import (
	"fmt"
	"os"
	"strconv"
)

const (
	readlineLine  = "READLINE_LINE"
	readlinePoint = "READLINE_POINT"
)

type BashProxy struct {
}

func GetBashProxy() Proxy {
	return &BashProxy{}
}

// GetCommandLine return the current command line and the cursor position from the READLINE_* environmental variables.
func (b *BashProxy) GetCommandLine() (string, int) {
	s := os.Getenv(readlineLine)
	l, err := strconv.Atoi(os.Getenv(readlinePoint))
	if err != nil || l < 0 {
		l = len(s)
	}
	return s, l
}

// UpdateCommandLine prints a string that can be evaled by bash to update the READLINE_* environmental variables
// to update the current command line.
func (b *BashProxy) PrintUpdateCommandLineEvalStr(commandLine string, cursorPos int) {
	fmt.Print(readlineLine)
	fmt.Print("=")
	fmt.Println(Escape(commandLine))

	fmt.Print(readlinePoint)
	fmt.Print("=")
	fmt.Println(strconv.Itoa(cursorPos))
}
