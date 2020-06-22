// Package zetasketch provides simplified Go implementation of https://github.com/google/zetasketch ,
// compatible with BigQuery HyperLogLog++ https://cloud.google.com/bigquery/docs/reference/standard-sql/hll_functions
package zetasketch

import (
	"encoding/base64"
	"errors"

	"github.com/bsm/zetasketch/internal/fingerprint"
	"github.com/bsm/zetasketch/internal/hyperloglog"
	"github.com/bsm/zetasketch/internal/zetasketch"
	"google.golang.org/protobuf/proto"
)

// HLL.Marshal global vars (Go cannot take pointers of consts).
var (
	hllEncodingVersion int32 = 2
	hllAggregatorType        = zetasketch.AggregatorType_HYPERLOGLOG_PLUS_UNIQUE
)

// HLL is a HyperLogLog++ aggregator.
// It is NOT thread-safe.
type HLL struct {
	hll       *hyperloglog.HyperLogLogPlus
	precision int32
	numValues int64
}

// NewHLL initializes a new HyperLogLog++ aggregator.
func NewHLL(precision uint8) (*HLL, error) {
	hll, err := hyperloglog.NewPlus(precision)
	if err != nil {
		return nil, err
	}
	return &HLL{
		hll:       hll,
		precision: int32(precision),
	}, nil
}

// Add adds a byte-slice value.
//
// WARNING: hashing of data more than 32 bytes long is not implemented yet.
func (a *HLL) Add(value []byte) {
	a.numValues++
	a.hll.Add(hasher(value))
}

// Count returns current estimated count.
func (a *HLL) Count() uint64 {
	return a.hll.Count()
}

// Merge merges another HyperLogLog++ aggregator into current one.
func (a *HLL) Merge(another *HLL) error {
	return a.hll.Merge(another.hll)
}

// Clear clears internal state of an aggregator so it can be reused.
func (a *HLL) Clear() {
	a.hll.Clear()
}

// Marshal marshals an aggregator to a binary proto message (raw bytes, not base64 encoded).
func (a *HLL) Marshal() ([]byte, error) {
	a.hll.ToNormal() // force normal (dense) representation

	aggState := &zetasketch.AggregatorStateProto{
		Type:            &hllAggregatorType,
		NumValues:       &a.numValues,        // TODO: we must track each addition in a wrapper around HLL++
		EncodingVersion: &hllEncodingVersion, // fixed
		ValueType:       nil,                 // looks to be a type of values being added - strings, bytes, ints etc - I think, can be omitted (may need to check though)
	}
	hllState := &zetasketch.HyperLogLogPlusUniqueStateProto{
		SparseSize:                  nil, // we do not support sparse (yet)
		PrecisionOrNumBuckets:       &a.precision,
		SparsePrecisionOrNumBuckets: nil, // we do not support sparse (yet)
		Data:                        a.hll.Reg,
		SparseData:                  nil, // we do not support sparse (yet)
	}
	proto.SetExtension(aggState, zetasketch.E_HyperloglogplusUniqueState, hllState)

	return proto.Marshal(aggState)
}

// MarshalJSON serializes aggregator to JSON.
func (a *HLL) MarshalJSON() ([]byte, error) {
	data, err := a.Marshal()
	if err != nil {
		return nil, err
	}

	enc := base64.StdEncoding
	n := enc.EncodedLen(len(data))
	buf := make([]byte, n+2) // base64 data + two double quotes
	buf[0], buf[n+1] = '"', '"'
	enc.Encode(buf[1:], data)
	return buf, nil
}

// UnmarshalJSON is a dummy method - unmarshaling from JSON is not supported.
func (a *HLL) UnmarshalJSON([]byte) error {
	return errors.New("marshalling HLL aggregator from JSON is not supported")
}

// ----------------------------------------------------------------------------

type hasher []byte

func (h hasher) Sum64() uint64 {
	return fingerprint.Hash64(h)
}
