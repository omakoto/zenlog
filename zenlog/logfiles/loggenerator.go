package logfiles

// Create log files and symlinks given a Command.

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/omakoto/go-common/src/fileutils"
	"github.com/omakoto/zenlog/zenlog/config"
	"github.com/omakoto/zenlog/zenlog/util"
)

const (
	SanDir = "SAN"
	RawDir = "RAW"
	EnvDir = "ENV"

	TodayLink     = "TODAY"
	ThisMonthLink = "THISMONTH"

	maxPrevLinks = 10
)

// LogFiles represents a set of log files (san, raw and env) for a single command.
type LogFiles struct {
	SanFile string
	RawFile string
	EnvFile string

	SanF *os.File
	RawF *os.File
	EnvF *os.File

	San *bufio.Writer
	Raw *bufio.Writer
	Env *bufio.Writer
}

func clamp(v string, maxLen int) string {
	if len(v) <= maxLen {
		return v
	}
	return v[0:maxLen]
}

// Create and open a log file with the parent directory, if needed.
func open(name string, truncate, append bool) (*os.File, *bufio.Writer) {
	os.MkdirAll(filepath.Dir(name), 0700)
	mode := os.O_WRONLY | os.O_CREATE
	if truncate {
		mode |= os.O_TRUNC
	}
	if append {
		mode |= os.O_APPEND
	}
	f, err := os.OpenFile(name, mode, 0600)
	util.Check(err, "Cannot create logfile %s", name)
	return f, bufio.NewWriter(f)
}

func makeSymlink(from, to string, warnIfExists bool) {
	if !fileutils.FileExists(to) {
		util.Warn(os.Symlink(from, to), "Symlink failed")
	} else if warnIfExists {
		util.Say("%s already exists", to)
	}
}

func ensureSymLink(from, to string) {
	fi, err := os.Lstat(to)
	if err != nil {
		if os.IsNotExist(err) {
			// okay
		} else {
			util.Warn(err, "lstat() failed")
		}
	} else {
		if fi.Mode()&os.ModeSymlink != 0 {
			e, _ := os.Readlink(to)
			if e == from {
				return // already exists.
			}
		}
		util.Warn(os.Remove(to), "Can't remove %s", to)
	}
	makeSymlink(from, to, true)
}

// Create "previous" links.
// fullDirName: Parent directory, such as "/zenlog/"
// logType: e.g. "SAN"
// logFullFileName: Symlink target log filename.
func createPrevLink(fullDirName, logType, logFullFileName string) {
	if !fileutils.FileExists(logFullFileName) {
		return // just in case.
	}
	oneLetter := logType[0:1]
	for i := maxPrevLinks; i >= 2; i-- {
		from := fullDirName + (strings.Repeat(oneLetter, i-1))
		if !fileutils.FileExists(from) {
			continue
		}
		to := fullDirName + (strings.Repeat(oneLetter, i))

		// No need to check the error.
		if fileutils.FileExists(to) {
			util.Warn(os.Remove(to), "Remove failed")
		}
		util.Warn(os.Rename(from, to), "Rename failed")
	}
	makeSymlink(logFullFileName, fullDirName+oneLetter, true)
}

// Create TODAY and THISMONTH links.
// fullDirName: Parent directory, such as "/zenlog/"
// logType: e.g. "SAN"
// logFullFileName: Symlink target log filename.
func createDayLinks(fullDirName, logType, logFullFileName string) {
	if !fileutils.FileExists(logFullFileName) {
		return // just in case.
	}
	todayFrom := filepath.Dir(logFullFileName)
	thisMonthFrom := filepath.Dir(todayFrom)

	todayTo := fullDirName + "/" + logType + "/" + TodayLink
	thisMonthTo := fullDirName + "/" + logType + "/" + ThisMonthLink

	ensureSymLink(todayFrom, todayTo)
	ensureSymLink(thisMonthFrom, thisMonthTo)
}

// Create auxiliary links.
// parentDirName: e.g. "cmds"
// childDirName: e.g. "cat"
// logType: e.g. "SAN"
// logFullFileName: Symlink target log filename.
func createLinks(config *config.Config, parentDirName, childDirName, logType, logFullFileName string, now time.Time) {
	if !fileutils.FileExists(logFullFileName) {
		return // just in case.
	}
	childDirName = clamp(util.FilenameSafe(childDirName), 64)
	if childDirName == "." || childDirName == ".." || childDirName == "" {
		return
	}

	fullChildDir := config.LogDir + parentDirName + "/" + childDirName + "/"

	fullDirName := fullChildDir + logType + "/" +
		fmt.Sprintf("%04d/%02d/%02d", now.Year(), now.Month(), now.Day()) + "/"

	util.Warn(os.MkdirAll(fullDirName, 0700), "MkdirAll failed")
	makeSymlink(logFullFileName, fullDirName+filepath.Base(logFullFileName), false)
	createPrevLink(fullChildDir, logType, logFullFileName)
}

func removePath(s string) string {
	return regexp.MustCompile(`^\S+/`).ReplaceAllString(s, "")
}

// CreateAndOpenLogFiles opens the log files for a command.
func CreateAndOpenLogFiles(config *config.Config, now time.Time, command *Command) LogFiles {
	l := LogFiles{}

	tag := ""
	if command.Comment != "" {
		tag = "_+" + clamp(util.FilenameSafe(command.Comment), 32)
	}

	now = now.Local()

	const M = "@@@"
	f := fmt.Sprintf("%s%s/%04d/%02d/%02d/%02d-%02d-%02d.%03d-%05d%s_+%s.log",
		config.LogDir,
		M,
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/1000000,
		config.ZenlogPid,
		tag,
		clamp(util.FilenameSafe(removePath(command.CommandLine)), 32))

	l.SanFile = strings.Replace(f, M, SanDir, 1)
	l.RawFile = strings.Replace(f, M, RawDir, 1)
	l.EnvFile = strings.Replace(f, M, EnvDir, 1)

	l.Open(true)

	spid := strconv.Itoa(config.ZenlogPid)
	items := []struct {
		fullLogFilename string
		logType         string
	}{
		{l.SanFile, SanDir},
		{l.RawFile, RawDir},
		{l.EnvFile, EnvDir},
	}
	for _, item := range items {
		createDayLinks(config.LogDir, item.logType, item.fullLogFilename)
		createPrevLink(config.LogDir, item.logType, item.fullLogFilename)
		createLinks(config, "pids", spid, item.logType, item.fullLogFilename, now)
		for _, exe := range command.ExeNames {
			createLinks(config, "cmds", exe, item.logType, item.fullLogFilename, now)
		}
		if command.Comment != "" {
			createLinks(config, "tags", command.Comment, item.logType, item.fullLogFilename, now)
		}
	}

	return l
}

// Open actually opens the set of the log files.
func (l *LogFiles) Open(truncate bool) {
	l.SanF, l.San = open(l.SanFile, truncate, true)
	l.RawF, l.Raw = open(l.RawFile, truncate, true)
	l.EnvF, l.Env = open(l.EnvFile, truncate, true)
}

func closeSingle(f **os.File, w **bufio.Writer) {
	if *w != nil {
		(*w).Flush()
		*w = nil
	}
	if *f != nil {
		(*f).Close()
		*f = nil
	}
}

// Close opens all the log files.
func (l *LogFiles) Close() {
	closeSingle(&l.SanF, &l.San)
	closeSingle(&l.RawF, &l.Raw)
	closeSingle(&l.EnvF, &l.Env)
}

// TODO Move the following functions to somewhere else...

func (l *LogFiles) writeTimeToEnv(key string, t time.Time) {
	l.Env.WriteString(key)
	l.Env.WriteString(": ")
	l.Env.WriteString(fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
	l.Env.WriteString("\n")
}

// WriteEnv writes the command start information to the ENV file, and a string given via 'env', to the ENV log file.
func (l *LogFiles) WriteEnv(command *Command, envs string, startTime time.Time) {
	l.Env.WriteString("Command: ")
	l.Env.WriteString(command.CommandLine)
	l.Env.WriteString("\n")

	l.writeTimeToEnv("Start time", startTime)

	l.Env.WriteString(envs)
	if envs[len(envs)-1] != '\n' {
		l.Env.WriteString("\n")
	}
}

// WriteFinishToEnv write the command finish information to the ENV file.
func (l *LogFiles) WriteFinishToEnv(exitCode int, startTime, finishTime time.Time) {
	l.Env.WriteString("Exit status: ")
	l.Env.WriteString(strconv.Itoa(exitCode))
	l.Env.WriteString("\n")

	l.writeTimeToEnv("Finish time", finishTime)

	l.Env.WriteString("Duration: ")
	l.Env.WriteString(finishTime.Sub(startTime).String())
	l.Env.WriteString("\n")
}
