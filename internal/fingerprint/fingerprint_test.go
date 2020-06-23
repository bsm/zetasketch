package fingerprint_test

import (
	"testing"

	"github.com/bsm/zetasketch/internal/fingerprint"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hash64", func() {
	// Simplified copy of:
	// https://github.com/google/zetasketch/blob/master/javatests/com/google/zetasketch/internal/hash/HashTest.java
	It("should hash", func() {
		Expect(fingerprint.Hash64(nil)).To(Equal(uint64(0x23ad7c904aa665e3)))
		Expect(fingerprint.Hash64([]byte{0x66, 0x6f, 0x6f, 0x62, 0x61, 0x72})).To(Equal(uint64(0x36a1e57a138e4467)))
		Expect(fingerprint.Hash64([]byte("foo"))).To(Equal(uint64(0xd0bcbfe261b36504)))
		Expect(fingerprint.Hash64([]byte("Z\u00fcrich"))).To(Equal(uint64(0x27efc00f7d2ce548)))
		Expect(fingerprint.Hash64([]byte("Zu\u0308rich"))).To(Equal(uint64(0x7dfa3067e55c7e8a)))
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "fingerprint")
}
