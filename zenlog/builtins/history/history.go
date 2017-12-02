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

type logFileType int

const (
	logTypeSan logFileType = iota
	logTypeRaw
	logTypeEnv
)

func logChar(logType logFileType) string {
	switch logType {
	case logTypeSan:
		return "S"
	case logTypeRaw:
		return "R"
	case logTypeEnv:
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

func history(pid, nth int, logType logFileType, writer io.Writer) bool {
	if nth < 0 {
		util.Fatalf("Invalid argument for nth: %d", nth)
	}
	config := config.InitConfigForCommands()

	if pid <= 0 {
		pid = config.ZenlogPid
	}
	dir := fmt.Sprintf("%spids/%d/", config.LogDir, pid)
	w := bufio.NewWriter(writer)
	defer w.Flush()

	ch := logChar(logType)

	success := false
	if nth > 0 {
		success = writeIfLink(w, dir+strings.Repeat(ch, nth))
	} else {
		for i := 10; i >= 1; i-- {
			success = writeIfLink(w, dir+strings.Repeat(ch, i)) || success
		}
	}
	return success
}

func flagsToLogType(flagR, flagE bool) logFileType {
	if flagR {
		return logTypeRaw
	}
	if flagE {
		return logTypeEnv
	}
	return logTypeSan
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
