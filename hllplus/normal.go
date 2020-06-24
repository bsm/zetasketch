package hllplus

import "math/bits"

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

func downgradeRhoW(index int, rhoW, sourceP, targetP uint8) uint8 {
	// Preserve 0 rhoW in the normal encoding since this represents any unset register.
	if rhoW == 0 {
		return 0
	}

	// Splice off the new index by bit shifting just past the index prefix. If the new suffix is
	// not all zeros, then the new rhoW is just the number of leading zeros + 1 in the new suffix.
	//
	// Otherwise, the old rhoW needs to be updated to account for the additional number of leading
	// zeros.
	w := uint64(index) << uint8(64-sourceP+targetP)
	if w == 0 {
		return rhoW + sourceP - targetP
	}
	return uint8(bits.LeadingZeros64(w)) + 1
}
