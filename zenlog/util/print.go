package util

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

var (
	// Debug is whether the debug output is enabled or not.
	Debug       = false
	outputIsRaw = false

	//debugOutSet = false
	debugOut = os.Stderr
)

func init() {
	// If ZENLOG_DEBUG is set to '1', enable debug log.
	if os.Getenv("ZENLOG_DEBUG") == "1" {
		Debug = true
	}
}

func maybePrintStackTrack() {
	if Debug || os.Getenv("ZENLOG_PRINT_STACK") == "1" {
		debug.PrintStack()
	}
}

// SetOutputIsRaw sets whether stdout is in raw mode or not.
func SetOutputIsRaw(raw bool) {
	outputIsRaw = raw
}

func getNewLine() string {
	if outputIsRaw {
		return "\r\n"
	}
	return "\n"
}

func replaceLf(s string) string {
	if outputIsRaw {
		s = strings.Replace(s, "\n", "\r\n", -1)
	}
	return s
}

func formatMessage(format string, a ...interface{}) string {
	return replaceLf(fmt.Sprintf("zenlog: "+format, a...))
}

func Debugf(format string, a ...interface{}) {
	if Debug {
		DebugfForce(format, a...)
	}
}

func DebugfForce(format string, a ...interface{}) {
	color := ""
	end := ""
	if outputIsRaw {
		// Logger side
		color = "\x1b[0m\x1b[1;32m[L]" // Append [L]
		end = "\x1b[0m\r\n"            // Note the \r.
	} else {
		color = "\x1b[0m\x1b[1;33m"
		end = "\x1b[0m\n"
	}
	fmt.Fprint(debugOut, color)
	fmt.Fprint(debugOut, formatMessage(format, a...))
	fmt.Fprint(debugOut, end)
}

func Dump(prefix string, obj interface{}) {
	if !Debug {
		return
	}
	Debugf("%s%s", prefix, spew.Sdump(obj))
}

func Fatalf(format string, a ...interface{}) {
	//fmt.Fprint(os.Stderr, "\x1b[0m\x1b[1;31m")
	fmt.Fprint(os.Stderr, formatMessage(format, a...))
	//fmt.Fprint(os.Stderr, "\x1b[0m")
	fmt.Fprint(os.Stderr, getNewLine())
	maybePrintStackTrack()
	ExitFailure()
}

func Check(err error, format string, a ...interface{}) {
	if Warn(err, format, a...) {
		ExitFailure()
	}
}

func Say(format string, a ...interface{}) {
	message := formatMessage(format, a...)
	fmt.Fprint(os.Stderr, "\x1b[0m\x1b[1;33m")
	fmt.Fprint(os.Stderr, message)
	fmt.Fprint(os.Stderr, "\x1b[0m")
	fmt.Fprint(os.Stderr, getNewLine())
}

func Warn(err error, format string, a ...interface{}) bool {
	if err != nil {
		message := formatMessage(format, a...)
		fmt.Fprint(os.Stderr, "\x1b[0m\x1b[1;33m")
		fmt.Fprint(os.Stderr, message)
		fmt.Fprint(os.Stderr, ": ")
		fmt.Fprint(os.Stderr, err.Error())
		fmt.Fprint(os.Stderr, "\x1b[0m")
		fmt.Fprint(os.Stderr, getNewLine())
		maybePrintStackTrack()
		return true
	}
	return false
}
