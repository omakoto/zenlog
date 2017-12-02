package zenlog

import (
	"fmt"
	"github.com/kr/pty"
	"github.com/omakoto/zenlog-go/zenlog/config"
	"github.com/omakoto/zenlog-go/zenlog/envs"
	"github.com/omakoto/zenlog-go/zenlog/logger"
	"github.com/omakoto/zenlog-go/zenlog/shell"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

const resurrectCode = 13

func maybeStartEmergencyShell(startTime time.Time, r interface{}, childStatus int) {
	if r == nil && childStatus <= 0 {
		return // okay
	}
	startShell := false
	if r != nil {
		util.Say("Panic detected: %v", r)
		startShell = true
	} else {
		// If the child dies too early, something may be wrong, so start a shell...
		if util.NewClock().Now().Sub(startTime).Seconds() < 30 {
			util.Say("Child finished unsuccessfully, too soon?: code=%d%s", childStatus)
			startShell = true
		}
	}
	if startShell {
		util.StartEmergencyShell()
	}
}

// StartZenlog starts a new zenlog session.
func StartZenlog(args []string) (commandExitCode int, resurrect bool) {
	var childStatus = -1

	startTime := util.NewClock().Now()
	defer func() {
		if !resurrect {
			maybeStartEmergencyShell(startTime, recover(), childStatus)
		}
	}()

	config := config.InitConfigiForLogger()

	fmt.Printf("Zenlog starting... [ZENLOG_DIR=%s ZENLOG_PID=%d]\n", config.LogDir, config.ZenlogPid)

	l := logger.NewLogger(config)
	defer l.CleanUp()

	// Set up signal handler.
	sigch := make(chan os.Signal)
	signal.Notify(sigch, syscall.SIGCHLD, syscall.SIGWINCH, syscall.SIGHUP)

	// Set up environmental variables.
	l.ExportEnviron()

	// Create a pty and start the child command.
	util.Debugf("Executing: %s", config.StartCommand)
	c := exec.Command("/bin/sh", "-c",
		envs.ZenlogSignature+
			fmt.Sprintf("=\"$(tty)\":%s ", shell.Escape(Signature()))+
			config.StartCommand)
	m, err := pty.Start(c)
	util.Check(err, "Unable to create pty or execute /bin/sh")
	defer m.Close()

	util.PropagateTerminalSize(os.Stdin, m)

	// Signal handler.
	go func() {
		for s := range sigch {
			switch s {
			case syscall.SIGWINCH:
				util.Debugf("Caught SIGWINCH")

				util.PropagateTerminalSize(os.Stdin, m)
				l.SendFlushRequest()

			case syscall.SIGHUP:
				util.Debugf("Caught SIGHUP")

				l.SendCloseRequest()

			case syscall.SIGCHLD:
				util.Debugf("Caught SIGCHLD")
				ps, err := c.Process.Wait()
				if err != nil {
					util.Warn(err, "Wait failed")
					childStatus = 255
				} else {
					childStatus = ps.Sys().(syscall.WaitStatus).ExitStatus()
				}
				l.OnChildDied()
			default:
				util.Say("Caught unexpected signal: %+v", s)
			}
		}
	}()

	// Forward the input from stdin to the l.
	go func() {
		io.Copy(m, os.Stdin)
	}()

	// Read the output, and write to the STDOUT, and also to the pipe.
	go func() {
		buf := make([]byte, 32*1024)

		for {
			nr, er := m.Read(buf)
			if nr > 0 {
				// First, write to stdout.
				nw, ew := os.Stdout.Write(buf[0:nr])
				if util.Warn(ew, "Stdout.Write failed") {
					break
				}
				if nr != nw {
					err = io.ErrShortWrite
					break
				}
				// Then, write to l.
				nw, ew = l.ForwardPipe.Write(buf[0:nr])
				if util.Warn(ew, "ForwardPipe.Write failed") {
					break
				}
				if nr != nw {
					err = io.ErrShortWrite
					break
				}
			}
			if er != nil {
				break
			}
		}
	}()
	// Logger.
	l.DoLogger()

	util.Debugf("Zenlog exiting with=%d", childStatus)
	if childStatus == resurrectCode {
		return 0, true
	}

	return childStatus, false
}
