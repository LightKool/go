package cache

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/golang/groupcache/lru"
)

// Entry status constants
const (
	invalid int32 = iota
	loading
	active
)

type entry struct {
	key    interface{}
	value  interface{}
	status int32
	signal chan bool
	err    error
}

func newEntry(key interface{}) *entry {
	return &entry{key: key, signal: make(chan bool)}
}

func (e *entry) getOrLoad(c *LoadableLRUCache) (interface{}, error) {
	switch atomic.LoadInt32(&e.status) {
	case invalid:
		if atomic.CompareAndSwapInt32(&e.status, invalid, loading) {
			value, err := c.loader(e.key)
			// nil value is not allowed because it will cause the
			// data inconsistency with the persistence layer.
			if value == nil && err == nil {
				err = fmt.Errorf("can't load value of key: [%s]", e.key)
			}
			e.value = value
			e.err = err
			if err == nil {
				atomic.CompareAndSwapInt32(&e.status, loading, active)
			} else {
				c.Remove(e.key)
			}
			// close the signal chan to wake up all blocking goroutines
			close(e.signal)
			return value, err
		}
		<-e.signal
		return e.value, e.err
	case loading:
		<-e.signal
		return e.value, e.err
	default:
		return e.value, e.err
	}
}

// Loader is responsible for loading the cache value if it
// hasn't been cached yet.
type Loader func(interface{}) (interface{}, error)

// LoadableLRUCache is an LRUCache which will automatically
// load the value with provided Loader if the value hasn't been
// cached yet.
type LoadableLRUCache struct {
	mu     sync.RWMutex
	data   *lru.Cache // Use LRU Cache to store the entries.
	loader Loader
}

// Get returns the value associated with the key if it exists,
// otherwise try to load/cache/return the value through the loading
// function provided when creating the instance of the implementation.
// The function will return the value if it's loaded successfully,
// or the nil value and an error if the value couldn't be loaded.
func (c *LoadableLRUCache) Get(key interface{}) (interface{}, error) {
	e := c.getEntry(key)
	return e.getOrLoad(c)
}

// Remove deletes the cached value associated with the key from the cache.
func (c *LoadableLRUCache) Remove(key interface{}) {
	c.mu.Lock()
	c.data.Remove(key)
	c.mu.Unlock()
}

// RemoveOldest deletes the least-recently-used(LRU) value from the cache.
func (c *LoadableLRUCache) RemoveOldest() {
	c.mu.Lock()
	c.data.RemoveOldest()
	c.mu.Unlock()
}

// RemoveOldestN deletes the n least-recently-used values from the cache.
func (c *LoadableLRUCache) RemoveOldestN(n int) {
	c.mu.Lock()
	for i := 0; i < n; i++ {
		c.data.RemoveOldest()
	}
	c.mu.Unlock()
}

// Clear deletes all values from the cache.
func (c *LoadableLRUCache) Clear() {
	c.mu.Lock()
	c.data.Clear()
	c.mu.Unlock()
}

// Retrieve or create entry instance. Use double reference check
// to handle concurrent access.
func (c *LoadableLRUCache) getEntry(key interface{}) *entry {
	c.mu.RLock()
	e, ok := c.data.Get(key)
	if !ok {
		c.mu.RUnlock()
		c.mu.Lock()
		e, ok = c.data.Get(key)
		if !ok {
			e = newEntry(key)
			c.data.Add(key, e)
		}
		c.mu.Unlock()
	} else {
		c.mu.RUnlock()
	}
	return e.(*entry)
}

// NewLoadableLRUCache creates a new LoadableLRUCache.
// The loader mustn't be nil.
func NewLoadableLRUCache(loader Loader, maxEntries int) *LoadableLRUCache {
	return &LoadableLRUCache{loader: loader, data: lru.New(maxEntries)}
}

// NewLoadableLRUCacheWithOnEvicted creates a new LoadableLRUCache.
// The loader mustn't be nil.
//
// The onEvicted callback will be called after the key and value
// being deleted, manually or automatically.
func NewLoadableLRUCacheWithOnEvicted(loader Loader, maxEntries int, onEvicted func(lru.Key, interface{})) *LoadableLRUCache {
	return &LoadableLRUCache{loader: loader, data: &lru.Cache{MaxEntries: maxEntries, OnEvicted: onEntryEvicted(onEvicted)}}
}

func onEntryEvicted(onEvicted func(lru.Key, interface{})) func(lru.Key, interface{}) {
	return func(key lru.Key, value interface{}) {
		e := value.(*entry)
		if atomic.LoadInt32(&e.status) == active {
			onEvicted(key, e.value)
		}
	}
}
