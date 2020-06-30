package hllplus_test

import (
	"math/rand"
	"testing"

	"github.com/bsm/zetasketch/hllplus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("HLL", func() {
	var subject *hllplus.HLL
	var rnd *rand.Rand

	BeforeEach(func() {
		rnd = rand.New(rand.NewSource(33))
	})

	DescribeTable("estimate normal (800 unique)",
		func(p int, exp int) {
			subject, _ = hllplus.NewNormal(uint8(p))
			for i := 0; i < 800; i++ {
				subject.Add(rnd.Uint64())
			}
			Expect(subject.IsSparse()).To(BeFalse())
			Expect(subject.Estimate()).To(Equal(int64(exp)))
		},
		Entry("p=10", 10, 800),
		Entry("p=11", 11, 794),
		Entry("p=12", 12, 788),
		Entry("p=13", 13, 793),
		Entry("p=14", 14, 791),
		Entry("p=15", 15, 793),
		Entry("p=16", 16, 795),
		Entry("p=17", 17, 797),
		Entry("p=18", 18, 799),
	)

	DescribeTable("estimate sparse (800 unique)",
		func(p int, exp int) {
			subject, _ = hllplus.New(uint8(p-5), uint8(p))
			for i := 0; i < 800; i++ {
				subject.Add(rnd.Uint64())
			}
			Expect(subject.IsSparse()).To(BeTrue())
			Expect(subject.Estimate()).To(Equal(int64(exp)))
		},
		Entry("p=16", 16, 799),
		Entry("p=17", 17, 798),
		Entry("p=18", 18, 799),
		Entry("p=19", 19, 799),
		Entry("p=20", 20, 799),
		Entry("p=21", 21, 799),
		Entry("p=22", 22, 800),
		Entry("p=23", 23, 800),
	)

	DescribeTable("estimate normal (200k unique)",
		func(p int, exp int) {
			subject, _ = hllplus.NewNormal(uint8(p))
			for i := 0; i < 200_000; i++ {
				subject.Add(rnd.Uint64())
			}
			Expect(subject.IsSparse()).To(BeFalse())
			Expect(subject.Estimate()).To(Equal(int64(exp)))
		},
		Entry("p=10", 10, 199197),
		Entry("p=11", 11, 204855),
		Entry("p=12", 12, 204356),
		Entry("p=13", 13, 202958),
		Entry("p=14", 14, 202977),
		Entry("p=15", 15, 201496),
		Entry("p=16", 16, 201398),
		Entry("p=17", 17, 201208),
		Entry("p=18", 18, 200999),
		Entry("p=19", 19, 200567),
		Entry("p=20", 20, 200032),
		Entry("p=21", 21, 200013),
		Entry("p=22", 22, 200003),
		Entry("p=23", 23, 199989),
		Entry("p=24", 24, 200026),
	)

	DescribeTable("estimate sparse (200k unique)",
		func(p int, exp int) {
			subject, _ = hllplus.New(uint8(p-5), uint8(p))
			for i := 0; i < 200_000; i++ {
				subject.Add(rnd.Uint64())
			}
			Expect(subject.IsSparse()).To(BeTrue())
			Expect(subject.Estimate()).To(Equal(int64(exp)))
		},
		Entry("p=24", 24, 200535),
		Entry("p=25", 25, 200048),
	)

	DescribeTable("estimate normal (100k unique + 2x50k)",
		func(p int, exp int) {
			subject, _ = hllplus.NewNormal(uint8(p))
			for i := 0; i < 100_000; i++ {
				subject.Add(rnd.Uint64())
			}
			for i := 0; i < 50_000; i++ {
				h := rnd.Uint64()
				subject.Add(h)
				subject.Add(h)
			}
			Expect(subject.IsSparse()).To(BeFalse())
			Expect(subject.Estimate()).To(Equal(int64(exp)))
		},
		Entry("p=10", 10, 148131),
		Entry("p=11", 11, 152153),
		Entry("p=12", 12, 152338),
		Entry("p=13", 13, 152453),
		Entry("p=14", 14, 150853),
		Entry("p=15", 15, 150458),
		Entry("p=16", 16, 150811),
		Entry("p=17", 17, 150795),
		Entry("p=18", 18, 150592),
		Entry("p=19", 19, 150265),
		Entry("p=20", 20, 149879),
		Entry("p=21", 21, 149939),
		Entry("p=22", 22, 150005),
		Entry("p=23", 23, 149946),
		Entry("p=24", 24, 149988),
	)

	DescribeTable("estimate sparse (100k unique + 2x50k)",
		func(p int, exp int) {
			subject, _ = hllplus.New(uint8(p-5), uint8(p))
			for i := 0; i < 100_000; i++ {
				subject.Add(rnd.Uint64())
			}
			for i := 0; i < 50_000; i++ {
				h := rnd.Uint64()
				subject.Add(h)
				subject.Add(h)
			}
			Expect(subject.IsSparse()).To(BeTrue())
			Expect(subject.Estimate()).To(Equal(int64(exp)))
		},
		Entry("p=23", 23, 150823),
		Entry("p=24", 24, 150144),
		Entry("p=25", 25, 150012),
	)

	FIt("should normalize", func() {
		subject, _ = hllplus.New(12, 17)
		for i := 0; i < 3_084; i++ {
			subject.Add(rnd.Uint64())
		}
		Expect(subject.IsSparse()).To(BeTrue())
		Expect(subject.Estimate()).To(BeNumerically("==", 3_108))

		subject.Add(rnd.Uint64())
		Expect(subject.IsSparse()).To(BeFalse())
		Expect(subject.Estimate()).To(BeNumerically("==", 3_100))
	})

	It("should downgrade", func() {
		s1, _ := hllplus.NewNormal(15)
		s2, _ := hllplus.NewNormal(12)

		for i := 0; i < 100_000; i++ {
			n := rnd.Uint64()
			s1.Add(n)
			s2.Add(n)
		}

		Expect(s1.Estimate()).To(Equal(int64(99879)))
		Expect(s2.Estimate()).To(Equal(int64(100680)))

		Expect(s1.Precision()).To(Equal(uint8(15)))
		Expect(s1.SparsePrecision()).To(Equal(uint8(20)))
		Expect(s1.Downgrade(12, 17)).To(Succeed())
		Expect(s1.Precision()).To(Equal(uint8(12)))
		Expect(s1.SparsePrecision()).To(Equal(uint8(17)))

		Expect(s1.Estimate()).To(Equal(int64(100680)))
		Expect(s2.Estimate()).To(Equal(int64(100680)))
	})

	Describe("merge", func() {
		var s1, s2, s3 *hllplus.HLL

		BeforeEach(func() {
			s1, _ = hllplus.NewNormal(15)
			s2, _ = hllplus.NewNormal(15)
			s3, _ = hllplus.NewNormal(12)

			for i := 0; i < 50_000; i++ {
				n := rnd.Uint64()
				s1.Add(n)
				s2.Add(n)
				s3.Add(n)
			}
			for i := 0; i < 50_000; i++ {
				s1.Add(rnd.Uint64())
				s2.Add(rnd.Uint64())
				s3.Add(rnd.Uint64())
			}

			Expect(s1.Estimate()).To(Equal(int64(100324)))
			Expect(s2.Estimate()).To(Equal(int64(100168)))
			Expect(s3.Estimate()).To(Equal(int64(100464)))
		})

		It("should support equal precision", func() {
			s1.Merge(s2)
			Expect(s1.Estimate()).To(Equal(int64(150794)))
		})

		It("should support targets with lower precision", func() {
			s1.Merge(s3)
			Expect(s1.Estimate()).To(Equal(int64(154744)))
			Expect(s1.Precision()).To(Equal(uint8(12)))
			Expect(s1.SparsePrecision()).To(Equal(uint8(17)))
		})

		It("should support targets with higher precision", func() {
			s3.Merge(s1)
			Expect(s3.Estimate()).To(Equal(int64(154744)))
			Expect(s3.Precision()).To(Equal(uint8(12)))
			Expect(s3.SparsePrecision()).To(Equal(uint8(17)))
		})

		It("should succeed if target is empty", func() {
			subject, _ = hllplus.NewNormal(15)
			Expect(func() { subject.Merge(s1) }).NotTo(Panic())

			// just a straight copy of s1:
			Expect(subject.Estimate()).To(Equal(s1.Estimate()))
			Expect(subject.Precision()).To(Equal(s1.Precision()))
			Expect(subject.SparsePrecision()).To(Equal(s1.SparsePrecision()))
		})
	})

	It("should init from proto", func() {
		subject, _ = hllplus.New(12, 17)
		for i := 0; i < 10_000; i++ {
			subject.Add(rnd.Uint64())
		}
		Expect(subject.IsSparse()).To(BeFalse())
		Expect(subject.Estimate()).To(BeNumerically("==", 6_843))

		msg := subject.Proto()

		// both precisions are always stored:
		Expect(msg.GetPrecisionOrNumBuckets()).To(BeNumerically("==", 12))
		Expect(msg.GetSparsePrecisionOrNumBuckets()).To(BeNumerically("==", 17))

		// expect dense representation:
		Expect(msg.GetData()).NotTo(BeEmpty()) // TODO: maybe better check exact value?

		// expect NO sparse representation:
		Expect(msg.SparseSize).To(BeNil())
		Expect(msg.GetSparseData()).To(BeNil())

		// init back from proto:
		subject, err := hllplus.NewFromProto(msg)
		Expect(err).NotTo(HaveOccurred())
		Expect(subject.IsSparse()).To(BeFalse())
		Expect(subject.Precision()).To(BeNumerically("==", 12))
		Expect(subject.SparsePrecision()).To(BeNumerically("==", 17))
		Expect(subject.Estimate()).To(BeNumerically("==", 6_843))
	})
})

// --------------------------------------------------------------------

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "zetasketch/hllplus")
}
