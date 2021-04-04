package main

import "testing"

func TestAlignmentFormat(t *testing.T) {
	var tests = []struct {
		alignment
		expected string
	}{
		{alignment{}, "True Neutral"},
		{alignment{-1., 0}, "Chaotic Neutral"},
		{alignment{1., 0}, "Lawful Neutral"},
		{alignment{0, 1.}, "Neutral Good"},
		{alignment{0, -1.}, "Neutral Evil"},
		{alignment{-1., 1.}, "Chaotic Good"},
		{alignment{1., 1.}, "Lawful Good"},
		{alignment{-1., -1.}, "Chaotic Evil"},
	}
	for i := range tests {
		if tests[i].Format() != tests[i].expected {
			t.Errorf("expected %q, got %q for %#v", tests[i].expected, tests[i].Format(), tests[i].alignment)
		}
	}
}
