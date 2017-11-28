package logfiles

import (
	"github.com/omakoto/zenlog-go/zenlog/config"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"testing"
)

func TestSplitComment(t *testing.T) {
	config := config.Config{}
	tests := []struct {
		input   string
		pre     string
		comment string
	}{
		{"", "", ""},
		{"abc # def", "abc", "def"},
		{"abc;#def", "abc", "def"},
		{"abc|#def", "abc", "def"},
		{"abc&#def", "abc", "def"},
		{"abc)#def", "abc", "def"},
		{"abc(#def", "abc", "def"},
		{"abc}#def", "abc", "def"},
		{"abc{#def", "abc", "def"},
	}
	for _, v := range tests {
		pre, comment := splitComment(&config, v.input)
		if v.pre != pre {
			t.Errorf("input=%s expected=%s actual=%s", v.input, v.pre, pre)
		}
		if v.comment != comment {
			t.Errorf("input=%s expected=%s actual=%s", v.input, v.comment, comment)
		}
	}
}

func TestParseCommandLine(t *testing.T) {
	config := config.Config{}
	tests := []struct {
		input   string
		exes    []string
		comment string
	}{
		{"", util.Ar(), ""},
		{"abc def #xyz", util.Ar("abc"), "xyz"},
		{"abc def|ABC DEF #xyz", util.Ar("abc", "ABC"), "xyz"},
		{"a x ; b x ; | c x |& d x && e x || f x |& g x >A <B >&A <&B >>>A", util.Ar("a", "b", "c", "d", "e", "f", "g"), ""},
		{"a x;b x|c x|&d x&&e x||f x|&g x>A<B>&A<&B>>>A", util.Ar("a", "b", "c", "d", "e", "f", "g"), ""},
	}
	for _, v := range tests {
		res := ParseCommandLine(&config, v.input)
		if !util.SlicesEqual(res.ExeNames, v.exes) {
			t.Errorf("input=%s expected=%v actual=%v", v.input, v.exes, res.ExeNames)
		}
		if v.comment != res.Comment {
			t.Errorf("input=%s expected=%s actual=%s", v.input, v.comment, res.Comment)
		}
	}
}
