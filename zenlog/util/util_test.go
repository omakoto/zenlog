package util

import "testing"

func TestNumLines(t *testing.T) {
	inputs := []struct {
		s     string
		lines int
	}{
		{"", 0},
		{"\n", 1},
		{"\n\n", 2},
		{"  \n  \n  ", 2},
	}
	for _, v := range inputs {
		actual := NumLines([]byte(v.s))
		if v.lines != actual {
			t.Errorf("Input=%s expected=%d, actual=%d", v.s, v.lines, actual)
		}
	}
}
