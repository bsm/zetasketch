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

	DescribeTable("dense (800 unique)",
		func(p int, exp int) {
			subject, _ = hllplus.New(uint8(p), 20)
			rnd = rand.New(rand.NewSource(33))
			for i := 0; i < 800; i++ {
				subject.Add(rnd.Uint64())
			}
			Expect(subject.Estimate()).To(Equal(uint64(exp)))
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

	DescribeTable("dense (200k unique)",
		func(p int, exp int) {
			subject, _ = hllplus.New(uint8(p), 20)
			rnd = rand.New(rand.NewSource(33))
			for i := 0; i < 200_000; i++ {
				subject.Add(rnd.Uint64())
			}
			Expect(subject.Estimate()).To(Equal(uint64(exp)))
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

	DescribeTable("dense (100k unique + 2x50k)",
		func(p int, exp int) {
			subject, _ = hllplus.New(uint8(p), 20)
			rnd = rand.New(rand.NewSource(33))
			for i := 0; i < 100_000; i++ {
				subject.Add(rnd.Uint64())
			}
			for i := 0; i < 50_000; i++ {
				h := rnd.Uint64()
				subject.Add(h)
				subject.Add(h)
			}
			Expect(subject.Estimate()).To(Equal(uint64(exp)))
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
})

// --------------------------------------------------------------------

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "zetasketch/hllplus")
}
