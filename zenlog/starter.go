package zenlog

import (
	"fmt"
	"github.com/kr/pty"
	"github.com/omakoto/zenlog-go/zenlog/config"
	"github.com/omakoto/zenlog-go/zenlog/envs"
	"github.com/omakoto/zenlog-go/zenlog/logger"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func maybeStartEmergencyShell(startTime time.Time, r interface{}, childStatus int) {
	if r == nil && childStatus <= 0 {
		return // okay
	}
	startShell := false
	if r != nil {
		util.Say("Panic detected: %v", r)
		startShell = true
	} else {
		util.Say("Child finished unsuccessfully: code=%d", childStatus)

		// If the child dies too early, something may be wrong, so start a shell...
		startShell = (util.NewClock().Now().Sub(startTime).Seconds() < 30)
	}
	if startShell {
		util.StartEmergencyShell()
	}
}

func StartZenlog(args []string) int {
	var childStatus int = -1

	startTime := util.NewClock().Now()
	defer func() {
		maybeStartEmergencyShell(startTime, recover(), childStatus)
	}()

	config := config.InitConfigiForLogger()
	util.Dump("config=", config)

	logger := logger.NewLogger(config)
	defer logger.CleanUp()
	util.Dump("Logger=", logger)

	// Set up signal handler.
	sigch := make(chan os.Signal)
	signal.Notify(sigch, syscall.SIGCHLD, syscall.SIGWINCH)

	// Set up environmental variables.
	logger.ExportEnviron()

	// Create a pty and start the child command.
	util.Debugf("Executing: %s", config.StartCommand)
	c := exec.Command("/bin/sh", "-c",
		envs.ZENLOG_SIGNATURE+
			fmt.Sprintf("=\"$(tty)\":%s ", util.Shescape(util.Signature()))+
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
				logger.SendFlushRequest()

			case syscall.SIGCHLD:
				util.Debugf("Caught SIGCHLD")
				ps, err := c.Process.Wait()
				if err != nil {
					util.Fatalf("Wait failed: %s", err)
				} else {
					childStatus = ps.Sys().(syscall.WaitStatus).ExitStatus()
				}
				logger.OnChildDied()
			default:
				util.Debugf("Caught unexpected signal: %+v", s)
			}
		}
	}()

	// Forward the input from stdin to the logger.
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
				// Then, write to logger.
				nw, ew = logger.ForwardPipe.Write(buf[0:nr])
				if util.Warn(ew, "Stdout.Write failed") {
					break
				}
				if nr != nw {
					err = io.ErrShortWrite
					break
				}
			}
			if er != nil {
				break // Ignore read error.
			}
		}
	}()
	// Logger.
	logger.DoLogger()

	util.Debugf("Zenlog exitting with=%d", childStatus)
	return childStatus
}
