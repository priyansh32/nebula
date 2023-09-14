package coordinator

import (
	"crypto/sha256"
	"encoding/binary"
	"sort"
)

// hashKey hashes the key to determine its position on the ring
func hashKey(key string) uint64 {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	hashBytes := hasher.Sum(nil)
	// Take the first 8 bytes (64 bits) and convert to uint64
	truncatedHash := binary.BigEndian.Uint64(hashBytes[:8])

	return truncatedHash
}

func insertSorted(slice []uint64, element uint64) []uint64 {
	index := sort.Search(len(slice), func(i int) bool {
		return slice[i] >= element
	})

	// Insert the element at the found index
	slice = append(slice, 0)             // Add a new element to the slice
	copy(slice[index+1:], slice[index:]) // Shift elements to make space
	slice[index] = element               // Insert the new element

	return slice
}

func removeSorted(slice []uint64, element uint64) []uint64 {
	index := sort.Search(len(slice), func(i int) bool {
		return slice[i] >= element
	})

	if index < len(slice) && slice[index] == element {
		slice = append(slice[:index], slice[index+1:]...)
	}

	return slice
}
