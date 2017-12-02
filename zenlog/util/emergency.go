package util

import "syscall"

func StartEmergencyShell() {
	Say("Starting emergency shell...")

	shell := "/bin/sh"
	syscall.Exec(shell, StringSlice(shell), nil)
}
