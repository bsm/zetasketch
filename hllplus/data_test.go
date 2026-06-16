package hllplus_test

import (
	"math"
	"testing"

	"github.com/bsm/zetasketch/hllplus"
)

func TestEstimateBias(t *testing.T) {
	cases := []struct {
		e     float64
		p     uint8
		exp   float64
		delta float64
	}{
		{0, 15, 0.0, 0},
		{1, 15, 0.0, 0},
		{10_000, 15, 0.0, 0},
		{100_000, 15, 888.1, 0.1},
		{200_000, 15, 0.0, 0},

		{50_000, 13, 0.0, 0},
		{50_000, 14, 449.7, 0.1},
		{50_000, 15, 7820.2, 0.1},
		{50_000, 16, 44513.2, 0.1},
		{50_000, 17, 0.0, 0},
	}
	for _, tc := range cases {
		got := hllplus.EstimateBias(tc.e, tc.p)
		if math.Abs(got-tc.exp) > tc.delta {
			t.Errorf("EstimateBias(%v, %d) = %v, want %v (±%v)", tc.e, tc.p, got, tc.exp, tc.delta)
		}
	}
}
