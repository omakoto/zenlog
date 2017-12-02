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
		{"aaa 'b  b'  ccc", util.Ar("aaa", "b  b", "ccc")},
		{`aaa 'b  b'\''  d'  ccc`, util.Ar("aaa", "b  b'  d", "ccc")},
		{"`ab\"`  ccc", util.Ar("ab\"", "ccc")},
		{`a\ \'\ \"`, util.Ar(`a ' "`)},
		{`$HOME/abc`, util.Ar(`$HOME/abc`)},
		{`${HOME}/abc`, util.Ar(`${HOME}/abc`)},
		{`  $(cat  ok  "$(next   "de  f")")/abc  xyz`, util.Ar(`$(cat  ok  $(next   de  f))/abc`, "xyz")},
		{`$ \`, util.Ar(`$`)},
		{`$`, util.Ar(`$`)},
		{`"\`, util.Ar()},
		{`cat|&grep`, util.Ar(`cat`, `|&`, `grep`)},
	}
	for _, v := range inputs {
		actual := ShellSplit(v.source)
		util.AssertStringSlicesEqual(t, "Source="+v.source, v.expected, actual)
	}
}
