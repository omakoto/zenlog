package history

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/omakoto/zenlog-go/zenlog/config"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"io"
	"os"
	"strings"
)

type LogFileType int

const (
	LogTypeSan LogFileType = iota
	LogTypeRaw
	LogTypeEnv
)

func logChar(logType LogFileType) string {
	switch logType {
	case LogTypeSan:
		return "S"
	case LogTypeRaw:
		return "R"
	case LogTypeEnv:
		return "E"
	}
	util.Fatalf("Unknown type %d", logType)
	return ""
}

func writeIfLink(w *bufio.Writer, filename string) bool {
	util.Debugf("Checking %s", filename)
	to, err := os.Readlink(filename)
	if err != nil {
		return false
	}
	util.Debugf(" -> %s", to)
	w.WriteString(to)
	w.WriteString("\n")
	return true
}

func NthLastLog(config *config.Config, pid, nth int, logType LogFileType) string {
	if pid <= 0 {
		pid = config.ZenlogPid
	}
	ch := logChar(logType)
	dir := fmt.Sprintf("%spids/%d/", config.LogDir, pid)
	link := dir + strings.Repeat(ch, nth)
	if !util.FileExists(link) {
		return ""
	}
	file, err := os.Readlink(link)
	if err != nil {
		util.Warn(err, "Readlink failed")
		return ""
	}
	return file
}

func history(pid, nth int, logType LogFileType, writer io.Writer) bool {
	if nth < 0 {
		util.Fatalf("Invalid argument for nth: %d", nth)
	}
	config := config.InitConfigForCommands()

	if pid <= 0 {
		pid = config.ZenlogPid
	}
	w := bufio.NewWriter(writer)
	defer w.Flush()

	success := false
	if nth > 0 {
		file := NthLastLog(config, pid, nth, logType)
		if file != "" {
			w.WriteString(file)
			w.WriteString("\n")
			success = true
		}
	} else {
		dir := fmt.Sprintf("%spids/%d/", config.LogDir, pid)
		ch := logChar(logType)
		for i := 10; i >= 1; i-- {
			success = writeIfLink(w, dir+strings.Repeat(ch, i)) || success
		}
	}
	return success
}

func flagsToLogType(flagR, flagE bool) LogFileType {
	if flagR {
		return LogTypeRaw
	}
	if flagE {
		return LogTypeEnv
	}
	return LogTypeSan
}

// AllHistoryCommand is the implementation of "zenlog history".
func AllHistoryCommand(args []string) {
	flags := flag.NewFlagSet("zenlog history", flag.ExitOnError)
	r := flags.Bool("r", false, "Print RAW filename")
	e := flags.Bool("e", false, "Print ENV filename")
	n := flags.Int("n", 0, "Print nth (>=1) last filename")
	p := flags.Int("p", 0, "Specify ZENLOG_PID")

	flags.Parse(args)

	util.Exit(history(*p, *n, flagsToLogType(*r, *e), os.Stdout))
}

// CurrentLogCommand is the implementation of "zenlog current-log".
func CurrentLogCommand(args []string) {
	flags := flag.NewFlagSet("zenlog current-log", flag.ExitOnError)
	r := flags.Bool("r", false, "Print RAW filename")
	e := flags.Bool("e", false, "Print ENV filename")
	p := flags.Int("p", 0, "Specify ZENLOG_PID")

	flags.Parse(args)

	util.Exit(history(*p, 1, flagsToLogType(*r, *e), os.Stdout))
}

// LastLogCommand is the implementation of "zenlog last-log".
func LastLogCommand(args []string) {
	flags := flag.NewFlagSet("zenlog last-log", flag.ExitOnError)
	r := flags.Bool("r", false, "Print RAW filename")
	e := flags.Bool("e", false, "Print ENV filename")
	p := flags.Int("p", 0, "Specify ZENLOG_PID")

	flags.Parse(args)

	util.Exit(history(*p, 2, flagsToLogType(*r, *e), os.Stdout))
}
