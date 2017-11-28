package logger

import (
	"bufio"
	"fmt"
	"github.com/mattn/go-isatty"
	"github.com/omakoto/zenlog-go/zenlog/builtins"
	"github.com/omakoto/zenlog-go/zenlog/config"
	"github.com/omakoto/zenlog-go/zenlog/envs"
	"github.com/omakoto/zenlog-go/zenlog/logfiles"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"github.com/pkg/term"
	"io"
	"os"
	"regexp"
	"strconv"
	"syscall"
)

type Logger struct {
	Config *config.Config

	OuterTty  string
	stdinTerm *term.Term

	ForwardPipe *os.File
	ReversePipe *os.File

	startRequest        *builtins.StartRequest
	logFiles            *logfiles.LogFiles
	numLines            int
	hasDanglingLastLine bool

	rePrefixCommands      *regexp.Regexp
	reAlwaysNoLogCommands *regexp.Regexp

	sanitizer *util.Sanitizer

	clock util.Clock
}

func mustMakeFifo(config *config.Config, suffix string) *os.File {
	filename := fmt.Sprintf("%szenlog.%d%s.pipe", config.TempDir, config.ZenlogPid, suffix)
	os.Remove(filename)

	util.Debugf("Making fifo '%s'...", filename)
	err := syscall.Mkfifo(filename, 0600)
	util.Check(err, "Makefifo failed for '%s'", filename)

	file, err := os.OpenFile(filename, os.O_RDWR, 0600)
	util.Check(err, "OpenFile failed for '%s'", filename)
	return file
}

func NewLogger(config *config.Config) *Logger {
	if !isatty.IsTerminal(os.Stdin.Fd()) {
		util.Fatalf("Stdin must be terminal.")
	}

	l := Logger{Config: config}

	l.rePrefixCommands = regexp.MustCompile(config.PrefixCommands)
	l.reAlwaysNoLogCommands = regexp.MustCompile(config.AlwaysNoLogCommands)
	l.sanitizer = util.NewSanitizer()

	l.OuterTty = util.Ttyname(os.Stdin.Fd())
	stdinTerm, err := term.Open(l.OuterTty)
	util.Check(err, "Cannot open tty '%s'", l.OuterTty)

	util.Debugf("stdinTerm=%+v", stdinTerm)
	l.stdinTerm = stdinTerm

	err = l.stdinTerm.SetRaw()
	util.Check(err, "SetRaw failed")
	util.SetOutputIsRaw()

	// Make the pipes.
	l.ForwardPipe = mustMakeFifo(config, "f")
	l.ReversePipe = mustMakeFifo(config, "r")

	l.clock = util.NewClock()

	return &l
}

func (l *Logger) ExportEnviron() {
	os.Setenv(envs.ZENLOG_PID, strconv.Itoa(os.Getpid()))
	os.Setenv(envs.ZENLOG_OUTER_TTY, l.OuterTty)
	os.Setenv(envs.ZENLOG_LOGGER_IN, l.ForwardPipe.Name())
	os.Setenv(envs.ZENLOG_LOGGER_OUT, l.ReversePipe.Name())
}

func (l *Logger) CleanUp() {
	l.stdinTerm.Restore()

	l.ForwardPipe.Close()
	l.ReversePipe.Close()

	os.Remove(l.ForwardPipe.Name())
	os.Remove(l.ReversePipe.Name())
}

func (l *Logger) MustReply(config *config.Config, vals []string) {
	reply := util.Encode(vals)
	util.Debugf("Replying: %v", vals)
	_, err := l.ReversePipe.WriteString(reply)
	util.Check(err, "Failed to reply from logger")
}

func (l *Logger) isOpen() bool {
	return l.logFiles != nil
}

// Open log files.
func (l *Logger) openLogs(request *builtins.StartRequest) {
	// If the previous log is still open, close it.
	l.closeLogs(nil)

	// Re-init the fields for the command.
	l.startRequest = request
	l.logFiles = &request.LogFiles

	l.logFiles.Open()

	l.write([]byte("$ " + request.Command.CommandLine + "\n"))

	l.numLines = 0 // Don't count the first line. Start with 0 here.
	l.hasDanglingLastLine = false
}

// Close log files.
func (l *Logger) closeLogs(req *builtins.StopRequest) {
	if !l.isOpen() {
		return
	}
	if req != nil {
		l.logFiles.WriteFinishToEnv(req.ExitStatus, l.startRequest.StartTime, l.clock.Now())
	}
	l.logFiles.Close()

	l.startRequest = nil
	l.logFiles = nil

	// If the last line didn't finish with NL, then add one line.
	if l.hasDanglingLastLine {
		l.numLines++
		l.hasDanglingLastLine = false
	}
}

// Write a log line.
func (l *Logger) write(line []byte) {
	if !l.isOpen() || len(line) == 0 {
		return
	}
	_, err := l.logFiles.Raw.Write(line)
	util.Warn(err, "Write failed")

	_, err = l.logFiles.San.Write(l.sanitizer.Sanitize(line))
	util.Warn(err, "Write failed")

	l.numLines += util.NumLines(line)
	l.hasDanglingLastLine = line[len(line)-1] != '\n'
}

func (l *Logger) DoLogger() {
	bout := bufio.NewReader(l.ForwardPipe)
	for {
		line, err := bout.ReadBytes('\n')
		if len(line) != 0 {
			decoded, pre, args := util.TryDecodeBytes(line)
			if !decoded {
				l.write(line)
			} else {
				l.write(pre)

				if len(args) == 0 {
					util.Say("Received empty command.")
					continue
				}
				switch args[0] {
				case builtins.COMMAND_START_COMMAND:
					if len(args) != 2 {
						util.Say("Invalid number of args (%d) for %s.", len(args), builtins.COMMAND_START_COMMAND)
						continue
					}
					// Parse request.
					req, err := builtins.DecodeStartRequest(args[1])
					if util.Warn(err, "Decode error") {
						continue
					}
					util.Dump("StartRequest=", req)

					// Open log.
					l.openLogs(req)
					continue
				case builtins.COMMAND_END_COMMAND:
					if len(args) != 3 {
						util.Say("Invalid number of args (%d) for %s.", len(args), builtins.COMMAND_END_COMMAND)
						continue
					}
					// Parse request.
					req, err := builtins.DecodeStopRequest(args[2])
					if util.Warn(err, "Decode error") {
						continue
					}

					// Close log.
					l.closeLogs(req)

					// Send reply.
					r := builtins.StopReply{l.numLines}

					fingerprint := args[1]
					rep := make([]string, 3)
					rep[0] = builtins.COMMAND_END_COMMAND
					rep[1] = fingerprint
					rep[2] = string(r.MustEncode())
					l.MustReply(l.Config, rep)
					continue
				}
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			util.Fatalf("ReadString failed: %s", err)
		}
	}
}
