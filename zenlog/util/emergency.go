package util

import "syscall"

func StartEmergencyShell() {
	Say("Something may be wrong; starting emergency shell...")

	shell := "/bin/sh"
	syscall.Exec(shell, StringSlice(shell), nil)
}