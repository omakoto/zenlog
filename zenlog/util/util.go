package util

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

var (
	// Debug is whether the debug output is enabled or not.
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

	rand.Seed(time.Now().Unix())
}

// SetOutputIsRaw sets whether stdout is in raw mode or not.
func SetOutputIsRaw(raw bool) {
	outputIsRaw = raw
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
	fmt.Fprintf(os.Stderr, formatMessage(format, a...))
	fmt.Fprint(os.Stderr, "\n")
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
	fmt.Fprint(os.Stderr, replaceLf("\x1b[0m\n"))
}

func Warn(err error, format string, a ...interface{}) bool {
	if err != nil {
		message := formatMessage(format, a...)
		fmt.Fprint(os.Stderr, "\x1b[0m\x1b[1;33m")
		fmt.Fprint(os.Stderr, message)
		fmt.Fprint(os.Stderr, ": ")
		fmt.Fprint(os.Stderr, err.Error())
		fmt.Fprint(os.Stderr, replaceLf("\x1b[0m\n"))
		return true
	}
	return false
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

func DirExists(file string) bool {
	stat, err := os.Stat(file)
	return err == nil && ((stat.Mode() & os.ModeDir) != 0)
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

// FindZenlogBin returns the fullpath of the zenlog executable.
func FindZenlogBin() string {
	me, err := os.Executable()
	Check(err, "Executable failed")
	Debugf("$0=%s", me)

	path, err := filepath.EvalSymlinks(me)
	Check(err, "EvalSymlinks failed")

	path, err = filepath.Abs(path)
	Check(err, "Abs failed")
	Debugf("$0(resolved)=%s", path)

	return path
}

// FindZenlogBinDir returns the fullpath of the directory where the zenlog executable is.
func FindZenlogBinDir() string {
	return filepath.Dir(FindZenlogBin())
}

// ZenlogBinCtime returns the ctime of the zenlog executable.
func ZenlogBinCtime() time.Time {
	stat, err := os.Stat(FindZenlogBin())
	Check(err, "Stat failed")
	return stat.ModTime()
}

// ZenlogSrcTopDir returns the fullpath of the source top directory.
func ZenlogSrcTopDir() string {
	zenlogBinDir := FindZenlogBinDir()

	for _, d := range StringSlice("/../", "/../src/github.com/omakoto/zenlog-go/") {
		candidate := zenlogBinDir + d
		candidate, err := filepath.Abs(candidate)
		Check(err, "Abs failed")

		if DirExists(candidate + "/subcommands") {
			return candidate
		}
	}
	log.Fatalf("Zenlog source directory not found.")
	return ""
}

// StringSlice is a convenient way to build a string slice.
func StringSlice(arr ...string) []string {
	return arr
}

// GetIntEnv extracts an integer from an environmental variable.
func GetIntEnv(name string, def int) int {
	ret, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		return def
	}
	return ret
}
