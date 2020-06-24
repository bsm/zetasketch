package hash

import (
	"encoding/binary"
)

const (
	// hash seed values/components:
	c0 = 0xa5b85c5e198ed849
	c1 = 0x8d58ac26afe12e47
	c2 = 0xc47b6e9e3a970ed3
	c3 = 0xc6a4a7935bd1e995
)

// Bytes computes a hash of data.
//
// WARNING: it supports data with up to 64 bytes for now. Will panic if data is larger.
func Bytes(data []byte) uint64 {
	var h uint64
	if n := len(data); n <= 32 {
		h = mm64(data, c0^c1^c2)
	} else if n < 64 {
		h = hash33to64(data)
	} else {
		panic("hashing of data more than 64 bytes long is not implemented yet")
	}

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
		k := load64(data[b*8:]) * c3
		k = shiftMix(k) * c3

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

	h = shiftMix(h) * c3
	return shiftMix(h)
}

func hash128to64(hi, lo uint64) uint64 {
	h := (lo ^ hi) * c3
	h = (hi ^ shiftMix(h)) * c3
	return shiftMix(h) * c3
}

func hash33to64(data []byte) uint64 {
	z := load64(data[24:])
	a := load64(data) + (uint64(len(data))+load64(data[len(data)-16:]))*c0
	b := rotateRight(a+z, 52)
	c := rotateRight(a, 37)

	a += load64(data[8:])
	c += rotateRight(a, 7)
	a += load64(data[16:])
	vf := a + z
	vs := b + rotateRight(a, 31) + c

	a = load64(data[16:]) + load64(data[len(data)-32:])
	z = load64(data[len(data)-8:])
	b = rotateRight(a+z, 52)
	c = rotateRight(a, 37)

	a += load64(data[len(data)-24:])
	c += rotateRight(a, 7)
	a += load64(data[len(data)-16:])
	wf := a + z
	ws := b + rotateRight(a, 31) + c

	r := shiftMix((vf+ws)*c2 + (wf+vs)*c0)
	return shiftMix(r*c0+vs) * c2
}

func load64(data []byte) uint64 {
	return binary.LittleEndian.Uint64(data)
}

func rotateRight(v uint64, k uint) uint64 {
	return (v >> k) | (v << (64 - k))
}

func shiftMix(v uint64) uint64 {
	return v ^ (v >> 47)
}
