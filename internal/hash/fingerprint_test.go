package hash_test

import (
	"bytes"
	"testing"

	"github.com/bsm/zetasketch/internal/hash"
)

func TestBytes(t *testing.T) {
	foobar := []byte("foobar")

	// Partially taken from:
	// https://github.com/google/zetasketch/blob/master/javatests/com/google/zetasketch/internal/hash/HashTest.java
	cases := []struct {
		in  []byte
		exp uint64
	}{
		// short sequences
		{nil, 0x23ad7c904aa665e3},
		{foobar, 0x36a1e57a138e4467},
		{[]byte("foo"), 0xd0bcbfe261b36504},
		{[]byte("f"), 16100291902947574842},
		{[]byte("Zürich"), 0x27efc00f7d2ce548},
		{[]byte("Zürich"), 0x7dfa3067e55c7e8a},
		{bytes.Repeat(foobar, 3), 0xd7b08d94eeefdccd},

		// longer sequences
		{bytes.Repeat(foobar, 8), 0x94386e8403038649},
		{bytes.Repeat(foobar, 24), 0xd019f3291f1d4d37},
		{bytes.Repeat(foobar, 45), 0x549b2e228e80ee1a},
		{bytes.Repeat(foobar, 48), 0x874fce45a1e2a8ae},
	}
	for _, tc := range cases {
		if got := hash.Bytes(tc.in); got != tc.exp {
			t.Errorf("Bytes(%q) = %#x, want %#x", tc.in, got, tc.exp)
		}
	}
}
