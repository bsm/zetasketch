// Package zetasketch provices a collection of libraries for single-pass, distributed,
// sublinear-space approximate aggregation and sketching algorithms.
package zetasketch

import (
	"encoding"

	"github.com/bsm/zetasketch/internal/hash"
)

// Aggregator provides an interface that wraps
// distributed, online aggregation algorithm.
type Aggregator interface {
	// Add adds a value.
	Add(v Value)
	// NumValues returns the total number of input values that this aggregator has seen.
	NumValues() uint64
	// Merge merges two aggregators.
	Merge(other Aggregator) error

	encoding.BinaryMarshaler
}

// Value is a hashable value.
type Value interface {
	Sum64() uint64
}

type hashSum uint64

func (v hashSum) Sum64() uint64 { return uint64(v) }

// StringValue converts a string to a Value.
func StringValue(s string) Value {
	return BinaryValue([]byte(s))
}

// BinaryValue converts a byte slice to a Value.
func BinaryValue(p []byte) Value {
	return hashSum(hash.Bytes(p))
}

// Uint32Value converts a number to a Value.
func Uint32Value(v uint32) Value {
	return hashSum(hash.Uint32(v))
}

// Uint64Value converts a number slice to a Value.
func Uint64Value(v uint64) Value {
	return hashSum(hash.Uint64(v))
}
