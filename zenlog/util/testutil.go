package util

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Ar(a ...string) []string {
	return a
}

func SlicesEqual(a []string, b[] string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true;
}

func AssertStringsEqual(t *testing.T, input string, expected string, actual string) {
	if expected != actual {
		t.Errorf("input=%s expected=%s actual=%s", input, expected, actual)
	}
}

func AssertStringSlicesEqual(t *testing.T, input string, expected []string, actual []string) {
	if !SlicesEqual(expected, actual) {
		t.Errorf("input=%s expected=%s actual=%s", input, spew.Sdump(expected), spew.Sdump(actual))
	}
}

func AssertFileExist(t *testing.T, file string) {
	if !FileExists(file) {
		t.Errorf("File %s not createtd.", file)
	}
}
