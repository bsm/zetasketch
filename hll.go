package zetasketch

import (
	"fmt"

	"github.com/bsm/zetasketch/hllplus"
	"github.com/bsm/zetasketch/internal/zetasketch"
	pb "github.com/bsm/zetasketch/internal/zetasketch"
	"google.golang.org/protobuf/proto"
)

// HLL implements a HLL++ aggregator for estimating cardinalities of multisets.
//
// The precision defines the accuracy of the HLL++ aggregator at the cost of the memory used. The
// upper bound on the memory required is 2^precision bytes, but less memory is used for
// smaller cardinalities. The relative error is 1.04 / sqrt(2^precision).
// A typical value used at Google is 15, which gives an error of about 0.6% while requiring an upper
// bound of 32KiB of memory.
//
// Note that this aggregator is not designed to be thread safe.
type HLL struct {
	h *hllplus.HLL
	n uint64
}

// NewHLL inits a new HLL++ aggregator.
func NewHLL(cfg *HLLConfig) *HLL {
	h, err := hllplus.New(cfg.precision(), cfg.sparsePrecision())
	if err != nil {
		panic(err)
	}
	return &HLL{h: h}
}

// Add adds value v to the aggregator.
func (h *HLL) Add(v Value) {
	h.n++
	h.h.Add(v.Sum64())
}

// NumValues returns the number of values seen.
func (h *HLL) NumValues() uint64 {
	return h.n
}

// Merge merges aggregator other into h.
func (h *HLL) Merge(other Aggregator) error {
	h2, ok := other.(*HLL)
	if !ok {
		return fmt.Errorf("cannot merge %T into %T", other, h)
	}

	h.h.Merge(h2.h)
	h.n += h2.n
	return nil
}

// Result returns an estimate of the unique of values.
func (h *HLL) Result() uint64 {
	return h.h.Estimate()
}

// MarshalBinary implements encoding.BinaryMarshaler interface.
func (h *HLL) MarshalBinary() ([]byte, error) {
	return proto.Marshal(h.Proto())
}

// Proto returns a marshalable protobuf message.
func (h *HLL) Proto() proto.Message {
	var (
		encodingVersion int32 = 2
		aggType               = pb.AggregatorType_HYPERLOGLOG_PLUS_UNIQUE
		numValues             = int64(h.n)
	)
	msg := &pb.AggregatorStateProto{
		Type:            &aggType,
		EncodingVersion: &encodingVersion,
		NumValues:       &numValues,
		ValueType:       nil, // looks to be a type of values being added - strings, bytes, ints etc - seems that BQ is OK with it being omitted
	}
	proto.SetExtension(msg, zetasketch.E_HyperloglogplusUniqueState, h.h.Proto())
	return msg
}

// -----------------------------------------------------------------------

// HLLConfig speficies the configuration parameters for the HLL++ aggregator.
type HLLConfig struct {
	// Defaults to 15.
	Precision uint8

	// If no sparse precision is specified, the default is calculated as precision + 5.
	SparsePrecision uint8
}

func (c *HLLConfig) precision() uint8 {
	if c != nil && c.Precision >= hllplus.MinPrecision && c.Precision <= hllplus.MaxPrecision {
		return c.Precision
	}
	return 15
}

func (c *HLLConfig) sparsePrecision() uint8 {
	min := c.precision()
	if c != nil && c.SparsePrecision >= min && c.SparsePrecision <= hllplus.MaxSparsePrecision {
		return c.SparsePrecision
	}
	if n := min + 5; n <= hllplus.MaxSparsePrecision {
		return n
	}
	return hllplus.MaxSparsePrecision
}
