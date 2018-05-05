package util

import (
	"github.com/omakoto/go-common/src/common"
	"os"
	"syscall"
)

// MustExec is a must-version of syscall.Exec.
func MustExec(args []string) {
	common.Debugf("Executing: %v", args)

	err := syscall.Exec(args[0], args, os.Environ())
	common.Checkf(err, "Exec failed args=%v", args)
}
