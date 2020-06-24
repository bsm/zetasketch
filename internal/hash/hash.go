// Package hash implements hashing as done in zetasketch Java library
// https://github.com/google/zetasketch/blob/master/java/com/google/zetasketch/internal/hash/.
package hash

import (
	"encoding/binary"
)

// Uint32 hashes uint32 numbers.
func Uint32(v uint32) uint64 {
	buf := make([]byte, 5)
	binary.LittleEndian.PutUint32(buf, v)
	return Bytes(buf)
}

// Uint64 hashes uint64 numbers.
func Uint64(v uint64) uint64 {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, v)
	return Bytes(buf)
}

// String hashes strings.
func String(v string) uint64 {
	return Bytes([]byte(v))
}
