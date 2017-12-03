package builtins

import (
	"flag"

	"github.com/omakoto/zenlog-go/zenlog/logger"
	"github.com/omakoto/zenlog-go/zenlog/util"
)

// StartCommand tells zenlog to start logging for a command.
func startCommand(args []string) {
	flags := flag.NewFlagSet("zenlog start", flag.ExitOnError)
	e := flags.String("e", "", "Pass string to write to ENV file")
	flags.Parse(args)

	if len(args) < 1 {
		util.Fatalf("start-command expects at least 1 argument.")
	}
	logger.StartCommand(*e, args[:], util.NewClock())
}

// StartWithEnvCommand tells zenlog to start logging for a command.
func startWithEnvCommand(args []string) {
	if len(args) < 2 {
		util.Fatalf("start-command-with-env expects at least 2 arguments.")
	}
	logger.StartCommand(args[0], args[1:], util.NewClock())
}
