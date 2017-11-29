package main

import (
	"github.com/omakoto/zenlog-go/zenlog"
	"github.com/omakoto/zenlog-go/zenlog/builtins"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"os"
)

func realMain() int {
	command, args := util.GetSubcommand()

	if command == "" {
		return zenlog.StartZenlog(args)
	}
	builtins.MaybeRunBuiltin(command, args)
	zenlog.MaybeRunExtetrnalCommand(command, args)

	util.Fatalf("Unknown subcommand: '%s'", command)
	return 0
}

func main() {
	os.Exit(util.RunWithRescue(realMain))
}
