package util

import (
	"github.com/omakoto/go-common/src/utils"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	// E2E tests set _ZENLOG_TIME_INJECTION_FILE to the file that contains an injected time.
	timeOverrideFile = os.Getenv("_ZENLOG_TIME_INJECTION_FILE")
)

// GetInjectedNow returns an injected time if _ZENLOG_TIME_INJECTION_FILE is set, or otherwise just returns a passed Clock.
func GetInjectedNow(clock utils.Clock) time.Time {
	if timeOverrideFile == "" {
		return clock.Now()
	}
	bytes, err := ioutil.ReadFile(timeOverrideFile)
	Check(err, "ReadFile failed")
	i, err := strconv.ParseInt(strings.TrimRight(string(bytes), "\n"), 10, 64)
	Check(err, "ParseInt failed")

	err = ioutil.WriteFile(timeOverrideFile, []byte(strconv.FormatInt(i+1, 10)), 0600)
	Check(err, "WriteFile failed")

	return time.Unix(i, 0)
}
