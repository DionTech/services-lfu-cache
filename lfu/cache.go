package lfu

import "sync"

type Cache struct {
	size     uint
	items    map[string]interface{}
	tracking map[string]uint
	lock     sync.RWMutex
}

func (cache *Cache) Put(key string, item interface{}) (bool, error) {

}

func (cache *Cache) Get(key string) (interface{}, error) {

}

func (cache *Cache) Forget(key string) (bool, error) {

}
