package util

import (
	"os"
)

// GetSubcommand extracts a zenlog subcommand name from os.Args.
func GetSubcommand() (command string, args []string) {
	Debugf("os.Args=%+v", os.Args)
	if len(os.Args) == 1 {
		command = ""
		args = os.Args[1:]
	} else {
		command = os.Args[1]
		args = os.Args[2:]
	}
	Debugf("command='%s', args=%+v", command, args)
	return
}
