package transform

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func HashSlice[T comparable](slice []T) string {
	hasher := sha256.New()

	// Iterate over each element of the slice
	for _, elem := range slice {
		// Convert each element to bytes
		// You may need to handle different types more explicitly in a production scenario
		_, err := hasher.Write([]byte(fmt.Sprintf("%v", elem)))
		if err != nil {
			panic(err)
		}
	}

	// Compute the hash and return it as a hex string
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
