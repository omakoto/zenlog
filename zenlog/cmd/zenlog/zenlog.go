package main

import (
	"github.com/omakoto/zenlog-go/zenlog"
	"github.com/omakoto/zenlog-go/zenlog/builtins"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"os"
)

func main() {
	command, args := util.GetSubcommand()

	if command == "" {
		os.Exit(zenlog.StartZenlog(args))
	}
	builtins.MaybeRunBuiltin(command, args)
	zenlog.MaybeRunExtetrnalCommand(command, args)

	util.Fatalf("Unknown subcommand: '%s'", command)
}
