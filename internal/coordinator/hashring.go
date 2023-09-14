package coordinator

import (
	"errors"
	"sort"

	"github.com/google/uuid"
)

type node struct {
	storeClient *StoreClient
}

type HashRing struct {
	nodes             map[uint64]*node
	sortedKeys        []uint64
	replicationFactor int
}

func NewHashRing(rf int) *HashRing {
	return &HashRing{
		nodes:             make(map[uint64]*node),
		sortedKeys:        make([]uint64, 0),
		replicationFactor: rf,
	}
}

// adds a store to the ring with count nodes
func (hr *HashRing) AddStoreNodes(s *StoreClient) {

	for i := 0; i < hr.replicationFactor; i++ {

		// identify the position of the node on the ring
		hash := hashKey(s.name + "-" + uuid.New().String())

		hr.nodes[hash] = &node{
			storeClient: s,
		}

		// add the node to the sorted keys
		// high cost operation at insertion is fine
		// since this is done very few times compared to reads
		hr.sortedKeys = insertSorted(hr.sortedKeys, hash)
		s.nodeKeys = append(s.nodeKeys, hash)
	}
}

// finds the store for the given key
// this operation is O(log n), n = number of nodes in the ring
func (hr *HashRing) GetStore(key string) (*StoreClient, error) {
	hash := hashKey(key)

	if len(hr.nodes) == 0 {
		return nil, errors.New("no stores in the ring")
	}

	var s *StoreClient

	// find upper bound of hash in sortedKeys, if not found return the first element
	index := sort.Search(len(hr.sortedKeys), func(i int) bool {
		return hr.sortedKeys[i] >= hash
	})

	if index < len(hr.sortedKeys) {
		s = hr.nodes[hr.sortedKeys[index]].storeClient
	} else {
		s = hr.nodes[hr.sortedKeys[0]].storeClient
	}

	return s, nil
}
