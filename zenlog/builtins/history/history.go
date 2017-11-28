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

func history(pid, nth int, logType LogFileType, writer io.Writer) bool {
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
		for i := 1; i <= 10; i++ {
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

func HistoryCommand(args []string) {
	flags := flag.NewFlagSet("zenlog history", flag.ExitOnError)
	r := flags.Bool("r", false, "Print RAW filename")
	e := flags.Bool("e", false, "Print ENV filename")
	n := flags.Int("n", 0, "Print nth (>=1) last filename")
	p := flags.Int("p", 0, "Specify ZENLOG_PID")

	flags.Parse(args)

	util.Exit(history(*p, *n, flagsToLogType(*r, *e), os.Stdout))
}

func CurrentLogCommand(args []string) {
	flags := flag.NewFlagSet("zenlog current-log", flag.ExitOnError)
	r := flags.Bool("r", false, "Print RAW filename")
	e := flags.Bool("e", false, "Print ENV filename")
	p := flags.Int("p", 0, "Specify ZENLOG_PID")

	flags.Parse(args)

	util.Exit(history(*p, 1, flagsToLogType(*r, *e), os.Stdout))
}

func LastLogCommand(args []string) {
	flags := flag.NewFlagSet("zenlog last-log", flag.ExitOnError)
	r := flags.Bool("r", false, "Print RAW filename")
	e := flags.Bool("e", false, "Print ENV filename")
	p := flags.Int("p", 0, "Specify ZENLOG_PID")

	flags.Parse(args)

	util.Exit(history(*p, 2, flagsToLogType(*r, *e), os.Stdout))
}
