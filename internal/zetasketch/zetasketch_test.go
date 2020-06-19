package zetasketch_test

import (
	"encoding/base64"
	"testing"

	. "github.com/bsm/zetasketch"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HLL", func() {
	var subject *HLL

	BeforeEach(func() {
		var err error
		subject, err = NewHLL(10)
		Expect(err).NotTo(HaveOccurred())

		// real count = 3 uniques:
		subject.Add([]byte("foo"))
		subject.Add([]byte("bar"))
		subject.Add([]byte("foo")) // dupe
		subject.Add([]byte("baz"))
	})

	It("should return estimated count", func() {
		Expect(subject.Count()).To(BeNumerically("==", 3))
	})

	It("should merge another aggregator", func() {
		another, err := NewHLL(10)
		Expect(err).NotTo(HaveOccurred())
		another.Add([]byte("baz"))
		another.Add([]byte("baz"))
		another.Add([]byte("qux"))

		Expect(subject.Merge(another)).To(Succeed())
		Expect(subject.Count()).To(BeNumerically("==", 4)) // 3 uniques from `subject`, 2 uniques in `another`, but 1 overlap
	})

	It("should clear", func() {
		Expect(subject.Count()).To(BeNumerically("==", 3))
		subject.Clear()
		Expect(subject.Count()).To(BeZero())
	})

	It("should marshal", func() {
		data, err := subject.Marshal()
		Expect(err).NotTo(HaveOccurred())

		Expect(base64.StdEncoding.EncodeToString(data)).To(Equal(
			// THIS FAILS IN BQ:
			// "Invalid input bytes to HLL_COUNT.EXTRACT":
			//
			// SELECT HLL_COUNT.EXTRACT(
			//   FROM_BASE64("ggcCGAoIcBAEGAI=")
			// )
			//
			// and it really looks way too short - haven't figured out what's wrong yet
			`ggcCGAoIcBAEGAI=`,
		))
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "zetasketch")
}
