package hllplus

import "math/bits"

const (
	minDensePrecision = 10
	maxDensePrecision = 24
)

// Returns the HyperLogLog++ index for the given hash.
func calcIndex(hash uint64, precision uint8) int {
	idx := hash >> (64 - precision)
	return int(idx)
}

// Returns the HyperLogLog++ ρ(w) for the given hash, which is the number of
// leading zero bits + 1 for the bits after the normal index.
func calcRhoW(hash uint64, precision uint8) uint8 {
	return computeRhoW(hash, 64-precision)
}

func computeRhoW(value uint64, offset uint8) uint8 {
	// Strip of the index and move the rhoW to a higher order.
	w := value << (64 - offset)
	if w == 0 {
		return offset + 1
	}
	return uint8(bits.LeadingZeros64(w)) + 1
}
