package lfu

import (
	"errors"
	"sync"
)

type Cache struct {
	size     uint
	items    map[string]interface{}
	tracking map[string]uint
	lock     sync.RWMutex
}

var LFU Cache

func Init(size uint) Cache {
	return Cache{size: size}
}

func (cache *Cache) Put(key string, item interface{}) (bool, error) {
	if cache.items == nil {
		cache.items = make(map[string]interface{})
	}

	cache.items[key] = item

	return true, nil
}

func (cache *Cache) Get(key string) (interface{}, error) {
	item, exists := cache.items[key]

	if !exists {
		return nil, errors.New("key not exists")
	}

	return item, nil
}

func (cache *Cache) Forget(key string) error {
	_, exists := cache.items[key]

	//true means: we do not need to delete, but is the same as when deletion was successful
	if !exists {
		return nil
	}

	delete(cache.items, key)

	return nil
}
