package hash_test

import (
	"testing"

	"github.com/bsm/zetasketch/internal/hash"

	. "github.com/bsm/ginkgo"
	. "github.com/bsm/gomega"
)

var _ = Describe("Hash", func() {
	It("should hash uint64s", func() {
		Expect(hash.Uint64(0)).To(Equal(uint64(0x853a22bd6e14a48f)))
		Expect(hash.Uint64(1)).To(Equal(uint64(0xb91968b83211c978)))
		Expect(hash.Uint64(2)).To(Equal(uint64(0x83e2c1afe085d87a)))
		Expect(hash.Uint64((1 << 62) - 127)).To(Equal(uint64(0xfd2303b188e412d9)))
	})

	It("should hash uint32s", func() {
		Expect(hash.Uint32(0)).To(Equal(uint64(0x905df40cd02611cb)))
		Expect(hash.Uint32(1)).To(Equal(uint64(0xba4e724d9d787a26)))
		Expect(hash.Uint32(2)).To(Equal(uint64(0x6ef6414ac0c8858e)))
		Expect(hash.Uint32((1 << 29) - 127)).To(Equal(uint64(0xbda9eae625a1584)))
	})

	It("should hash strings", func() {
		Expect(hash.String("foo")).To(Equal(uint64(0xd0bcbfe261b36504)))
		Expect(hash.String("Z\u00fcrich")).To(Equal(uint64(0x27efc00f7d2ce548)))
		Expect(hash.String("Zu\u0308rich")).To(Equal(uint64(0x7dfa3067e55c7e8a)))
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "zetasketch/internal/hash")
}
