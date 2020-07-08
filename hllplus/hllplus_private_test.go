package hllplus

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("sparseState", func() {
	var subject *sparseState

	BeforeEach(func() {
		subject = newSparseState(12, 17, nil)
	})

	It("should handle duplicates properly", func() {
		subject.Add(1)
		subject.Flush()
		Expect(subject.data.Count()).To(BeNumerically("==", 1))

		subject.Add(1)
		subject.Flush()
		Expect(subject.data.Count()).To(BeNumerically("==", 1))

		subject.Add(1)
		subject.Add(2)
		subject.Add(1)
		subject.Flush()
		Expect(subject.data.Count()).To(BeNumerically("==", 2))
	})
})

var _ = Describe("sparseState", func() {
	var subject *deltaSlice

	BeforeEach(func() {
		subject = recycleDeltaSlice(0) // pre-allocation doesn't matter in tests
	})

	It("should ignore if appended element is the same as last one", func() {
		subject.Append(1)
		Expect(subject.Count()).To(BeNumerically("==", 1))

		subject.Append(1)
		Expect(subject.Count()).To(BeNumerically("==", 1))

		subject.Append(1)
		subject.Append(2)
		Expect(subject.Count()).To(BeNumerically("==", 2))
	})
})
