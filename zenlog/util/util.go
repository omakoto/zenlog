package util

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	reFilenameSafe = NewLazyRegexp(`[^a-zA-Z0-9\-\_\.\+]+`)
	reSlashes      = NewLazyRegexp(`//+`)

	reRegexCleaner = NewLazyRegexp(`(?:\s+|\s*#[^\n]*\n\s*)`)
)

func init() {
	rand.Seed(time.Now().Unix())
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
