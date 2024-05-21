package utils

import (
	"fmt"
	"sync"
)

type cache struct {
	mutex   sync.Mutex
	storage map[string]interface{}
}

var Cache = cache{
	mutex: sync.Mutex{},
}

// set cache value by key, return bool result whether value was set or not.
func (c *cache) Set(key string, value interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.storage != nil {
		c.storage[key] = value
		// c.exportCache()
		return nil
	}

	return fmt.Errorf("cache storage is nil")
}

// get cache value by key, return false second param if no value.
func (c *cache) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	val, ok := c.storage[key]
	// c.exportCache()
	return val, ok
}

func (c *cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.storage, key)
}

func (c *cache) Init() {
	c.storage = make(map[string]interface{})
}
