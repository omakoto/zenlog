package builtins

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/omakoto/zenlog-go/zenlog"
	"github.com/omakoto/zenlog-go/zenlog/builtins/history"
	"github.com/omakoto/zenlog-go/zenlog/envs"
	"github.com/omakoto/zenlog-go/zenlog/logger"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"io"
	"os"
	"strconv"
	"strings"
)

func InZenlog() bool {
	sig := util.Tty() + ":" + zenlog.Signature()
	util.Debugf("Signature=%s", sig)
	return sig == os.Getenv(envs.ZENLOG_SIGNATURE)
}

func FailIfInZenlog() {
	if InZenlog() {
		util.Fatalf("Already in zenlog.")
	}
}

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

func WriteToLogger() {
	FailUnlessInZenlog()
	copyStdinToFile(os.Getenv(envs.ZENLOG_LOGGER_IN))
}

func WriteToOuter() {
	FailUnlessInZenlog()
	file := os.Getenv(envs.ZENLOG_OUTER_TTY)
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

func OuterTty() {
	FailUnlessInZenlog()
	fmt.Println(os.Getenv(envs.ZENLOG_OUTER_TTY))
}

func LoggerPipe() {
	FailUnlessInZenlog()
	fmt.Println(os.Getenv(envs.ZENLOG_LOGGER_IN))
}

func CheckUpdate() {
	if strconv.FormatInt(util.SelfCtime().Unix(), 10) == os.Getenv(envs.ZENLOG_BIN_CTIME) {
		util.ExitSuccess()
	}
	util.Say("Zenlog updated. Run \"zenlog_restart\" (or \"exit 13\") to restart a zenlog session.")
	util.ExitFailure()
}

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
		history.HistoryCommand(args)

	case "current-log":
		FailUnlessInZenlog()
		history.CurrentLogCommand(args)

	case "last-log":
		FailUnlessInZenlog()
		history.LastLogCommand(args)

	case "check-update":
		FailUnlessInZenlog()
		CheckUpdate()

		// TODO Refactor these commands for testability.
	case "start-command":
		FailUnlessInZenlog()
		if len(args) < 1 {
			util.Fatalf("start-command expects at least 1 argument.")
		}
		logger.StartCommand("", args[:], util.NewClock())

	case "start-command-with-env":
		FailUnlessInZenlog()
		if len(args) < 2 {
			util.Fatalf("start-command-with-env expects at least 2 arguments.")
		}
		logger.StartCommand(args[0], args[1:], util.NewClock())

	case "stop-log", "end-command":
		FailUnlessInZenlog()

		wantLineNumber := false
		i := 0
		if i < len(args) && args[i] == "-n" {
			i++
			wantLineNumber = true
		}
		exitStatus := -1
		var err error
		if i < len(args) {
			exitStatus, err = strconv.Atoi(args[i])
			util.Check(err, "Exit status must be integer; '%s' given.", args[i])
			i++
		}
		logger.EndCommand(exitStatus, wantLineNumber, util.NewClock())
	default:
		return
	}
	util.ExitSuccess()
}
