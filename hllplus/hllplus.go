package hllplus

import (
	"fmt"
	"math"
)

// HLL is a HyperLogLog++ sketch implementation.
type HLL struct {
	dense []byte

	p  uint8 // dense precision
	pp uint8 // sparse precision
}

// New inits a new sketch.
// The dense precision p must be between 10 and 24.
// The sparse precision pp must be between 0 and 25.
// This function only returns an error when an invalid precision is provided.
func New(p, pp uint8) (*HLL, error) {
	if p < minDensePrecision || p > maxDensePrecision {
		return nil, fmt.Errorf("invalid dense precision %d", p)
	}
	if pp > maxSparsePrecision {
		return nil, fmt.Errorf("invalid sparse precision %d", pp)
	}

	return &HLL{p: p, pp: pp}, nil
}

// Add adds the uniform hash value to the representation.
func (s *HLL) Add(hash uint64) {
	if len(s.dense) == 0 {
		s.dense = make([]byte, 1<<s.p)
	}

	pos := calcIndex(hash, s.p)
	rho := calcRhoW(hash, s.p)
	if rho > s.dense[pos] {
		s.dense[pos] = rho
	}
}

// Estimate computes the cardinality estimate according to the algorithm in Figure 6 of the HLL++ paper
// (https://goo.gl/pc916Z).
func (s *HLL) Estimate() uint64 {
	if len(s.dense) == 0 {
		return 0
	}

	// Compute the summation component of the harmonic mean for the HLL++ algorithm while also
	// keeping track of the number of zeros in case we need to apply LinearCounting instead.
	numZeros := 0
	sum := 0.0

	for _, c := range s.dense {
		if c == 0 {
			numZeros++
		}

		// Compute sum += math.pow(2, -v) without actually performing a floating point exponent
		// computation (which is expensive). v can be at most 64 - precision + 1 and the minimum
		// precision is larger than 2 (see MINIMUM_PRECISION), so this left shift can not overflow.
		x := 1 << c
		sum += 1.0 / float64(x)
	}

	// Return the LinearCount for small cardinalities where, as explained in the HLL++ paper
	// (https://goo.gl/pc916Z), the results with LinearCount tend to be more accurate than with HLL.
	x := 1 << s.p
	m := float64(x)
	if numZeros != 0 {
		h := m * math.Log(m/float64(numZeros))
		if int(h) <= linearCountingThreshold(s.p) {
			return uint64(math.Round(h))
		}
	}

	// The "raw" estimate, designated by E in the HLL++ paper (https://goo.gl/pc916Z).
	raw := alpha(s.p) * m * m / sum

	// Perform bias correction on small estimates. HyperLogLogPlusPlusData only contains bias
	// estimates for small cardinalities and returns 0 for anything else, so the "E < 5m" guard from
	// the HLL++ paper (https://goo.gl/pc916Z) is superfluous here.
	return uint64(math.Round(raw - estimateBias(raw, s.p)))
}
