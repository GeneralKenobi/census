package util

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

// RngSeed generates a random seed using the crypto random number generator.
func RngSeed() (int64, error) {
	var seedBytes [8]byte
	_, err := rand.Read(seedBytes[:])
	if err != nil {
		return 0, fmt.Errorf("error generating random seed: %w", err)
	}

	seed := binary.LittleEndian.Uint64(seedBytes[:])
	return int64(seed), nil
}
