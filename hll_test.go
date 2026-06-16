package zetasketch_test

import (
	"testing"

	"github.com/bsm/zetasketch"
)

var _ zetasketch.Aggregator = (*zetasketch.HLL)(nil)

func newTestHLL() *zetasketch.HLL {
	subject := zetasketch.NewHLL(nil)
	for i := range 1_000 {
		subject.Add(zetasketch.Uint64Value(uint64(i)))
	}
	for i := 500; i < 1_000; i++ {
		subject.Add(zetasketch.Uint64Value(uint64(i)))
	}
	return subject
}

func TestHLL_NumValues(t *testing.T) {
	subject := newTestHLL()
	if got, exp := subject.NumValues(), int64(1_500); got != exp {
		t.Errorf("got %d, want %d", got, exp)
	}
}

func TestHLL_Result(t *testing.T) {
	subject := newTestHLL()
	if got, exp := subject.Result(), int64(1_000); got != exp {
		t.Errorf("got %d, want %d", got, exp)
	}
}

func TestHLL_Merge(t *testing.T) {
	subject := newTestHLL()

	other := zetasketch.NewHLL(nil)
	for i := 800; i < 1_200; i++ {
		other.Add(zetasketch.Uint64Value(uint64(i)))
	}

	if err := subject.Merge(other); err != nil {
		t.Fatal(err)
	}
	if got, exp := subject.NumValues(), int64(1_900); got != exp {
		t.Errorf("NumValues: got %d, want %d", got, exp)
	}
	if got, exp := subject.Result(), int64(1_207); got != exp {
		t.Errorf("Result: got %d, want %d", got, exp)
	}

	// `other` is not modified:
	if got, exp := other.NumValues(), int64(400); got != exp {
		t.Errorf("other.NumValues: got %d, want %d", got, exp)
	}
	if got, exp := other.Result(), int64(400); got != exp {
		t.Errorf("other.Result: got %d, want %d", got, exp)
	}
}

func TestHLL_MarshalBinary(t *testing.T) {
	subject := newTestHLL()

	data, err := subject.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	subject = new(zetasketch.HLL)
	if err := subject.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
	if got, exp := subject.NumValues(), int64(1_500); got != exp {
		t.Errorf("NumValues: got %d, want %d", got, exp)
	}
	if got, exp := subject.Result(), int64(1_000); got != exp {
		t.Errorf("Result: got %d, want %d", got, exp)
	}
}
