package util

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func getIntEenv(name string, def int) int {
	ret, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		return def
	}
	return ret
}

var (
	zenlogPidInjected = getIntEenv("_ZENLOG_LOGGER_PID", os.Getpid())
	timeOverrideFile  = os.Getenv("_ZENLOG_TIME_INJECTION_FILE")
)

func GetLoggerPid() int {
	return zenlogPidInjected
}

func GetInjectedNow(clock Clock) time.Time {
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
