package main

import "testing"

func TestNewInterval(t *testing.T) {
	const m = 1
	const h = 60 * m
	const d = 24 * h
	cases := []struct {
		CurrentInterval  int64
		Ease             int64
		IntervalModifier int64
		NewInterval      int64
	}{
		{10 * d, 250, 100, 25 * d},
		{25 * d, 250, 100, 62 * d},
		{62 * d, 250, 100, 155},
	}

	for _, c := range cases {
		newInterval := GetNewInterval(c.CurrentInterval, c.Ease, c.IntervalModifier)
		if newInterval != c.NewInterval {
			t.Fatalf("for %#v: expect(%d) != %d", c, newInterval, c.NewInterval)
		}

	}
}
