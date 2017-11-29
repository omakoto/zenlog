package logger

import (
	"bufio"
	"fmt"
	"github.com/mattn/go-isatty"
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

	startRequest        *StartRequest
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

	// Update config with the pipe names.
	config.LoggerIn = l.ForwardPipe.Name()
	config.LoggerOut = l.ReversePipe.Name()
	config.OuterTty = util.Tty()

	return &l
}

func (l *Logger) ExportEnviron() {
	os.Setenv(envs.ZENLOG_DIR, l.Config.LogDir)
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
func (l *Logger) openLogs(request *StartRequest) {
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
func (l *Logger) closeLogs(req *StopRequest) {
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

	if l.Config.AutoFlush {
		l.flush()
	}
}
func (l *Logger) flush() {
	if !l.isOpen() {
		return
	}
	l.logFiles.Raw.Flush()
	l.logFiles.San.Flush()
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
				case CHILD_DIED_COMMAND:
					util.Debugf("Child died.")
					l.closeLogs(nil)
					return

				case FLUSH_COMMAND:
					util.Debugf("Flushing...")
					l.flush()
					continue

				case COMMAND_START_COMMAND:
					if len(args) != 2 {
						util.Say("Invalid number of args (%d) for %s.", len(args), COMMAND_START_COMMAND)
						continue
					}
					// Parse request.

					req := StartRequest{}
					if !util.TryUnmarshal(args[1], &req) {
						continue
					}
					util.Dump("StartRequest=", req)

					// Open log.
					l.openLogs(&req)
					continue
				case COMMAND_END_COMMAND:
					if len(args) != 3 {
						util.Say("Invalid number of args (%d) for %s.", len(args), COMMAND_END_COMMAND)
						continue
					}
					fingerprint := args[1]

					// Parse request.
					req := StopRequest{}
					if !util.TryUnmarshal(args[1], &req) {
						continue
					}
					util.Dump("StopRequest=", req)

					// Close log.
					l.closeLogs(&req)

					// Send reply.
					l.MustReply(l.Config, util.StringSlice(COMMAND_END_COMMAND, fingerprint, util.MustMarshal(StopReply{l.numLines})))
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

func (l *Logger) OnChildDied() {
	args := make([]string, 1)
	args[0] = CHILD_DIED_COMMAND
	MustSendToLogger(l.Config, args)
}
