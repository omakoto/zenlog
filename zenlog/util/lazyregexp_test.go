package util

import "testing"

func TestLazyRegexp_Pattern(t *testing.T) {
	p := NewLazyRegexp("")
	if p.regexp != nil {
		t.Errorf("p.regexp != nil")
	}
	p.Pattern()
	if p.regexp == nil {
		t.Errorf("p.regexp == nil")
	}
	first := p.regexp

	p.Pattern()

	if first != p.regexp {
		t.Errorf("p.regexp changed.")
	}
}
