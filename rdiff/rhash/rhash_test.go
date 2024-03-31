package rhash

import (
	"testing"
)

// Hash taken from https://en.wikipedia.org/wiki/Adler-32
func TestUpdate(t *testing.T) {
	wiki := []byte("Wikipedia")
	expectedHash := uint32(300286872)

	rh := New()

	for _, b := range wiki {
		rh.Update(b)
	}

	hash := rh.Sum()
	if expectedHash != hash {
		t.Errorf("calculated hash different than expected, got %d, expected %d", hash, expectedHash)
	}
}

func TestReset(t *testing.T) {
	helloThere := []byte("Hello Thy World")

	rh := New()

	for _, b := range helloThere {
		rh.Update(b)
	}

	if rh.Sum() == 1 {
		t.Errorf("hash not calculated")
	}

	rh.Reset()

	sum := rh.Sum()
	if rh.Sum() != 1 {
		t.Errorf("hash was not reset, got sum %d", sum)
	}
}

func TestRoll(t *testing.T) {
	helloWorld := []byte("Hello World")
	world := []byte("World")

	rhHelloWorld := New()
	rhWorld := New()

	for _, b := range helloWorld {
		rhHelloWorld.Update(b)
	}

	for _, b := range world {
		rhWorld.Update(b)
	}

	for i := 0; i < 6; i++ {
		b, err := rhHelloWorld.Roll()
		if err != nil {
			t.Errorf("has is empty: %s", err.Error())
		}

		if b != helloWorld[i] {
			t.Errorf("unexpected rolled out byte, expected %d, got %d", helloWorld[i], b)
		}
	}

	rolledHash := rhHelloWorld.Sum()
	worldHash := rhWorld.Sum()

	if rolledHash != worldHash {
		t.Errorf("rolled hash different than non rolled hash, rolledHash=%d, nonRolledHash=%d", rolledHash, worldHash)
	}
}
