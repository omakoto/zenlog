package shell

import (
	"testing"

	"github.com/omakoto/zenlog-go/zenlog/util"
)

func TestShellSplit(t *testing.T) {
	inputs := []struct {
		source   string
		expected []string
	}{
		{"", util.Ar()},
		{"a", util.Ar("a")},
		{"aaa", util.Ar("aaa")},
		{"aaa b  ccc", util.Ar("aaa", "b", "ccc")},
	}
	for _, v := range inputs {
		actual := ShellSplit(v.source)
		util.AcssertStringSlicesEqual(t, "Source="+v.source, v.expected, actual)
	}
}
