package hllplus

// EstimateBias test export.
func EstimateBias(e float64, p uint8) float64 {
	return estimateBias(e, p)
}

func NewNormal(precision uint8) (*HLL, error) {
	pp := min(precision+5, MaxSparsePrecision)

	s, err := New(precision, pp)
	if err != nil {
		return nil, err
	}
	s.normalize()
	return s, nil
}

func (s *HLL) IsSparse() bool {
	return s.sparse != nil
}
