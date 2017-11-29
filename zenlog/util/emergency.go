package util

import "syscall"

func StartEmergencyShell() {
	Say("Starting emergency shell. (TODO: Reconsider this feature.)")

	shell := "/bin/sh"
	syscall.Exec(shell, StringSlice(shell), nil)
}