package util

import "testing"

func CheckRunWithRescue(t *testing.T, expected int, f func() int) {
	actual := RunWithRescue(f)
	if expected != actual {
		t.Errorf("Expected=%d actual=%d f=%v\n", expected, actual, f)
	}
}

func TestRunWithRescue(t *testing.T) {
	CheckRunWithRescue(t, 0, func() int {return 0})
	CheckRunWithRescue(t, 5, func() int {return 5})
	CheckRunWithRescue(t, 0, func() int {ExitSuccess(); return -1})
	CheckRunWithRescue(t, 1, func() int {ExitFailure(); return -1})
}
