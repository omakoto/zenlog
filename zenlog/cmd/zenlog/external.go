package main

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/omakoto/zenlog/zenlog/config"
	"github.com/omakoto/zenlog/zenlog/util"
)

func tryRunExternalCommand(path string, command string, args []string, setX bool) {
	f, err := filepath.Abs(path + "/zenlog-" + command)
	util.Check(err, "Abs failed")

	util.Debugf("Checking %s", f)

	stat, err := os.Stat(f)
	if err != nil {
		return // File doesn't exist
	}

	if (stat.Mode() & syscall.S_IXUSR) == 0 {
		if !setX {
			return
		}
		// Files under ~/go/pkg/mod/ doesn't have X bits set, so just set them as needed.
		err = os.Chmod(f, stat.Mode()|syscall.S_IXUSR|syscall.S_IXGRP|syscall.S_IXOTH)
		util.Check(err, "Unable to set X bits on file '%s'", f)
	}

	execArgs := make([]string, 0, len(args)+1)
	execArgs = append(execArgs, f)
	execArgs = append(execArgs, args...)

	util.MustExec(execArgs)
}

// MaybeRunExternalCommand looks for "zenlog-SUBCOMMAND", first in the zenlog "subcommand" directory,
// then in PATH and executes it.
func MaybeRunExternalCommand(command string, args []string) {
	tryRunExternalCommand(config.ZenlogSrcTopDir()+"/subcommands", command, args, true)

	for _, path := range strings.Split(os.Getenv("PATH"), ":") {
		tryRunExternalCommand(path, command, args, false)
	}
}
