package hllplus

import (
	"encoding/binary"
	"math"
	"sort"
	"sync"
)

const sparseRhoWBits = 6

type sparseState struct {
	normalPrecision uint8
	sparsePrecision uint8

	data   *deltaSlice
	buffer uint32Set

	encodedFlag  uint32
	maxDataLen   int
	maxBufferLen int
}

func newSparseState(normalPrecision, sparsePrecision uint8) *sparseState {
	m := 1 << normalPrecision
	maxDataLen := m * 3 / 4
	maxBufferLen := m / 4

	encodedFlag := uint32(1 << sparsePrecision)
	if n := normalPrecision + sparseRhoWBits; n > sparsePrecision {
		encodedFlag = 1 << n
	}

	return &sparseState{
		normalPrecision: normalPrecision,
		sparsePrecision: sparsePrecision,

		data:   recycleDeltaSlice(maxDataLen),
		buffer: make(uint32Set, maxBufferLen),

		encodedFlag:  encodedFlag,
		maxDataLen:   maxDataLen,
		maxBufferLen: maxBufferLen,
	}
}

func (s *sparseState) Add(hash uint64) {
	val := s.encode(hash)
	if s.buffer.Add(val); s.buffer.Len() >= s.maxBufferLen {
		s.Flush()
	}
}

// Linear counting over the number of empty sparse buckets.
func (s *sparseState) Estimate() int64 {
	mm := 1 << s.sparsePrecision
	numBuckets := float64(mm)
	numZeros := numBuckets - float64(s.data.Count())
	return int64(numBuckets*math.Log(numBuckets/numZeros) + 0.5)
}

func (s *sparseState) Clone() *sparseState {
	return &sparseState{
		normalPrecision: s.normalPrecision,
		sparsePrecision: s.sparsePrecision,

		data:   s.data.Clone(),
		buffer: s.buffer.Clone(),

		encodedFlag:  s.encodedFlag,
		maxDataLen:   s.maxDataLen,
		maxBufferLen: s.maxBufferLen,
	}
}

func (s *sparseState) Flush() {
	if s.buffer.Len() == 0 {
		return
	}

	result := recycleDeltaSlice(s.data.Len())
	buffered := s.buffer.Flush()

	// merge existing data and buffered
	s.data.Iterate(func(x uint32) {
		for {
			if len(buffered) != 0 && buffered[0] <= x {
				result.Append(buffered[0])
				buffered = buffered[1:]
			} else {
				result.Append(x)
				break
			}
		}
	})

	// append remaining
	for _, x := range buffered {
		result.Append(x)
	}

	// replace data
	s.data.Release()
	s.data = result
}

func (s *sparseState) OverMax() bool {
	return s.data.Len() > s.maxDataLen
}

func (s *sparseState) encode(hash uint64) uint32 {
	sparsePos, rho := computePosRhoW(hash, s.sparsePrecision)
	delta := s.sparsePrecision - s.normalPrecision

	// Check if the normal rhoW can be re-constructed from the lowest sp-p bits of the sparse
	// index. In that case, we do not need to encode it explicitly.
	if mask := uint32(1<<delta) - 1; sparsePos&mask != 0 {
		return sparsePos
	}

	// Use the normal index instead of the sparse index since the lowest sp-p bits are all 0
	// anyway (see the mask above).
	normPos := sparsePos >> delta
	return s.encodedFlag | normPos<<sparseRhoWBits | uint32(rho)
}

// --------------------------------------------------------------------

type uint32Set map[uint32]struct{}

func (s uint32Set) Add(n uint32) {
	s[n] = struct{}{}
}

func (s uint32Set) Len() int {
	return len(s)
}

func (s uint32Set) Clone() uint32Set {
	if s == nil {
		return nil
	}

	t := make(uint32Set, len(s))
	for n := range s {
		t.Add(n)
	}
	return t
}

func (s uint32Set) Flush() []uint32 {
	nums := make(uint32Slice, 0, len(s))
	for n := range s {
		nums = append(nums, n)
		delete(s, n)
	}
	sort.Sort(nums)
	return nums
}

type uint32Slice []uint32

func (p uint32Slice) Len() int           { return len(p) }
func (p uint32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p uint32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// --------------------------------------------------------------------

// Varint encoded series of uint32s.
type uvarintSlice []byte

func (s uvarintSlice) Append(x uint32) uvarintSlice {
	for x >= 0x80 {
		s = append(s, byte(x)|0x80)
		x >>= 7
	}
	return append(s, byte(x))
}

func (s uvarintSlice) Iterate(fn func(uint32)) {
	t := s
	for {
		x, m := binary.Uvarint(t)
		if m < 1 {
			break
		}
		fn(uint32(x))
		t = t[m:]
	}
}

// --------------------------------------------------------------------

var deltaSlicePool sync.Pool

// Delta encoded slice of uint32s.
type deltaSlice struct {
	nums uvarintSlice
	last uint32
	size int
}

func recycleDeltaSlice(size int) *deltaSlice {
	if v := deltaSlicePool.Get(); v != nil {
		return v.(*deltaSlice)
	}
	return &deltaSlice{nums: make(uvarintSlice, 0, size)}
}

func (s *deltaSlice) Len() int {
	return len(s.nums)
}

func (s *deltaSlice) Count() int {
	return s.size
}

func (s *deltaSlice) Reset() {
	s.nums = s.nums[:0]
	s.last = 0
	s.size = 0
}

func (s *deltaSlice) Release() {
	s.Reset()
	deltaSlicePool.Put(s)
}

func (s *deltaSlice) Clone() *deltaSlice {
	if s == nil {
		return nil
	}

	t := &deltaSlice{
		nums: make(uvarintSlice, len(s.nums)),
		last: s.last,
		size: s.size,
	}
	copy(t.nums, s.nums)
	return t
}

func (s *deltaSlice) Append(x uint32) {
	s.nums = s.nums.Append(x - s.last)
	s.last = x
	s.size++
}

func (s *deltaSlice) Iterate(fn func(uint32)) {
	var last uint32
	s.nums.Iterate(func(u uint32) {
		x := u + last
		fn(x)
		last = x
	})
}
