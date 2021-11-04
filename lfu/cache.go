package lfu

import (
	"errors"
	"runtime"
	"sort"
	"sync"
	"time"
)

type Tracking struct {
	heap          uint64
	hits          uint
	lastUpdatedAt time.Time
}

type Cache struct {
	size        uint64
	items       map[string]interface{}
	tracking    map[string]Tracking
	lock        sync.RWMutex
	heap        uint64
	heapRuntime uint64
}

var LFU Cache

// Init will init a new cache. Define the max size of the cache in bytes.
// Minimum is 500000 (0.5 MB), but this only a theoretic one and not recommended.
// Minimum is required to can store some elements, there are basic heap allocations caused by maps.
func Init(size uint64) Cache {
	if size < 500000 {
		panic("minimum size is 500000")
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return Cache{size: size, heapRuntime: m.HeapAlloc}
}

// Put will put a cache key to the list.
func (cache *Cache) Put(key string, item interface{}) (bool, error) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if cache.items == nil {
		cache.items = make(map[string]interface{})
	}

	if cache.tracking == nil {
		cache.tracking = map[string]Tracking{}
	}

	actualHeap := cache.heapRuntime

	cache.items[key] = item

	cache.calcHeapRuntime()

	itemHeap := cache.heapRuntime - actualHeap

	//@TODO: think about it - what to do when item exists and a put was made again? Really create new Tracking? Update Tracking?
	cache.tracking[key] = Tracking{heap: itemHeap, lastUpdatedAt: time.Now()}

	cache.heap = cache.heap + itemHeap

	//when cache is oversized, we must reduce it - but we will not delete the new putted key!
	if cache.heap > cache.size {
		cache.reduce(cache.size, key)
	}

	return true, nil
}

// Get will retrieve a cache key from the list.
func (cache *Cache) Get(key string) (interface{}, error) {
	cache.lock.RLock()
	defer cache.lock.RUnlock()

	item, exists := cache.items[key]

	if !exists {
		return nil, errors.New("key not exists")
	}

	tracking, trackExists := cache.tracking[key]
	if !trackExists {
		tracking = Tracking{hits: 0, lastUpdatedAt: time.Now()}
	}
	tracking.hits = tracking.hits + 1
	cache.tracking[key] = tracking

	return item, nil
}

// Forget will delete a cache key from the list.
func (cache *Cache) Forget(key string) error {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	//actualize the heap size
	if tracking, exists := cache.tracking[key]; exists {
		cache.heap = cache.heap - tracking.heap
		delete(cache.tracking, key)
	}

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

	cache.calcHeapRuntime()

	return nil
}

// Reduce will reduce the map to a max heap size.
// Our first algorithm here is the LFU one.
// @TODO: later we can try to optimize it; there will be a few more possibilities (integrating when was the last hit and when was the creation date)
func (cache *Cache) reduce(max uint64, ignore string) (bool, error) {
	//first we will build a slice
	keys := cache.getSortedTrackingKeys(ignore)

	for _, key := range keys {
		cache.Forget(key)

		if cache.heap <= max {
			break
		}
	}

	return true, nil
}

// getSortedTrackingKeys will sort desc by hits
func (cache *Cache) getSortedTrackingKeys(ignore string) []string {
	//first we will build a slice
	keys := make([]string, 0, len(cache.items))

	for key := range cache.tracking {
		if key == ignore {
			continue
		}
		keys = append(keys, key)
	}

	//next we must order this slice by the hits of its keys
	sort.Slice(keys, func(i, j int) bool {
		return cache.tracking[keys[i]].hits > cache.tracking[keys[j]].hits
	})

	return keys
}

// calcHeap will store the HeapAlloc bytes at cache.heap
func (cache *Cache) calcHeapRuntime() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	cache.heapRuntime = m.HeapAlloc
}
