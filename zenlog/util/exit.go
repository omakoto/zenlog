package util

import "os"

type exitStatus struct {
	code int
}

// ExitSuccess should be used within RunAndExit to cleanly finishes the process with a success code.
func ExitSuccess() {
	Exit(true)
}

// ExitFailure should be used within RunAndExit to cleanly finishes the process with a failure code.
func ExitFailure() {
	Exit(false)
}

// Exit should be used within RunAndExit to cleanly finishes the process.
func Exit(success bool) {
	status := 1
	if success {
		status = 0
	}
	panic(exitStatus{status})
}

func runWithRescue(f func() int) (result int) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(exitStatus); ok {
				result = e.code
			} else {
				panic(r)
			}
		}
	}()
	result = f()
	return
}

// RunAndExit executes a given function. Within the function, util.Exit* functions can be used to finish the process cleanly.
func RunAndExit(f func() int) {
	os.Exit(runWithRescue(f))
}
