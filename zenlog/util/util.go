package util

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"math/rand"
	"os"
	"strings"
	"path/filepath"
)

var (
	Debug       = false
	outputIsRaw = false

	reFilenameSafe = NewLazyRegexp(`[^a-zA-Z0-9\-\_\.\+]+`)
	reSlashes      = NewLazyRegexp(`//+`)

	//debugOutSet = false
	debugOut = os.Stderr

	reRegexCleaner = NewLazyRegexp(`(?:\s+|\s*#[^\n]*\n\s*)`)
)

func init() {
	// If ZENLOG_DEBUG is set to '1', enable debug log.
	if os.Getenv("ZENLOG_DEBUG") == "1" {
		Debug = true
	}
}

func SetOutputIsRaw() {
	outputIsRaw = true
}

func formatMessage(format string, a ...interface{}) string {
	msg := fmt.Sprintf("zenlog: "+format, a...)

	if outputIsRaw {
		msg = strings.Replace(msg, "\n", "\r\n", -1)
	}

	return msg
}

func Debugf(format string, a ...interface{}) {
	if Debug {
		//if !debugOutSet {
		//	outer := os.Getenv(envs.ZENLOG_OUTER_TTY)
		//	if outer != "" {
		//		o, err := os.OpenFile(outer, os.O_WRONLY, 0)
		//		Warn(err, "Failed to openn")
		//		if o !=  nil {
		//			debugOut = o
		//		}
		//		debugOutSet = true
		//	}
		//}

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
}

func Dump(prefix string, obj interface{}) {
	if !Debug {
		return
	}
	Debugf("%s", prefix, spew.Sdump(obj))
}

func Fatalf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, formatMessage(format, a...))
	fmt.Fprint(os.Stderr, "\n")
	os.Exit(1)
}

func Check(err error, format string, a ...interface{}) {
	if Warn(err, format, a...) {
		os.Exit(1)
	}
}

func Say(format string, a ...interface{}) {
	message := formatMessage(format, a...)
	fmt.Fprint(os.Stderr, message)
	fmt.Fprint(os.Stderr, "\n")
}

func Warn(err error, format string, a ...interface{}) bool {
	if err != nil {
		message := formatMessage(format, a...)

		fmt.Fprint(os.Stderr, message)
		fmt.Fprint(os.Stderr, ": ")
		fmt.Fprint(os.Stderr, err.Error())
		fmt.Fprint(os.Stderr, "\n")
		return true
	}
	return false
}

func Exit(success bool) {
	status := 0
	if !success {
		status = 1
	}
	os.Exit(status)
}

func FirstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func FilenameSafe(s string) string {
	return reFilenameSafe.Pattern().ReplaceAllLiteralString(s, "_")
}

func FileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func CompressSlash(file string) string {
	return reSlashes.Pattern().ReplaceAllLiteralString(file, "/")
}

// Create a random string.
func Fingerprint() string {
	return fmt.Sprintf("%08x", rand.Int31())
}

// Remove whitespace and comments from a regex pattern.
func CleanUpRegexp(pattern string) string {
	return reRegexCleaner.Pattern().ReplaceAllLiteralString(pattern, "")
}

// Return number of LFs in bytes.
func NumLines(data []byte) int {
	ret := 0
	for _, b := range data {
		switch b {
		case '\n':
			ret++
		}
	}
	return ret
}

// Return the fullpath of the zenlog executable.
func FindSelf() string {
	me := os.Args[0]
	Debugf("$0=%s", me)

	path, err := filepath.EvalSymlinks(me)
	Check(err, "EvalSymlinks failed")

	path, err = filepath.Abs(path)
	Debugf("$0(resolved)=%s", path)

	return path
}

// Return the "signature" of the zenlog executable.
func Signature() string {
	self := FindSelf()
	stat, err := os.Stat(self)
	Check(err, "Stat failed")

	return fmt.Sprintf("%s:[%d]", self, stat.ModTime().Unix())
}