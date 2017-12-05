package builtins

import (
	"flag"
	"strconv"

	"github.com/omakoto/zenlog-go/zenlog/logger"
	"github.com/omakoto/zenlog-go/zenlog/util"
)

// startCommand tells zenlog to start logging for a command.
func startCommand(args []string) {
	flags := flag.NewFlagSet("zenlog start-command", flag.ExitOnError)
	e := flags.String("e", "", "Pass string to write to ENV file")
	flags.Parse(args)
	args = flags.Args()

	if len(args) < 1 {
		util.Fatalf("start-command expects at least 1 argument.")
	}
	logger.StartCommand(*e, args[:], util.NewClock())
}

// endCommand tells zenlog to stop logging for the current command.
func endCommand(args []string) {
	flags := flag.NewFlagSet("zenlog end-command", flag.ExitOnError)
	wantLineNumber := flags.Bool("n", false, "Print number of lines in log")
	flags.Parse(args)
	args = flags.Args()

	exitStatus := -1
	var err error

	if len(args) > 0 {
		exitStatus, err = strconv.Atoi(args[0])
		util.Check(err, "Exit status must be integer; '%s' given.", args[0])
	}
	logger.EndCommand(exitStatus, *wantLineNumber, util.NewClock())
}
