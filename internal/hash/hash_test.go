package hash_test

import (
	"testing"

	"github.com/bsm/zetasketch/internal/hash"
)

func TestUint64(t *testing.T) {
	cases := []struct {
		in  uint64
		exp uint64
	}{
		{0, 0x853a22bd6e14a48f},
		{1, 0xb91968b83211c978},
		{2, 0x83e2c1afe085d87a},
		{(1 << 62) - 127, 0xfd2303b188e412d9},
	}
	for _, tc := range cases {
		if got := hash.Uint64(tc.in); got != tc.exp {
			t.Errorf("Uint64(%d) = %#x, want %#x", tc.in, got, tc.exp)
		}
	}
}

func TestUint32(t *testing.T) {
	cases := []struct {
		in  uint32
		exp uint64
	}{
		{0, 0x905df40cd02611cb},
		{1, 0xba4e724d9d787a26},
		{2, 0x6ef6414ac0c8858e},
		{(1 << 29) - 127, 0xbda9eae625a1584},
	}
	for _, tc := range cases {
		if got := hash.Uint32(tc.in); got != tc.exp {
			t.Errorf("Uint32(%d) = %#x, want %#x", tc.in, got, tc.exp)
		}
	}
}

func TestString(t *testing.T) {
	cases := []struct {
		in  string
		exp uint64
	}{
		{"foo", 0xd0bcbfe261b36504},
		{"Zürich", 0x27efc00f7d2ce548},
		{"Zürich", 0x7dfa3067e55c7e8a},
	}
	for _, tc := range cases {
		if got := hash.String(tc.in); got != tc.exp {
			t.Errorf("String(%q) = %#x, want %#x", tc.in, got, tc.exp)
		}
	}
}
