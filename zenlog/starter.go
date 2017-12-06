package zenlog

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/omakoto/zenlog-go/zenlog/config"
	"github.com/omakoto/zenlog-go/zenlog/logger"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"runtime/pprof"
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
				fmt.Println("SIGUSR2")

				p := pprof.Lookup("goroutine")
				l.RunWithCookedTerminal(func() {
					p.WriteTo(os.Stdout, 1)
				})

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

	startTime := util.NewClock().Now()
	defer func() {
		if !resurrect {
			maybeStartEmergencyShell(startTime, recover(), childStatus)
		}
	}()

	config := config.InitConfigForLogger()

	// Create a logger and start a child.
	fmt.Printf("Zenlog starting... [ZENLOG_DIR=%s ZENLOG_PID=%d]\n", config.LogDir, config.ZenlogPid)

	l := logger.NewLogger(config)
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
