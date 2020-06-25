package hash_test

import (
	"bytes"

	"github.com/bsm/zetasketch/internal/hash"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bytes", func() {
	var foobar = []byte("foobar")

	// Partially taken from:
	// https://github.com/google/zetasketch/blob/master/javatests/com/google/zetasketch/internal/hash/HashTest.java
	It("should hash short sequences", func() {
		Expect(hash.Bytes(nil)).To(Equal(uint64(0x23ad7c904aa665e3)))
		Expect(hash.Bytes(foobar)).To(Equal(uint64(0x36a1e57a138e4467)))
		Expect(hash.Bytes([]byte("foo"))).To(Equal(uint64(0xd0bcbfe261b36504)))
		Expect(hash.Bytes([]byte("f"))).To(Equal(uint64(16100291902947574842)))
		Expect(hash.Bytes([]byte("Z\u00fcrich"))).To(Equal(uint64(0x27efc00f7d2ce548)))
		Expect(hash.Bytes([]byte("Zu\u0308rich"))).To(Equal(uint64(0x7dfa3067e55c7e8a)))
		Expect(hash.Bytes(bytes.Repeat(foobar, 3))).To(Equal(uint64(0xd7b08d94eeefdccd)))
	})

	It("should hash longer sequences", func() {
		Expect(hash.Bytes(bytes.Repeat(foobar, 8))).To(Equal(uint64(0x94386e8403038649)))
		Expect(hash.Bytes(bytes.Repeat(foobar, 24))).To(Equal(uint64(0xd019f3291f1d4d37)))
		Expect(hash.Bytes(bytes.Repeat(foobar, 45))).To(Equal(uint64(0x549b2e228e80ee1a)))
		Expect(hash.Bytes(bytes.Repeat(foobar, 48))).To(Equal(uint64(0x874fce45a1e2a8ae)))
	})
})
