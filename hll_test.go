package zetasketch_test

import (
	"github.com/bsm/zetasketch"

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
		Expect(subject.NumValues()).To(BeNumerically("==", 1_500))
	})

	It("should estimate uniques", func() {
		Expect(subject.Result()).To(BeNumerically("==", 1_003))
	})

	It("should merge", func() {
		other := zetasketch.NewHLL(nil)
		for i := 800; i < 1_200; i++ {
			other.Add(zetasketch.Uint64Value(uint64(i)))
		}

		Expect(subject.Merge(other)).To(Succeed())
		Expect(subject.NumValues()).To(BeNumerically("==", 1_900))
		Expect(subject.Result()).To(BeNumerically("==", 1_207))

		// `other` is not modified:
		Expect(other.NumValues()).To(BeNumerically("==", 400))
		Expect(other.Result()).To(BeNumerically("==", 400))
	})

	It("should marshal/unmarshal binary", func() {
		data, err := subject.MarshalBinary()
		Expect(err).NotTo(HaveOccurred())

		subject = new(zetasketch.HLL)
		Expect(subject.UnmarshalBinary(data)).To(Succeed())
		Expect(subject.NumValues()).To(BeNumerically("==", 1_500))
		Expect(subject.Result()).To(BeNumerically("==", 1_003))
	})
})
