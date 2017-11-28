package zenlog

import (
	"os"
	"strings"
	"syscall"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"path/filepath"
)

func tryRunExtetrnalCommand(path string, command string, args []string) {
	f := path + "/zenlog-" + command

	util.Debugf("Checking %s", f)

	stat, err := os.Stat(f)
	if (err == nil) && ((stat.Mode() & syscall.S_IXUSR) != 0) {
		execArgs := make([]string, 0, len(args) + 1)
		execArgs = append(execArgs, f)
		execArgs = append(execArgs, args...)

		util.Debugf("Executing: %s, %v", f, execArgs)

		err := syscall.Exec(f, execArgs, os.Environ())
		util.Check(err, "Exec failed")
	}
}

// Look for "zenlog-SUBCOMMAND" in PATH and execute it.
func MaybeRunExtetrnalCommand(command string, args []string) {
	exePath := util.FindSelf()

	tryRunExtetrnalCommand(filepath.Dir(exePath) + "/../subcommands", command, args)
	tryRunExtetrnalCommand(filepath.Dir(exePath) + "/../github.com/omakoto/zenlog-go/subcommands", command, args)

	for _, path := range strings.Split(os.Getenv("PATH"), ":") {
		tryRunExtetrnalCommand(path, command, args)
	}
}
