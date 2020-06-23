package hllplus

import (
	"fmt"
	"math"
)

// HLL is a HyperLogLog++ sketch implementation.
type HLL struct {
	normal []byte

	precision       uint8
	sparsePrecision uint8
}

// New inits a new sketch.
// The normal precision must be between 10 and 24.
// The sparse precision must be between 0 and 25.
// This function only returns an error when an invalid precision is provided.
func New(precision, sparsePrecision uint8) (*HLL, error) {
	if err := validate(precision, sparsePrecision); err != nil {
		return nil, err
	}

	return &HLL{precision: precision, sparsePrecision: sparsePrecision}, nil
}

// Precision returns the normal precision.
func (s *HLL) Precision() uint8 {
	return s.precision
}

// SparsePrecision returns the sparse precision.
func (s *HLL) SparsePrecision() uint8 {
	return s.sparsePrecision
}

// Add adds the uniform hash value to the representation.
func (s *HLL) Add(hash uint64) {
	if len(s.normal) == 0 {
		s.normal = make([]byte, 1<<s.precision)
	}

	pos := calcIndex(hash, s.precision)
	rho := calcRhoW(hash, s.precision)
	if rho > s.normal[pos] {
		s.normal[pos] = rho
	}
}

// Merge merges other into s.
func (s *HLL) Merge(other *HLL) {
	// Skip if there is nothing to merge.
	if len(other.normal) == 0 {
		return
	}

	// If other precision is higher.
	if s.precision < other.precision {
		other.eachRhoWDowngrade(s.precision, func(index int, rhoW uint8) {
			if s.normal[index] < rhoW {
				s.normal[index] = rhoW
			}
		})
		return
	}

	// If other precision is lower, downgrade.
	if s.precision > other.precision {
		_ = s.Downgrade(other.precision, other.sparsePrecision)
	}

	// Use largest rhoW.
	for i, rho := range other.normal {
		if s.normal[i] < rho {
			s.normal[i] = rho
		}
	}
}

// Clone creates a copy of the sketch.
func (s *HLL) Clone() *HLL {
	clone := &HLL{
		precision:       s.precision,
		sparsePrecision: s.sparsePrecision,
	}
	if len(s.normal) != 0 {
		clone.normal = make([]byte, len(s.normal))
		copy(clone.normal, s.normal)
	}
	return clone
}

// Estimate computes the cardinality estimate according to the algorithm in Figure 6 of the HLL++ paper
// (https://goo.gl/pc916Z).
func (s *HLL) Estimate() uint64 {
	if len(s.normal) == 0 {
		return 0
	}

	// Compute the summation component of the harmonic mean for the HLL++ algorithm while also
	// keeping track of the number of zeros in case we need to apply LinearCounting instead.
	numZeros := 0
	sum := 0.0

	for _, c := range s.normal {
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
	x := 1 << s.precision
	m := float64(x)
	if numZeros != 0 {
		n := uint64(m*math.Log(m/float64(numZeros)) + 0.5)
		if n <= linearCountingThreshold(s.precision) {
			return n
		}
	}

	// The "raw" estimate, designated by E in the HLL++ paper (https://goo.gl/pc916Z).
	raw := alpha(s.precision) * m * m / sum

	// Perform bias correction on small estimates. HyperLogLogPlusPlusData only contains bias
	// estimates for small cardinalities and returns 0 for anything else, so the "E < 5m" guard from
	// the HLL++ paper (https://goo.gl/pc916Z) is superfluous here.
	return uint64(raw - estimateBias(raw, s.precision) + 0.5)
}

// Downgrade tries to reduce the precision of the sketch.
// Attempts to increase precision will be ignored.
func (s *HLL) Downgrade(precision, sparsePrecision uint8) error {
	if err := validate(precision, sparsePrecision); err != nil {
		return err
	}

	if s.precision > precision {
		if len(s.normal) != 0 {
			normal := make([]byte, 1<<precision)
			s.eachRhoWDowngrade(precision, func(index int, rhoW uint8) {
				if normal[index] < rhoW {
					normal[index] = rhoW
				}
			})
			s.normal = normal
		}
		s.precision = precision
	}

	if s.sparsePrecision > sparsePrecision {
		s.sparsePrecision = sparsePrecision
	}
	return nil
}

func (s *HLL) eachRhoWDowngrade(targetPrecision uint8, iter func(int, uint8)) {
	for idx, rho := range s.normal {
		newI := idx >> (s.precision - targetPrecision)
		newR := downgradeRhoW(idx, rho, s.precision, targetPrecision)
		iter(newI, newR)
	}
}

func validate(precision, sparsePrecision uint8) error {
	if precision < minNormalPrecision || precision > maxNormalPrecision {
		return fmt.Errorf("invalid normal precision %d", precision)
	}
	if sparsePrecision > maxSparsePrecision {
		return fmt.Errorf("invalid sparse precision %d", sparsePrecision)
	}
	return nil
}

// GetData exposes underlying binary sketch.
//
// `sparseSize` is returned as `-1` if internal representation is dense.
func (s *HLL) GetData() (data []byte, sparseSize int32) {
	// TODO: alter this when sparse representation is implemented
	return s.dense, -1
}
