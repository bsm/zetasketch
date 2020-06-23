package hllplus_test

import (
	"github.com/bsm/zetasketch/hllplus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Data", func() {
	It("should estimate bias", func() {
		Expect(hllplus.EstimateBias(0, 15)).To(BeNumerically("==", 0.0))
		Expect(hllplus.EstimateBias(1, 15)).To(BeNumerically("==", 0.0))
		Expect(hllplus.EstimateBias(10_000, 15)).To(BeNumerically("==", 0.0))
		Expect(hllplus.EstimateBias(100_000, 15)).To(BeNumerically("~", 888.1, 0.1))
		Expect(hllplus.EstimateBias(200_000, 15)).To(BeNumerically("==", 0.0))

		Expect(hllplus.EstimateBias(50_000, 13)).To(BeNumerically("==", 0.0))
		Expect(hllplus.EstimateBias(50_000, 14)).To(BeNumerically("~", 449.7, 0.1))
		Expect(hllplus.EstimateBias(50_000, 15)).To(BeNumerically("~", 7820.2, 0.1))
		Expect(hllplus.EstimateBias(50_000, 16)).To(BeNumerically("~", 44513.2, 0.1))
		Expect(hllplus.EstimateBias(50_000, 17)).To(BeNumerically("==", 0.0))
	})
})
