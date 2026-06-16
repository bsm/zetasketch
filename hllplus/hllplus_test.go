package hllplus_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/bsm/zetasketch/hllplus"
)

func TestHLL_estimateNormal800(t *testing.T) {
	cases := []struct {
		p   int
		exp int64
	}{
		{10, 800},
		{11, 794},
		{12, 788},
		{13, 793},
		{14, 791},
		{15, 793},
		{16, 795},
		{17, 797},
		{18, 799},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("p=%d", tc.p), func(t *testing.T) {
			rnd := rand.New(rand.NewSource(33))
			subject, _ := hllplus.NewNormal(uint8(tc.p))
			for range 800 {
				subject.Add(rnd.Uint64())
			}
			if subject.IsSparse() {
				t.Error("expected normal representation")
			}
			if got := subject.Estimate(); got != tc.exp {
				t.Errorf("got %d, want %d", got, tc.exp)
			}
		})
	}
}

func TestHLL_estimateSparse800(t *testing.T) {
	cases := []struct {
		p   int
		exp int64
	}{
		{16, 795},
		{17, 798},
		{18, 799},
		{19, 799},
		{20, 799},
		{21, 799},
		{22, 800},
		{23, 800},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("p=%d", tc.p), func(t *testing.T) {
			rnd := rand.New(rand.NewSource(33))
			subject, _ := hllplus.New(uint8(tc.p-5), uint8(tc.p))
			for range 800 {
				subject.Add(rnd.Uint64())
			}
			if !subject.IsSparse() {
				t.Error("expected sparse representation")
			}
			if got := subject.Estimate(); got != tc.exp {
				t.Errorf("got %d, want %d", got, tc.exp)
			}
		})
	}
}

func TestHLL_estimateNormal200k(t *testing.T) {
	cases := []struct {
		p   int
		exp int64
	}{
		{10, 199197},
		{11, 204855},
		{12, 204356},
		{13, 202958},
		{14, 202977},
		{15, 201496},
		{16, 201398},
		{17, 201208},
		{18, 200999},
		{19, 200567},
		{20, 200032},
		{21, 200013},
		{22, 200003},
		{23, 199989},
		{24, 200026},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("p=%d", tc.p), func(t *testing.T) {
			rnd := rand.New(rand.NewSource(33))
			subject, _ := hllplus.NewNormal(uint8(tc.p))
			for range 200_000 {
				subject.Add(rnd.Uint64())
			}
			if subject.IsSparse() {
				t.Error("expected normal representation")
			}
			if got := subject.Estimate(); got != tc.exp {
				t.Errorf("got %d, want %d", got, tc.exp)
			}
		})
	}
}

func TestHLL_estimateSparse200k(t *testing.T) {
	cases := []struct {
		p   int
		exp int64
	}{
		{24, 200039},
		{25, 200048},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("p=%d", tc.p), func(t *testing.T) {
			rnd := rand.New(rand.NewSource(33))
			subject, _ := hllplus.New(uint8(tc.p-5), uint8(tc.p))
			for range 200_000 {
				subject.Add(rnd.Uint64())
			}
			if !subject.IsSparse() {
				t.Error("expected sparse representation")
			}
			if got := subject.Estimate(); got != tc.exp {
				t.Errorf("got %d, want %d", got, tc.exp)
			}
		})
	}
}

func TestHLL_estimateNormal100k2x50k(t *testing.T) {
	cases := []struct {
		p   int
		exp int64
	}{
		{10, 148131},
		{11, 152153},
		{12, 152338},
		{13, 152453},
		{14, 150853},
		{15, 150458},
		{16, 150811},
		{17, 150795},
		{18, 150592},
		{19, 150265},
		{20, 149879},
		{21, 149939},
		{22, 150005},
		{23, 149946},
		{24, 149988},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("p=%d", tc.p), func(t *testing.T) {
			rnd := rand.New(rand.NewSource(33))
			subject, _ := hllplus.NewNormal(uint8(tc.p))
			for range 100_000 {
				subject.Add(rnd.Uint64())
			}
			for range 50_000 {
				h := rnd.Uint64()
				subject.Add(h)
				subject.Add(h)
			}
			if subject.IsSparse() {
				t.Error("expected normal representation")
			}
			if got := subject.Estimate(); got != tc.exp {
				t.Errorf("got %d, want %d", got, tc.exp)
			}
		})
	}
}

func TestHLL_estimateSparse100k2x50k(t *testing.T) {
	cases := []struct {
		p   int
		exp int64
	}{
		{23, 149969},
		{24, 149998},
		{25, 150012},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("p=%d", tc.p), func(t *testing.T) {
			rnd := rand.New(rand.NewSource(33))
			subject, _ := hllplus.New(uint8(tc.p-5), uint8(tc.p))
			for range 100_000 {
				subject.Add(rnd.Uint64())
			}
			for range 50_000 {
				h := rnd.Uint64()
				subject.Add(h)
				subject.Add(h)
			}
			if !subject.IsSparse() {
				t.Error("expected sparse representation")
			}
			if got := subject.Estimate(); got != tc.exp {
				t.Errorf("got %d, want %d", got, tc.exp)
			}
		})
	}
}

func TestHLL_estimateSparseRepetitive(t *testing.T) {
	subject, _ := hllplus.New(12, 17)

	// same element added again and again - Add(N), Estimate() -> Flush(), repeat:
	for range 100 {
		subject.Add(1)
		if got := subject.Estimate(); got != 1 {
			t.Fatalf("got %d, want 1", got)
		}
	}

	// sanity check that HLL has not been normalized:
	if !subject.IsSparse() {
		t.Error("expected sparse representation")
	}
}

func TestHLL_normalize(t *testing.T) {
	rnd := rand.New(rand.NewSource(33))
	subject, _ := hllplus.New(12, 17)
	for range 3_084 {
		subject.Add(rnd.Uint64())
	}
	if !subject.IsSparse() {
		t.Error("expected sparse representation")
	}
	if got := subject.Estimate(); got != 3_083 {
		t.Errorf("got %d, want 3083", got)
	}

	subject.Add(rnd.Uint64())
	if subject.IsSparse() {
		t.Error("expected normal representation")
	}
	if got := subject.Estimate(); got != 3_072 {
		t.Errorf("got %d, want 3072", got)
	}
}

func TestHLL_downgrade(t *testing.T) {
	rnd := rand.New(rand.NewSource(33))
	s1, _ := hllplus.NewNormal(15)
	s2, _ := hllplus.NewNormal(12)

	for range 100_000 {
		n := rnd.Uint64()
		s1.Add(n)
		s2.Add(n)
	}

	if got := s1.Estimate(); got != 99879 {
		t.Errorf("s1.Estimate: got %d, want 99879", got)
	}
	if got := s2.Estimate(); got != 100680 {
		t.Errorf("s2.Estimate: got %d, want 100680", got)
	}

	if got := s1.Precision(); got != 15 {
		t.Errorf("s1.Precision: got %d, want 15", got)
	}
	if got := s1.SparsePrecision(); got != 20 {
		t.Errorf("s1.SparsePrecision: got %d, want 20", got)
	}
	if err := s1.Downgrade(12, 17); err != nil {
		t.Fatal(err)
	}
	if got := s1.Precision(); got != 12 {
		t.Errorf("s1.Precision: got %d, want 12", got)
	}
	if got := s1.SparsePrecision(); got != 17 {
		t.Errorf("s1.SparsePrecision: got %d, want 17", got)
	}

	if got := s1.Estimate(); got != 100680 {
		t.Errorf("s1.Estimate: got %d, want 100680", got)
	}
	if got := s2.Estimate(); got != 100680 {
		t.Errorf("s2.Estimate: got %d, want 100680", got)
	}
}

// newMergeFixture builds three sketches sharing 50k values and adds 50k distinct
// values to each, matching the original merge spec's BeforeEach.
func newMergeFixture(t *testing.T) (s1, s2, s3 *hllplus.HLL) {
	t.Helper()

	rnd := rand.New(rand.NewSource(33))
	s1, _ = hllplus.NewNormal(15)
	s2, _ = hllplus.NewNormal(15)
	s3, _ = hllplus.NewNormal(12)

	for range 50_000 {
		n := rnd.Uint64()
		s1.Add(n)
		s2.Add(n)
		s3.Add(n)
	}
	for range 50_000 {
		s1.Add(rnd.Uint64())
		s2.Add(rnd.Uint64())
		s3.Add(rnd.Uint64())
	}

	if got := s1.Estimate(); got != 100324 {
		t.Fatalf("s1.Estimate: got %d, want 100324", got)
	}
	if got := s2.Estimate(); got != 100168 {
		t.Fatalf("s2.Estimate: got %d, want 100168", got)
	}
	if got := s3.Estimate(); got != 100464 {
		t.Fatalf("s3.Estimate: got %d, want 100464", got)
	}
	return s1, s2, s3
}

func TestHLL_merge_equalPrecision(t *testing.T) {
	s1, s2, _ := newMergeFixture(t)
	s1.Merge(s2)
	if got := s1.Estimate(); got != 150794 {
		t.Errorf("got %d, want 150794", got)
	}
}

func TestHLL_merge_lowerPrecisionTarget(t *testing.T) {
	s1, _, s3 := newMergeFixture(t)
	s1.Merge(s3)
	if got := s1.Estimate(); got != 154744 {
		t.Errorf("got %d, want 154744", got)
	}
	if got := s1.Precision(); got != 12 {
		t.Errorf("Precision: got %d, want 12", got)
	}
	if got := s1.SparsePrecision(); got != 17 {
		t.Errorf("SparsePrecision: got %d, want 17", got)
	}
}

func TestHLL_merge_higherPrecisionTarget(t *testing.T) {
	s1, _, s3 := newMergeFixture(t)
	s3.Merge(s1)
	if got := s3.Estimate(); got != 154744 {
		t.Errorf("got %d, want 154744", got)
	}
	if got := s3.Precision(); got != 12 {
		t.Errorf("Precision: got %d, want 12", got)
	}
	if got := s3.SparsePrecision(); got != 17 {
		t.Errorf("SparsePrecision: got %d, want 17", got)
	}
}

func TestHLL_merge_emptyTarget(t *testing.T) {
	s1, _, _ := newMergeFixture(t)

	subject, _ := hllplus.NewNormal(15)
	subject.Merge(s1)

	// just a straight copy of s1:
	if got, exp := subject.Estimate(), s1.Estimate(); got != exp {
		t.Errorf("Estimate: got %d, want %d", got, exp)
	}
	if got, exp := subject.Precision(), s1.Precision(); got != exp {
		t.Errorf("Precision: got %d, want %d", got, exp)
	}
	if got, exp := subject.SparsePrecision(), s1.SparsePrecision(); got != exp {
		t.Errorf("SparsePrecision: got %d, want %d", got, exp)
	}
}

func TestHLL_proto_initNormal(t *testing.T) {
	rnd := rand.New(rand.NewSource(33))
	subject, _ := hllplus.New(12, 17)
	for range 10_000 {
		subject.Add(rnd.Uint64())
	}
	if subject.IsSparse() {
		t.Error("expected normal representation")
	}
	if got := subject.Estimate(); got != 9_912 {
		t.Errorf("got %d, want 9912", got)
	}

	msg := subject.Proto()

	// both precisions are always stored:
	if got := msg.GetPrecisionOrNumBuckets(); got != 12 {
		t.Errorf("precision: got %d, want 12", got)
	}
	if got := msg.GetSparsePrecisionOrNumBuckets(); got != 17 {
		t.Errorf("sparse precision: got %d, want 17", got)
	}

	// expect normal representation:
	if len(msg.GetData()) == 0 {
		t.Error("expected non-empty data")
	}

	// expect NO sparse representation:
	if msg.SparseSize != nil {
		t.Error("expected nil sparse size")
	}
	if msg.GetSparseData() != nil {
		t.Error("expected nil sparse data")
	}

	// init back from proto:
	restored, err := hllplus.NewFromProto(msg)
	if err != nil {
		t.Fatal(err)
	}
	if restored.IsSparse() {
		t.Error("expected normal representation")
	}
	if got := restored.Precision(); got != 12 {
		t.Errorf("precision: got %d, want 12", got)
	}
	if got := restored.SparsePrecision(); got != 17 {
		t.Errorf("sparse precision: got %d, want 17", got)
	}
	if got := restored.Estimate(); got != 9_912 {
		t.Errorf("got %d, want 9912", got)
	}
}

func TestHLL_proto_initSparse(t *testing.T) {
	rnd := rand.New(rand.NewSource(33))
	subject, _ := hllplus.New(12, 17)
	for range 800 {
		subject.Add(rnd.Uint64())
	}
	if !subject.IsSparse() {
		t.Error("expected sparse representation")
	}
	if got := subject.Estimate(); got != 798 {
		t.Errorf("got %d, want 798", got)
	}

	msg := subject.Proto()

	// both precisions are always stored:
	if got := msg.GetPrecisionOrNumBuckets(); got != 12 {
		t.Errorf("precision: got %d, want 12", got)
	}
	if got := msg.GetSparsePrecisionOrNumBuckets(); got != 17 {
		t.Errorf("sparse precision: got %d, want 17", got)
	}

	// expect NO normal representation:
	if len(msg.GetData()) != 0 {
		t.Error("expected empty data")
	}

	// expect sparse representation:
	// hash/rand collisions are fine, that's why it is != 800
	if got := msg.GetSparseSize(); got != 796 {
		t.Errorf("sparse size: got %d, want 796", got)
	}
	if len(msg.GetSparseData()) == 0 {
		t.Error("expected non-empty sparse data")
	}

	// init back from proto:
	restored, err := hllplus.NewFromProto(msg)
	if err != nil {
		t.Fatal(err)
	}
	if !restored.IsSparse() {
		t.Error("expected sparse representation")
	}
	if got := restored.Precision(); got != 12 {
		t.Errorf("precision: got %d, want 12", got)
	}
	if got := restored.SparsePrecision(); got != 17 {
		t.Errorf("sparse precision: got %d, want 17", got)
	}
	if got := restored.Estimate(); got != 798 {
		t.Errorf("got %d, want 798", got)
	}
}
