package builtins

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/omakoto/zenlog-go/zenlog/builtins/history"
	"github.com/omakoto/zenlog-go/zenlog/envs"
	"github.com/omakoto/zenlog-go/zenlog/logger"
	"github.com/omakoto/zenlog-go/zenlog/util"
)

// InZenlog true if the current process is running in a zenlog session.
func InZenlog() bool {
	sig := util.Tty() + ":" + logger.Signature()
	util.Debugf("Signature=%s", sig)
	return sig == os.Getenv(envs.ZenlogSignature)
}

// FailIfInZenlog quites the current process with an error code with an error message if it's running in a zenlog session.
func FailIfInZenlog() {
	if InZenlog() {
		util.Fatalf("Already in zenlog.")
	}
}

// FailUnlessInZenlog quites the current process with an error code with an error message unless it's running in a zenlog session.
func FailUnlessInZenlog() {
	if !InZenlog() {
		util.Fatalf("Not in zenlog.")
	}
}

func copyStdinToFile(file string) {
	out, err := os.OpenFile(file, os.O_WRONLY, 0)
	util.Check(err, "Unable to open "+file)
	io.Copy(out, os.Stdin)
}

// WriteToLogger read from STDIN and writes to the current logger. Implies FailUnlessInZenlog().
func WriteToLogger() {
	FailUnlessInZenlog()
	copyStdinToFile(os.Getenv(envs.ZenlogLoggerIn))
}

// WriteToOuter read from STDIN and writes to the console, without logging. Implies FailUnlessInZenlog().
func WriteToOuter() {
	FailUnlessInZenlog()
	file := os.Getenv(envs.ZenlogOuterTty)
	out, err := os.OpenFile(file, os.O_WRONLY, 0)
	util.Check(err, "Unable to open "+file)

	in := bufio.NewReader(os.Stdin)

	crlf := make([]byte, 2)
	crlf[0] = '\r'
	crlf[1] = '\n'

	for {
		line, err := in.ReadBytes('\n')
		if line != nil {
			line = bytes.TrimRight(line, "\r\n")
			out.Write(line)
			out.Write(crlf)
		}
		if err != nil {
			break
		}
	}
}

// OuterTty prints the outer TTY device filename. Implies FailUnlessInZenlog().
func OuterTty() {
	FailUnlessInZenlog()
	fmt.Println(os.Getenv(envs.ZenlogOuterTty))
}

// LoggerPipe prints named pile filename to the logger. Implies FailUnlessInZenlog().
func LoggerPipe() {
	FailUnlessInZenlog()
	fmt.Println(os.Getenv(envs.ZenlogLoggerIn))
}

func checkUpdate() {
	if strconv.FormatInt(util.ZenlogBinCtime().Unix(), 10) == os.Getenv(envs.ZenlogBinCtime) {
		util.ExitSuccess()
	}
	util.Say("Zenlog updated. Run \"zenlog_restart\" (or \"exit 13\") to restart a zenlog session.")
	util.ExitFailure()
}

// MaybeRunBuiltin runs a builtin command if a given command is a builtin subcommand.
func MaybeRunBuiltin(command string, args []string) {
	switch strings.Replace(command, "_", "-", -1) {
	case "in-zenlog":
		util.Exit(InZenlog())

	case "fail-if-in-zenlog":
		FailIfInZenlog()

	case "fail-unless-in-zenlog":
		FailUnlessInZenlog()

	case "write-to-logger":
		FailUnlessInZenlog()
		WriteToLogger()

	case "write-to-outer":
		FailUnlessInZenlog()
		WriteToOuter()

	case "outer-tty":
		FailUnlessInZenlog()
		OuterTty()

	case "logger-pipe":
		FailUnlessInZenlog()
		LoggerPipe()

		// History related commands.
	case "history":
		FailUnlessInZenlog()
		history.AllHistoryCommand(args)

	case "current-log":
		FailUnlessInZenlog()
		history.CurrentLogCommand(args)

	case "last-log":
		FailUnlessInZenlog()
		history.LastLogCommand(args)

	case "check-update":
		FailUnlessInZenlog()
		checkUpdate()

	case "start-command":
		FailUnlessInZenlog()
		startCommand(args)

	case "stop-log", "end-command":
		FailUnlessInZenlog()
		stopCommand(args)

	default:
		return
	}
	util.ExitSuccess()
}
