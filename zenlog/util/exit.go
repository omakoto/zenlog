package util

type exitStatus struct {
	code int
}

func ExitSuccess() {
	Exit(true)
}

func ExitFailure() {
	Exit(false)
}

func Exit(success bool) {
	status := 1
	if success {
		status = 0
	}
	panic(exitStatus{status})
}

func RunWithRescue(f func() int) (result int) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(exitStatus); ok {
				result = e.code
				return
			}
			panic(r)
		}
	}()
	return f()
}