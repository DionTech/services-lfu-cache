package lfu

import (
	"errors"
	"runtime"
	"sync"
)

type Tracking struct {
	heap uint64
	hits uint
}

type Cache struct {
	size     uint64
	items    map[string]interface{}
	tracking map[string]Tracking
	lock     sync.RWMutex
	heap     uint64
}

var LFU Cache

//Init will init a new cache. Define the max size of the cache in bytes
func Init(size uint64) Cache {
	return Cache{size: size}
}

func (cache *Cache) Put(key string, item interface{}) (bool, error) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if cache.items == nil {
		cache.items = make(map[string]interface{})
	}

	if cache.tracking == nil {
		cache.tracking = map[string]Tracking{}
	}

	actualHeap := cache.heap

	cache.items[key] = item

	cache.calcHeap()

	itemHeap := cache.heap - actualHeap

	cache.tracking[key] = Tracking{heap: itemHeap}

	return true, nil
}

func (cache *Cache) Get(key string) (interface{}, error) {
	cache.lock.RLock()
	defer cache.lock.RUnlock()

	item, exists := cache.items[key]

	if !exists {
		return nil, errors.New("key not exists")
	}

	return item, nil
}

func (cache *Cache) Forget(key string) error {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	delete(cache.tracking, key)

	_, exists := cache.items[key]

	//true means: we do not need to delete, but is the same as when deletion was successful
	if !exists {
		return nil
	}

	delete(cache.items, key)
	//triggering the garbage collector to reduce the heap size
	runtime.GC()

	cache.calcHeap()

	return nil
}

func (cache *Cache) calcHeap() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	cache.heap = m.HeapAlloc
}
