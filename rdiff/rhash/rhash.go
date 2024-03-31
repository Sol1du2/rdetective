package rhash

import (
	"fmt"
)

const (
	moduloPrime = 65521
)

// RollingHash represents a rolling hash computation using the adler32
// algorithm.
type RollingHash struct {
	a, b   uint32
	Size   int
	Window []byte
}

func New() *RollingHash {
	return &RollingHash{a: 1, b: 0, Size: 0}
}

// Update updates the rolling hash with the next byte.
func (r *RollingHash) Update(b byte) {
	r.a = (r.a + uint32(b)) % moduloPrime
	r.b = (r.b + r.a) % moduloPrime

	r.Window = append(r.Window, b)
	r.Size++
}

// Roll removes the first byte from the rolling hash.
func (r *RollingHash) Roll() (byte, error) {
	if len(r.Window) == 0 {
		return 0, fmt.Errorf("nothing to roll out") // Nothing to roll out
	}

	old := r.Window[0]

	r.a = (r.a + moduloPrime - uint32(old)) % moduloPrime
	r.b = (r.b + (uint32(len(r.Window))*uint32(old))/moduloPrime - (uint32(len(r.Window)) * uint32(old)) - 1) % moduloPrime

	r.Window = r.Window[1:]
	r.Size--

	return old, nil
}

// Returns the current hash value.
func (r *RollingHash) Sum() uint32 {
	return (r.b << 16) | r.a
}

// Resets the rolling hash calculations.
func (r *RollingHash) Reset() {
	r.a = 1
	r.b = 0
	r.Size = 0
	r.Window = []byte{}
}
