package util

import (
	"fmt"
	"github.com/omakoto/go-common/src/utils"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	reFilenameSafe = utils.NewLazyRegexp(`[^a-zA-Z0-9\-\_\.\+]+`)
	reSlashes      = utils.NewLazyRegexp(`//+`)
)

func init() {
	rand.Seed(time.Now().Unix())
}

func FilenameSafe(s string) string {
	return strings.TrimRight(reFilenameSafe.Pattern().ReplaceAllLiteralString(s, "_"), "_")
}

func CompressSlash(file string) string {
	return reSlashes.Pattern().ReplaceAllLiteralString(file, "/")
}

// Create a random string.
func Fingerprint() string {
	return fmt.Sprintf("%08x", rand.Int31())
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

// GetIntEnv extracts an integer from an environmental variable.
func GetIntEnv(name string, def int) int {
	ret, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		return def
	}
	return ret
}
