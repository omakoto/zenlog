package shell

import (
	"testing"
)

func TestShescape(t *testing.T) {
	inputs := []struct {
		expected string
		source   string
	}{
		{"", ""},
		{"abc", "abc"},
		{"'abc '", "abc "},
		{"'abc def \" '\\'' xyz '\\'''", "abc def \" ' xyz '"},
	}
	for _, v := range inputs {
		actual := Shescape(v.source)
		if v.expected != actual {
			t.Errorf("Expected=%s, actual=%s", v.expected, actual)
		}
	}
}
