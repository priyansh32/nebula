package coordinator

import (
	// "hash/fnv"

	"crypto/sha256"
	"encoding/binary"
	"errors"

	"github.com/google/uuid"
)

type node struct {
	hash        uint64
	storeClient *StoreClient
}

type HashRing struct {
	nodes             map[uint64]*node
	replicationFactor int
}

func NewHashRing(r int) *HashRing {
	return &HashRing{
		nodes:             make(map[uint64]*node),
		replicationFactor: r,
	}
}

// hashKey hashes the key to determine its position on the ring
func hashKey(key string) uint64 {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	hashBytes := hasher.Sum(nil)
	// Take the first 8 bytes (64 bits) and convert to uint64
	truncatedHash := binary.BigEndian.Uint64(hashBytes[:8])

	return truncatedHash
}

// adds a store to the ring with count nodes
func (hr *HashRing) AddStoreNodes(s *StoreClient) {

	for i := 0; i < hr.replicationFactor; i++ {

		// identify the position of the node on the ring
		hash := hashKey(s.name + "-" + uuid.New().String())

		n := node{
			hash:        hash,
			storeClient: s,
		}

		hr.nodes[hash] = &n
		s.nodeKeys = append(s.nodeKeys, hash)
	}
}

// finds the store for the given key
func (hr *HashRing) GetStore(key string) (*StoreClient, error) {
	hash := hashKey(key)

	if len(hr.nodes) == 0 {
		return nil, errors.New("no stores in the ring")
	}

	var s *StoreClient

	var diff uint64
	diff = ^uint64(0)

	for _, v := range hr.nodes {

		var curr uint64

		if v.hash >= hash {
			curr = v.hash - hash
		} else {
			// wrap around
			curr = hash - v.hash
			curr = ^curr
		}

		// find the node with the smallest difference forward
		if diff > curr {
			diff = curr
			s = v.storeClient
		}
	}

	return s, nil
}
