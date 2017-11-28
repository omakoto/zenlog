package util

import (
	"testing"
)

func TestEncodeSuccess(t *testing.T) {
	tests := []struct {
		prefix string
		args   []string
	}{
		{"", Ar()},
		{"abc", Ar()},
		{"", Ar("a")},
		{"", Ar("a", "b")},
		{"", Ar("a b", "b")},
		{"日本", Ar("漢 字", "によ")},
	}
	for i, v := range tests {
		s, p, a := TryDecodeBytes([]byte(v.prefix + Encode(v.args)))
		if s && (string(p) == v.prefix) && SlicesEqual(a, v.args) {
			// OK
		} else {
			t.Errorf("i=%d: [success=%s] Source pre='%s' args=[%d]%v, Actual pre='%s' args=[%d]%v", i, s, v.prefix, len(v.args), v.args, p, len(a), a)
		}
	}
}

func TestEncodeFailure(t *testing.T) {
	tests := []struct {
		value string
	}{
		{""},
		{"aa bb cc"},
		{"日  本"},
	}
	for _, v := range tests {
		s, _, _ := TryDecodeBytes([]byte(v.value))
		if s {
			t.Errorf("Success should be false for \"%s\".", v.value)
		}
	}
}
