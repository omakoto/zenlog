package logfiles

// Create log files and symlinks given a Command.

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"github.com/omakoto/zenlog-go/zenlog/config"
)

const (
	SAN = "SAN"
	RAW = "RAW"
	ENV = "ENV"

	MAX_PREV_LINKS = 10
)

type LogFiles struct {
	SanFile string
	RawFile string
	EnvFile string

	San *os.File
	Raw *os.File
	Env *os.File
}

func clamp(v string, maxLen int) string {
	if len(v) <= maxLen {
		return v
	}
	return v[0:maxLen]
}

// Create and open a log file with the parent directory, if needed.
func open(name string) *os.File {
	os.MkdirAll(filepath.Dir(name), 0700)
	f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	util.Check(err, "Cannot create logfile %s", name)
	return f
}

// Create "previous" links.
// fullDirName: Parent directory, such as "/zenlog/"
// logType: e.g. "SAN"
// logFullFileName: Symlink target log filename.
func createPrevLink(fullDirName, logType, logFullFileName string) {
	if !util.FileExists(logFullFileName) {
		return // just in case.
	}
	oneLetter := logType[0:1]
	for i := MAX_PREV_LINKS; i >= 2; i-- {
		from := fullDirName + (strings.Repeat(oneLetter, i-1))
		if !util.FileExists(from) {
			continue
		}
		to := fullDirName + (strings.Repeat(oneLetter, i))

		// No nened to check the error.
		if util.FileExists(to) {
			util.Warn(os.Remove(to), "Remove failed")
		}
		util.Warn(os.Rename(from, to), "Rename failed")
	}
	util.Warn(os.Symlink(logFullFileName, fullDirName+oneLetter), "Symlink failed")
}

// Create auxiliary links.
// parentDirName: e.g. "cmds"
// childDirName: e.g. "cat"
// logType: e.g. "SAN"
// logFullFileName: Symlink target log filename.
func createLinks(config *config.Config, parentDirName, childDirName, logType, logFullFileName string, now time.Time) {
	if !util.FileExists(logFullFileName) {
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
	util.Warn(os.Symlink(logFullFileName, fullDirName+filepath.Base(logFullFileName)), "Symlink failed")
	createPrevLink(fullChildDir, logType, logFullFileName)
}

func OpenLogFiles(config *config.Config, now time.Time, command *Command) LogFiles {
	ret := LogFiles{}

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
		clamp(util.FilenameSafe(command.CommandLine), 32))

	ret.SanFile = strings.Replace(f, M, SAN, 1)
	ret.RawFile = strings.Replace(f, M, RAW, 1)
	ret.EnvFile = strings.Replace(f, M, ENV, 1)

	ret.Open()

	spid := strconv.Itoa(config.ZenlogPid)
	items := []struct {
		fullLogFilename string
		logType    string
	}{
		{ret.SanFile, SAN},
		{ret.RawFile, RAW},
		{ret.EnvFile, ENV},
	}
	for _, item := range items {
		createPrevLink(config.LogDir, item.logType, item.fullLogFilename)
		createLinks(config, "pids", spid, item.logType, item.fullLogFilename, now)
		for _, exe := range command.ExeNames {
			createLinks(config, "cmds", exe, item.logType, item.fullLogFilename, now)
		}
		if command.Comment != "" {
			createLinks(config, "tags", command.Comment, item.logType, item.fullLogFilename, now)
		}
	}

	return ret
}

func (l *LogFiles) Open() {
	l.San = open(l.SanFile)
	l.Raw = open(l.RawFile)
	l.Env = open(l.EnvFile)
}

func closeSingle(f **os.File) {
	if *f != nil {
		(*f).Close()
		*f = nil
	}
}

func (l *LogFiles) Close() {
	closeSingle(&l.San)
	closeSingle(&l.Raw)
	closeSingle(&l.Env)
}

// TODO Move the following functions to somewhere else...

func (l *LogFiles) writeTimeToEnv(key string, t time.Time) {
	l.Env.WriteString(key)
	l.Env.WriteString(": ")
	l.Env.WriteString(fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
	l.Env.WriteString("\n")
}

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

func (l *LogFiles) WriteFinishToEnv(exitCode int, startTime, finishTime time.Time) {
	l.Env.WriteString("Exit status: ")
	l.Env.WriteString(strconv.Itoa(exitCode))
	l.Env.WriteString("\n")

	l.writeTimeToEnv("Finish time", finishTime)

	l.Env.WriteString("Duration: ")
	l.Env.WriteString(finishTime.Sub(startTime).String())
	l.Env.WriteString("\n")
}
