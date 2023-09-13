package store

import (
	"container/list"
	"errors"
)

type LRUCache struct {
	capacity uint32
	cache    map[string]*list.Element
	list     *list.List
}

type Pair struct {
	key   string
	value string
}

func LRUConstructor(capacity uint32) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (lr *LRUCache) Get(key string) (string, error) {
	if elem, ok := lr.cache[key]; ok {
		lr.list.MoveToFront(elem)
		return elem.Value.(Pair).value, nil
	}
	return "", errors.New("cache miss")
}

func (lr *LRUCache) Put(key string, value string) {
	if elem, ok := lr.cache[key]; ok {
		elem.Value = Pair{key, value}
		lr.list.MoveToFront(elem)
	} else {
		if uint32(len(lr.cache)) >= lr.capacity {
			// Remove the least recently used item
			tail := lr.list.Back()
			if tail != nil {
				delete(lr.cache, tail.Value.(Pair).key)
				lr.list.Remove(tail)
			}
		}

		// Add the new item to the front of the list
		newElem := lr.list.PushFront(Pair{key, value})
		lr.cache[key] = newElem
	}
}

func (lr *LRUCache) Remove(key string) {
	if elem, ok := lr.cache[key]; ok {
		delete(lr.cache, key)
		lr.list.Remove(elem)
	}
}
