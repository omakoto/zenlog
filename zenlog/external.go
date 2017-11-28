package zenlog

import (
	"os"
	"strings"
	"syscall"
	"github.com/omakoto/zenlog-go/zenlog/util"
)

// Look for "zenlog-SUBCOMMAND" in PATH and execute it.
func MaybeRunExtetrnalCommand(command string, args []string) {
	for _, path := range strings.Split(os.Getenv("PATH"), ":") {
		f := path + "/zenlog-" + command
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
}
