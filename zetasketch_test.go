package zetasketch_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	. "github.com/bsm/zetasketch"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HLL", func() {
	var subject *HLL

	BeforeEach(func() {
		var err error
		subject, err = NewHLL(10, 10)
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
		another, err := NewHLL(10, 10)
		Expect(err).NotTo(HaveOccurred())
		another.Add([]byte("baz"))
		another.Add([]byte("baz"))
		another.Add([]byte("qux"))

		Expect(subject.Merge(another)).To(Succeed())
		Expect(subject.Count()).To(BeNumerically("==", 4)) // 3 uniques from `subject`, 2 uniques in `another`, but 1 overlap
	})

	// It("should clear", func() {
	// 	Expect(subject.Count()).To(BeNumerically("==", 3))
	// 	subject.Clear()
	// 	Expect(subject.Count()).To(BeZero())
	// })

	It("should marshal", func() {
		data, err := subject.Marshal()
		Expect(err).NotTo(HaveOccurred())

		// sparse:
		Expect(base64.StdEncoding.EncodeToString(data)).To(Equal(
			// DOES NOT WORK:
			// `SELECT HLL_COUNT.EXTRACT(FROM_BASE64("THIS VALUE"))`
			// EXPECTED: 3
			// GOT:      Invalid input bytes to HLL_COUNT.EXTRACT
			`ggcEEAAgCghwEAQYAg==`,
		))

		for i := 0; i < 400; i++ {
			subject.Add([]byte(fmt.Sprintf("%04d", i)))
		}

		data, err = subject.Marshal()
		Expect(err).NotTo(HaveOccurred())

		// dense:
		Expect(base64.StdEncoding.EncodeToString(data)).To(Equal(
			// DOES NOT WORK:
			// `SELECT HLL_COUNT.EXTRACT(FROM_BASE64("THIS VALUE"))`
			// EXPECTED: 406
			// GOT:      Invalid input bytes to HLL_COUNT.EXTRACT
			`ggftBBDDAiAKMuUEBYIC/AGAAoABgAiEA/wCggGGBXiAAoAChAj8BYAChAH8BIAFggGAAYIDfoACgAr+A4ABgAGGAvoFhAP8BoIF/gGAAoIBfoADgAWEAYABgAH8BYQB/ASAC4oF+gGAA/4I/gaCBP4DgAKCBP4JggiAB4oFevoDhAR8gAGEA3yEA3yAAYIBigH0DIABgAKAAYACgAOEAoIIhgP8BPwB/AGAAYICgAGCA/wCgAWCAYAJ/gGAA4ABkgXwA/4GgAiEBHyAA4ADggH+AYQGgAJ8gAKCAYQG+gWAA4QBfIABgAOIAfgEiAH4AYYB+geCBX6IAfoHhgH6AYQCggJ6ggT8AYIBggJ8hAGAD36AAf4DgAKGA4ABgAH6AYABggF+hAJ8ggP+AYYE/gKABX6CBfwGhAJ8iAd6/gKCA4AChAGEAfoEggH+AoABfIABgAOABoQE/AGEAXyCA4IB/AOAAYgEeoALgAGSBOwJggf+BYADiAn4A4AEgAGAAYQGgAGAAfwCggWCAv4CgAKACP4BhASAAfwCgA6AAYACggH+BIAEhASAA/wBgA+AAowCdIABgAGAAYIB/gGGAvwD/gOCB/4FggOAAYAC/gaAAoIEfoADgAGEBvwBhgb+A3yEAfwCgAGGBPwDgAOCAYAE/gKCAX5+hAH8AoADggSCCIAChgH6Av4BiAF6/gN+ggSCAv4F/gGAA4ADgAWGAfwBgAGICPgI/gOCAoADgAGABn6IAvgDhgL6A4IDggGCBH78BIYB+gGAAYABgAWAA4AHggKAAv4DgAKAAoIBhASGAnSCA/4EggmABP4ChAGAAghwEJQDGAI=`,
		))
	})

	It("should marshal to JSON", func() {
		data, err := subject.MarshalJSON()
		Expect(err).NotTo(HaveOccurred())

		Expect(data).To(HavePrefix(`"ggcE`))
		Expect(data).To(HaveSuffix(`QYAg=="`))
	})

	It("should refuse to unmarshal JSON", func() {
		payload := []byte(`""`) // payload does not matter
		Expect(subject.UnmarshalJSON(payload)).To(MatchError("marshalling HLL aggregator from JSON is not supported"))
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "zetasketch")
}
