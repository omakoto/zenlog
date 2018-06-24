package main

import (
	"github.com/omakoto/zenlog/zenlog/util"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func tryRunExternalCommand(path string, command string, args []string) {
	f, err := filepath.Abs(path + "/zenlog-" + command)
	util.Check(err, "Abs failed")

	util.Debugf("Checking %s", f)

	stat, err := os.Stat(f)
	if (err == nil) && ((stat.Mode() & syscall.S_IXUSR) != 0) {
		execArgs := make([]string, 0, len(args)+1)
		execArgs = append(execArgs, f)
		execArgs = append(execArgs, args...)

		util.MustExec(execArgs)
	}
}

// MaybeRunExternalCommand looks for "zenlog-SUBCOMMAND", first in the zenlog "subcommand" directory,
// then in PATH and executes it.
func MaybeRunExternalCommand(command string, args []string) {
	tryRunExternalCommand(util.ZenlogSrcTopDir()+"/subcommands", command, args)

	for _, path := range strings.Split(os.Getenv("PATH"), ":") {
		tryRunExternalCommand(path, command, args)
	}
}
