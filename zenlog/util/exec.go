package util

import (
	"os"
	"syscall"
)

// MustExec is a must-version of syscall.Exec.
func MustExec(args []string) {
	Debugf("Executing: %v", args)

	err := syscall.Exec(args[0], args, os.Environ())
	Check(err, "Exec failed args=%v", args)
}
