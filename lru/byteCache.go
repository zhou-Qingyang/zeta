package main

import "sync"

type ByteCache struct {
	mu         sync.Mutex
	lru        *LruCache
	cacheBytes int64
}

func (c *ByteCache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = NewCache(c.cacheBytes)
	}
	c.lru.Insert(key, value)
}

func (c *ByteCache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
