package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/groupcache/lru"
)

func TestLoadableLRUCache(t *testing.T) {
	f := func(key interface{}) (interface{}, error) {
		time.Sleep(1 * time.Second)
		return fmt.Sprintf("Got a key: %s", key), nil
	}
	onEvicted := func(key lru.Key, value interface{}) {
		t.Logf("The key: %s has been evicted!\n", key)
	}
	c := NewLoadableLRUCacheWithOnEvicted(f, 2, onEvicted)

	go func() {
		v, _ := c.Get("key1")
		t.Logf("result in goroutin1: %s", v)
	}()

	go func() {
		v, _ := c.Get("key2")
		t.Logf("result in goroutin2: %s", v)
	}()

	v, _ := c.Get("key3")
	t.Logf("result in main goroutine: %s", v)
	time.Sleep(2 * time.Second)
}
