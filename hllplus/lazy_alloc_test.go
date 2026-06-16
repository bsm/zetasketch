package hllplus_test

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/bsm/zetasketch/hllplus"

	"google.golang.org/protobuf/proto"
)

// buildHLL constructs an HLL sketch with the given precisions and feeds it n
// values drawn deterministically from seed, so that two calls with the same
// arguments observe byte-for-byte identical input.
func buildHLL(tb testing.TB, precision, sparsePrecision uint8, n int, seed int64) *hllplus.HLL {
	tb.Helper()

	h, err := hllplus.New(precision, sparsePrecision)
	if err != nil {
		tb.Fatalf("New(%d, %d): %v", precision, sparsePrecision, err)
	}

	rnd := rand.New(rand.NewSource(seed))
	for range n {
		h.Add(rnd.Uint64())
	}
	return h
}

func marshalHLL(tb testing.TB, h *hllplus.HLL) []byte {
	tb.Helper()

	b, err := proto.Marshal(h.Proto())
	if err != nil {
		tb.Fatalf("marshal: %v", err)
	}
	return b
}

// TestLazyAllocInvariant guards that lazily allocating the sparse delta slice and
// buffer map does not change observable behavior. Across a wide low/high
// cardinality range and a couple of precisions it asserts that the cardinality
// estimate and the serialized bytes are stable and fully deterministic: two
// independently built sketches over identical input must agree byte-for-byte, and
// a serialize/deserialize round-trip must preserve the estimate. The sparse
// buffer is a map, so this also pins down that serialization stays independent of
// map iteration order.
func TestLazyAllocInvariant(t *testing.T) {
	precisions := []struct{ p, sp uint8 }{
		{12, 17},
		{15, 20},
	}
	cardinalities := []int{1, 10, 100, 10_000, 1_000_000}

	const seed = 42

	for _, pr := range precisions {
		for _, n := range cardinalities {
			if n >= 1_000_000 && testing.Short() {
				continue
			}

			a := buildHLL(t, pr.p, pr.sp, n, seed)
			b := buildHLL(t, pr.p, pr.sp, n, seed)

			// (b) cardinality estimate is identical across independent builds.
			estA, estB := a.Estimate(), b.Estimate()
			if estA != estB {
				t.Fatalf("p=%d sp=%d n=%d: estimates differ: %d != %d", pr.p, pr.sp, n, estA, estB)
			}

			// (a) serialized bytes are byte-identical across independent builds.
			bytesA, bytesB := marshalHLL(t, a), marshalHLL(t, b)
			if !bytes.Equal(bytesA, bytesB) {
				t.Fatalf("p=%d sp=%d n=%d: serialized bytes differ (%d vs %d bytes)",
					pr.p, pr.sp, n, len(bytesA), len(bytesB))
			}

			// round-trip: serialization faithfully restores the estimate.
			restored, err := hllplus.NewFromProto(a.Proto())
			if err != nil {
				t.Fatalf("p=%d sp=%d n=%d: NewFromProto: %v", pr.p, pr.sp, n, err)
			}
			if got := restored.Estimate(); got != estA {
				t.Fatalf("p=%d sp=%d n=%d: round-trip estimate changed: %d != %d", pr.p, pr.sp, n, got, estA)
			}

			// sanity: the estimate stays within HLL++ error bounds of the truth.
			if d := estA - int64(n); d < -(int64(n)/20+2) || d > int64(n)/20+2 {
				t.Fatalf("p=%d sp=%d n=%d: estimate %d outside tolerance", pr.p, pr.sp, n, estA)
			}
		}
	}
}

var benchSink *hllplus.HLL

// BenchmarkSparseLowCardinality measures the per-sketch heap cost of the common
// low-cardinality case: construct a fresh sparse sketch and add a handful of
// distinct values. Run with -benchmem to read B/op and allocs/op.
func BenchmarkSparseLowCardinality(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		h, err := hllplus.New(12, 17)
		if err != nil {
			b.Fatal(err)
		}
		for j := range 10 {
			h.Add(uint64(j) * 0x9e3779b97f4a7c15)
		}
		benchSink = h
	}
}
