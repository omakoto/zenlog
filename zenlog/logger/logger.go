package logger

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/kr/pty"
	"github.com/mattn/go-isatty"
	"github.com/omakoto/go-common/src/shell"
	"github.com/omakoto/go-common/src/textio"
	"github.com/omakoto/go-common/src/utils"
	"github.com/omakoto/zenlog/zenlog/config"
	"github.com/omakoto/zenlog/zenlog/envs"
	"github.com/omakoto/zenlog/zenlog/logfiles"
	"github.com/omakoto/zenlog/zenlog/util"
	"github.com/pkg/term"
)

type Logger struct {
	Config *config.Config

	child *exec.Cmd

	OuterTty  string
	master    *os.File
	stdinTerm *term.Term

	ForwardPipe *os.File
	ReversePipe *os.File

	startRequest        *StartRequest
	logFiles            *logfiles.LogFiles
	numLines            int
	hasDanglingLastLine bool

	sanitizer *textio.Sanitizer

	clock utils.Clock
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

	l.sanitizer = textio.NewSanitizer()

	l.OuterTty = util.Ttyname(os.Stdin.Fd())
	stdinTerm, err := term.Open(l.OuterTty)
	util.Check(err, "Cannot open tty '%s'", l.OuterTty)

	util.Debugf("stdinTerm=%+v", stdinTerm)
	l.stdinTerm = stdinTerm

	err = l.stdinTerm.SetRaw()
	util.Check(err, "SetRaw failed")
	util.SetOutputIsRaw(true)

	// Make the pipes.
	l.ForwardPipe = mustMakeFifo(config, "f")
	l.ReversePipe = mustMakeFifo(config, "r")

	l.clock = utils.NewClock()

	// Update config with the pipe names.
	config.LoggerIn = l.ForwardPipe.Name()
	config.LoggerOut = l.ReversePipe.Name()
	config.OuterTty = util.Tty()

	util.Dump("Logger=", l)

	return &l
}

func (l *Logger) exportEnviron() {
	os.Setenv(envs.ZenlogBin, util.FindZenlogBin())
	os.Setenv(envs.ZenlogSourceTop, util.ZenlogSrcTopDir())
	os.Setenv(envs.ZenlogBinCtime, strconv.FormatInt(util.ZenlogBinCtime().Unix(), 10))

	os.Setenv(envs.ZenlogDir, l.Config.LogDir)
	os.Setenv(envs.ZenlogPid, strconv.Itoa(l.Config.ZenlogPid))
	os.Setenv(envs.ZenlogOuterTty, l.OuterTty)
	os.Setenv(envs.ZenlogLoggerIn, l.ForwardPipe.Name())
	os.Setenv(envs.ZenlogLoggerOut, l.ReversePipe.Name())
}

// StartChild starts a child process in a new PTY.
func (l *Logger) StartChild() {
	l.exportEnviron()

	// Create a pty and start the child command.
	util.Debugf("Executing: %s", l.Config.StartCommand)
	l.child = exec.Command("/bin/sh", "-c",
		envs.ZenlogSignature+
			fmt.Sprintf("=\"$(tty)\":%s ", shell.Escape(Signature()))+
			l.Config.StartCommand)
	var err error
	l.master, err = pty.Start(l.child)
	util.Check(err, "Unable to create pty or execute /bin/sh")

	util.PropagateTerminalSize(os.Stdin, l.master)
}

// Child returns the child process.
func (l *Logger) Child() *exec.Cmd {
	return l.child
}

// Child returns the master tty.
func (l *Logger) Master() *os.File {
	return l.master
}

func (l *Logger) startForwarders() {
	m := l.master

	if l.Config.UseSplice {
		// Forward the input from stdin to the l.
		go forward(os.Stdin, m)

		// Read the output, and write to the STDOUT, and also to the pipe.
		go tee(m, l.ForwardPipe, os.Stdout)
	} else {
		go forwardSimple(os.Stdin, m)
		go teeSimple(m, l.ForwardPipe, os.Stdout)
	}
}

func (l *Logger) CleanUp() {
	if l.master != nil {
		l.master.Close()
	}

	l.stdinTerm.Restore()
	util.SetOutputIsRaw(false)

	l.ForwardPipe.Close()
	l.ReversePipe.Close()

	util.Warn(os.Remove(l.ForwardPipe.Name()), "Remove failed")
	util.Warn(os.Remove(l.ReversePipe.Name()), "Remove failed")
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

func (l *Logger) SendCloseRequest() {
	util.WriteToFile(l.Config.LoggerIn, utils.StringSlice(CloseSessionCommand))
}

func (l *Logger) SendFlushRequest() {
	util.WriteToFile(l.Config.LoggerIn, utils.StringSlice(FlushCommand))
}

// Open log files.
func (l *Logger) openLogs(request *StartRequest) {
	// If the previous log is still open, close it.
	l.closeLogs(nil)

	// Re-init the fields for the command.
	l.startRequest = request
	l.logFiles = &request.LogFiles

	l.logFiles.Open(false)

	l.write([]byte("$ " + request.Command.CommandLine + "\n"))

	l.numLines = 0 // Don't count the first line. Start with 0 here.
	l.hasDanglingLastLine = false

	// Check nolog.
	if request.Command.NoLog {
		l.write([]byte("[reducted]"))
		l.closeLogs(nil)

		// HACK: This is to update the injected clock even for reducted commands.
		util.GetInjectedNow(l.clock)
	}
}

// Close log files.
func (l *Logger) closeLogs(req *StopRequest) {
	if !l.isOpen() {
		return
	}
	if req != nil {
		l.logFiles.WriteFinishToEnv(req.ExitStatus, l.startRequest.StartTime, util.GetInjectedNow(l.clock))
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
	if l.isOpen() {
		l.logFiles.San.Flush()
		l.logFiles.Raw.Flush()
	}
}

// Start the forwarders, and do the main loop.
func (l *Logger) DoLogger() {
	l.startForwarders()

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
				case CloseSessionCommand:
					l.closeLogs(nil)
					return

				case FlushCommand:
					l.flush()
					continue

				case CommandStartCommand:
					if len(args) != 2 {
						util.Say("Invalid number of args (%d) for %s.", len(args), CommandStartCommand)
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
				case CommandEndCommand:
					if len(args) != 3 {
						util.Say("Invalid number of args (%d) for %s.", len(args), CommandEndCommand)
						continue
					}
					fingerprint := args[1]

					// Parse request.
					req := StopRequest{}
					if !util.TryUnmarshal(args[2], &req) {
						continue
					}
					util.Dump("StopRequest=", req)

					// Close log.
					l.closeLogs(&req)

					// Send reply.
					l.MustReply(l.Config, utils.StringSlice(CommandEndCommand, fingerprint, util.MustMarshal(StopReply{l.numLines})))
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
	l.SendCloseRequest()
}
