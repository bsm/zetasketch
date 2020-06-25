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
func (h *HLL) Result() int64 {
	return h.h.Estimate()
}

// Proto returns a marshalable protobuf message.
func (h *HLL) Proto() *pb.AggregatorStateProto {
	var (
		encodingVersion int32 = 2
		aggType               = pb.AggregatorType_HYPERLOGLOG_PLUS_UNIQUE
		numValues             = int64(h.n)
	)
	msg := &pb.AggregatorStateProto{
		Type:            &aggType,
		EncodingVersion: &encodingVersion,
		NumValues:       &numValues,
	}
	proto.SetExtension(msg, zetasketch.E_HyperloglogplusUniqueState, h.h.Proto())
	return msg
}

// FromProto populates aggregator from given proto message.
// It overrides existing state.
func (h *HLL) FromProto(msg *pb.AggregatorStateProto) error {
	if msg.GetType() != pb.AggregatorType_HYPERLOGLOG_PLUS_UNIQUE {
		return fmt.Errorf("incompatible binary message: unexpected type %s", msg.GetType().String())
	}
	if msg.GetEncodingVersion() != 2 {
		return fmt.Errorf("incompatible binary message: unsupported encoding version %#v", msg.GetEncodingVersion())
	}
	if msg.NumValues == nil {
		return fmt.Errorf("incompatible binary message: no num values")
	}

	ext := proto.GetExtension(msg, zetasketch.E_HyperloglogplusUniqueState)
	hState, ok := ext.(*pb.HyperLogLogPlusUniqueStateProto)
	if !ok {
		return fmt.Errorf("incompatible binary message: invalid HyperLogLog++ state")
	}

	hll, err := hllplus.NewFromProto(hState)
	if err != nil {
		return err
	}

	h.h = hll
	h.n = uint64(msg.GetNumValues())
	return nil
}

// MarshalBinary serializes aggregator to bytes.
func (h *HLL) MarshalBinary() ([]byte, error) {
	return proto.Marshal(h.Proto())
}

// UnmarshalBinary deserializes aggregator from bytes.
func (h *HLL) UnmarshalBinary(data []byte) error {
	msg := new(pb.AggregatorStateProto)
	if err := proto.Unmarshal(data, msg); err != nil {
		return err
	}
	return h.FromProto(msg)
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
