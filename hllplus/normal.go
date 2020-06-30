package hllplus

import "math/bits"

func normalDowngrade(pos int, rhoW, sourceP, targetP uint8) uint8 {
	// Preserve 0 rhoW in the normal encoding since this represents any unset register.
	if rhoW == 0 {
		return 0
	}

	// Splice off the new pos by bit shifting just past the pos prefix. If the new suffix is
	// not all zeros, then the new rhoW is just the number of leading zeros + 1 in the new suffix.
	//
	// Otherwise, the old rhoW needs to be updated to account for the additional number of leading
	// zeros.
	w := uint64(pos) << uint8(64-sourceP+targetP)
	if w == 0 {
		return rhoW + sourceP - targetP
	}
	return uint8(bits.LeadingZeros64(w)) + 1
}
