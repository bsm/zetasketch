package zetasketch_test

import (
	"github.com/bsm/zetasketch"
	pb "github.com/bsm/zetasketch/internal/zetasketch"
	"google.golang.org/protobuf/proto"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HLL", func() {
	var subject *zetasketch.HLL
	var _ zetasketch.Aggregator = subject

	BeforeEach(func() {
		subject = zetasketch.NewHLL(nil)

		for i := 0; i < 1_000; i++ {
			subject.Add(zetasketch.Uint64Value(uint64(i)))
		}
		for i := 500; i < 1_000; i++ {
			subject.Add(zetasketch.Uint64Value(uint64(i)))
		}
	})

	It("should count values", func() {
		Expect(subject.NumValues()).To(Equal(uint64(1_500)))
	})

	It("should estimate uniques", func() {
		Expect(subject.Result()).To(Equal(int64(1_003)))
	})

	It("should merge", func() {
		other := zetasketch.NewHLL(nil)
		for i := 800; i < 1_200; i++ {
			other.Add(zetasketch.Uint64Value(uint64(i)))
		}

		Expect(subject.Merge(other)).To(Succeed())
		Expect(subject.NumValues()).To(Equal(uint64(1_900)))
		Expect(subject.Result()).To(Equal(int64(1_207)))

		// `other` is not modified:
		Expect(other.NumValues()).To(Equal(uint64(400)))
		Expect(other.Result()).To(Equal(int64(400)))
	})

	It("should return / init from protobuf", func() {
		msg := subject.Proto()

		Expect(*msg.EncodingVersion).To(BeNumerically("==", 2))                // fixed/const
		Expect(*msg.Type).To(Equal(pb.AggregatorType_HYPERLOGLOG_PLUS_UNIQUE)) // fixed/const
		Expect(*msg.NumValues).To(BeNumerically("==", 1_500))
		Expect(msg.ValueType).To(BeNil()) // we don't populate it

		// check that we do not forget to populate HLL-specific extension:
		ext := proto.GetExtension(msg, pb.E_HyperloglogplusUniqueState)
		Expect(ext).NotTo(BeNil())

		// check that it can init back from proto message:
		subject := new(zetasketch.HLL)
		Expect(subject.FromProto(msg)).To(Succeed())
		Expect(subject.NumValues()).To(BeNumerically("==", 1_500))
		Expect(subject.Result()).To(BeNumerically("==", 1_003))

		// basic check for wrapper method to marshal this proto:
		data, err := subject.MarshalBinary()
		Expect(err).NotTo(HaveOccurred())
		Expect(len(data)).To(Equal(32786)) // on failure, this provides nicer message than HaveLen

		// and unmarshal:
		subject2 := new(zetasketch.HLL)
		Expect(subject2.UnmarshalBinary(data)).To(Succeed())
		Expect(subject.NumValues()).To(BeNumerically("==", 1_500))
		Expect(subject.Result()).To(BeNumerically("==", 1_003))
	})
})
