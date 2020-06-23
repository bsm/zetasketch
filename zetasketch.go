// Package zetasketch provides simplified Go implementation of https://github.com/google/zetasketch ,
// compatible with BigQuery HyperLogLog++ https://cloud.google.com/bigquery/docs/reference/standard-sql/hll_functions
package zetasketch

import (
	"github.com/bsm/zetasketch/hllplus"
	"github.com/bsm/zetasketch/internal/fingerprint"
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
	hll *hllplus.HLL // kept private not to forget to implement proxy methods (to bump numValues on any AddTYPE etc)

	// these values needed for marshalling:
	precision, sparsePrecision int32
	numValues                  int64
}

// NewHLL initializes a new HyperLogLog++ aggregator.
func NewHLL(precision /* 10..24 */, sparsePrecision /* 0..25 */ uint8) *HLL {
	hll, err := hllplus.New(precision, sparsePrecision)
	if err != nil {
		panic("hllplus init failed: " + err.Error()) // occurs only for bad precisions
	}
	return &HLL{
		hll:             hll,
		precision:       int32(precision),
		sparsePrecision: int32(sparsePrecision),
	}
}

// AddBytes adds a byte-slice value.
//
// WARNING: hashing of data more than 32 bytes long is not implemented yet, it will panic then.
func (a *HLL) AddBytes(value []byte) {
	a.numValues++
	a.hll.Add(fingerprint.Hash64(value))
}

// TODO: uncomment when/if implemented in internal/hllplus.
//
// // Merge merges another HyperLogLog++ aggregator into current one.
// func (a *HLL) Merge(another *HLL) error {
// 	if err := a.hll.Merge(another.hll); err != nil {
// 		return err
// 	}
// 	a.numValues += another.numValues
// 	return nil
// }

// Marshal marshals an aggregator to a binary proto message (raw bytes, not base64 encoded).
//
// This binary representation is compatible with BigQuery HyperLogLog++:
// https://cloud.google.com/bigquery/docs/reference/standard-sql/hll_functions
func (a *HLL) Marshal() ([]byte, error) {
	data, sparseSize := a.hll.GetData()

	aggState := &zetasketch.AggregatorStateProto{
		Type:            &hllAggregatorType,
		NumValues:       &a.numValues,        // TODO: we must track each addition in a wrapper around HLL++
		EncodingVersion: &hllEncodingVersion, // fixed
		ValueType:       nil,                 // looks to be a type of values being added - strings, bytes, ints etc - I think, can be omitted (may need to check though)
	}

	var hllState *zetasketch.HyperLogLogPlusUniqueStateProto
	if sparseSize != -1 {
		// sparse:
		sparseSize := int32(sparseSize)
		hllState = &zetasketch.HyperLogLogPlusUniqueStateProto{
			SparseSize:                  &sparseSize,
			SparsePrecisionOrNumBuckets: &a.precision,
			SparseData:                  data,
		}
	} else {
		// dense:
		hllState = &zetasketch.HyperLogLogPlusUniqueStateProto{
			PrecisionOrNumBuckets: &a.precision,
			Data:                  data,
		}
	}
	proto.SetExtension(aggState, zetasketch.E_HyperloglogplusUniqueState, hllState)

	return proto.Marshal(aggState)
}
