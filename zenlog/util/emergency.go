package util

import "syscall"

// StartEmergencyShell exec's /bin/sh.
func StartEmergencyShell() {
	Say("Starting emergency shell...")

	shell := "/bin/sh"
	syscall.Exec(shell, StringSlice(shell), nil)
}
