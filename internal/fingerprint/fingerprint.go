// Package fingerprint implements fingerprinting/hashing as done in zetasketch Java library
// https://github.com/google/zetasketch/blob/master/java/com/google/zetasketch/internal/hash/Fingerprint2011.java
package fingerprint

import (
	"encoding/binary"
)

const (
	// hash seed values/components:
	c0 = 0xa5b85c5e198ed849
	c1 = 0x8d58ac26afe12e47
	c2 = 0xc47b6e9e3a970ed3
	c3 = 0xc6a4a7935bd1e995

	rs = 47 // right shift bits
)

// Hash64 computes a hash of data.
func Hash64(data []byte) uint64 {
	/*
		Java implementation looks like this:

		if (length <= 32) {
			result = murmurHash64WithSeed(bytes, offset, length, K0 ^ K1 ^ K2);
		} else if (length <= 64) {
			result = hashLength33To64(bytes, offset, length);
		} else {
			result = fullFingerprint(bytes, offset, length);
		}
	*/
	if len(data) > 32 {
		panic("hashing of data more than 32 bytes long is not implemented yet")
	}

	h := mm64(data, c0^c1^c2)

	var u, v uint64 = c0, c0
	if len(data) >= 8 {
		u = binary.LittleEndian.Uint64(data)
	}

	if len(data) >= 9 {
		v = binary.LittleEndian.Uint64(data[len(data)-8:])
	}

	h = hash128to64(h+v, u)
	if h == 0 || h == 1 {
		return h + ^uint64(1)
	}
	return h
}

// mm64 computes 64-bit Murmur hash with given seed.
func mm64(data []byte, seed uint64) uint64 {
	h := seed ^ uint64(len(data))*c3

	nblocks := len(data) / 8
	for b := 0; b < nblocks; b++ {
		k := binary.LittleEndian.Uint64(data[b*8:])
		k *= c3
		k ^= k >> rs
		k *= c3

		h ^= k
		h *= c3
	}

	tail := data[nblocks*8:]
	switch len(tail) & 7 {
	case 7:
		h ^= uint64(tail[6]) << 48
		fallthrough
	case 6:
		h ^= uint64(tail[5]) << 40
		fallthrough
	case 5:
		h ^= uint64(tail[4]) << 32
		fallthrough
	case 4:
		h ^= uint64(tail[3]) << 24
		fallthrough
	case 3:
		h ^= uint64(tail[2]) << 16
		fallthrough
	case 2:
		h ^= uint64(tail[1]) << 8
		fallthrough
	case 1:
		h ^= uint64(tail[0])
		h *= c3
	}

	h ^= h >> rs
	h *= c3
	h ^= h >> rs
	return h
}

func hash128to64(hi, lo uint64) uint64 {
	h := (lo ^ hi) * c3
	h ^= h >> rs
	h = (hi ^ h) * c3
	h ^= h >> rs
	h *= c3
	return h
}
