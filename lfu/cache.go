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

//Put will put a cache key to the list.
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

//Get will retrieve a cache key from the list.
func (cache *Cache) Get(key string) (interface{}, error) {
	cache.lock.RLock()
	defer cache.lock.RUnlock()

	item, exists := cache.items[key]

	if !exists {
		return nil, errors.New("key not exists")
	}

	tracking, trackExists := cache.tracking[key]
	if !trackExists {
		tracking = Tracking{hits: 0}
	}
	tracking.hits = tracking.hits + 1
	cache.tracking[key] = tracking

	return item, nil
}

//Forget will delete a cache key from the list.
func (cache *Cache) Forget(key string) error {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	delete(cache.tracking, key)

	_, exists := cache.items[key]

	//true means: we do not need to delete, but is the same as when deletion was successful
	if !exists {
		//triggering the garbage collector to reduce the heap size
		runtime.GC()
		return nil
	}

	delete(cache.items, key)
	//triggering the garbage collector to reduce the heap size
	runtime.GC()

	cache.calcHeap()

	return nil
}

//calcHeap will store the HeapAlloc bytes at cache.heap
func (cache *Cache) calcHeap() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	cache.heap = m.HeapAlloc
}
