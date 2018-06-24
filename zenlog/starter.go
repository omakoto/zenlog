package zenlog

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/omakoto/go-common/src/utils"
	"github.com/omakoto/zenlog/zenlog/config"
	"github.com/omakoto/zenlog/zenlog/logger"
	"github.com/omakoto/zenlog/zenlog/util"
	"runtime/pprof"
)

const resurrectCode = 13

func maybeStartEmergencyShell(c *config.Config, startTime time.Time, r interface{}, childStatus int) {
	if r == nil && childStatus <= 0 {
		return // okay
	}
	startShell := false
	if r != nil {
		util.Say("Panic detected: %v", r)
		startShell = true
	} else {
		threshold := 30.0
		threshold = float64(c.CriticalCrashMaxSeconds)

		// If the child dies too early, something may be wrong, so start a shell...
		if utils.NewClock().Now().Sub(startTime).Seconds() < threshold {
			util.Say("Child finished unsuccessfully, too soon?: code=%d%s", childStatus)
			startShell = true
		}
	}
	if startShell {
		util.StartEmergencyShell()
	}
}

func dumpAllGoroutines() {
	b := bytes.Buffer{}

	p := pprof.Lookup("goroutine")
	p.WriteTo(&b, 1)

	util.Say(b.String())
}

func setupSignalHandler(l *logger.Logger, childStatus *int) {
	sigch := make(chan os.Signal)
	signal.Notify(sigch, syscall.SIGCHLD, syscall.SIGWINCH, syscall.SIGHUP, syscall.SIGUSR2)

	// Signal handler.
	go func() {
		for s := range sigch {
			switch s {
			case syscall.SIGWINCH:
				util.Debugf("Caught SIGWINCH")

				util.PropagateTerminalSize(os.Stdin, l.Master())
				l.SendFlushRequest()

			case syscall.SIGHUP:
				util.Debugf("Caught SIGHUP")

				l.SendCloseRequest()

			case syscall.SIGCHLD:
				util.Debugf("Caught SIGCHLD")
				ps, err := l.Child().Process.Wait()
				if err != nil {
					util.Warn(err, "Wait failed")
					*childStatus = 255
				} else {
					*childStatus = ps.Sys().(syscall.WaitStatus).ExitStatus()
				}
				l.OnChildDied()

			case syscall.SIGUSR2:
				util.Say("Caught SIGUSR2; dumping stacktraces.")
				dumpAllGoroutines()

			default:
				util.Say("Caught unexpected signal: %+v", s)
			}
		}
	}()
}

// StartZenlog starts a new zenlog session.
func StartZenlog(args []string) (commandExitCode int, resurrect bool) {
	// Initialize.
	var childStatus = -1

	var c *config.Config

	startTime := utils.NewClock().Now()
	defer func() {
		if !resurrect {
			maybeStartEmergencyShell(c, startTime, recover(), childStatus)
		}
	}()

	c = config.InitConfigForLogger()

	// Create a logger and start a child.
	fmt.Printf("Zenlog starting... [ZENLOG_DIR=%s ZENLOG_PID=%d]\n", c.LogDir, c.ZenlogPid)

	l := logger.NewLogger(c)
	defer l.CleanUp()

	l.StartChild()

	// Set up a signal handler
	setupSignalHandler(l, &childStatus)

	// Start the logger.
	l.DoLogger()

	// Child finished. (or maybe SIGHUP)
	util.Debugf("Zenlog exiting with=%d", childStatus)
	if childStatus == resurrectCode {
		return 0, true
	}

	return childStatus, false
}
