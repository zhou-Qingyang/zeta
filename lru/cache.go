package main

import (
	"container/list"
)

type LruCache struct {
	maxBytes int64
	nBytes   int64 // 是当前已使用的内存
	ll       *list.List
	cache    map[string]*list.Element
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

func NewCache(maxBytes int64) *LruCache {
	return &LruCache{
		maxBytes: maxBytes,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
	}
}

func (c *LruCache) Get(key string) (Value, bool) {
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToBack(e) // 将元素移动到队尾
		return e.Value.(*entry).value, true
	}
	return nil, false
}

func (c *LruCache) Remove(key string) bool {
	if e, ok := c.cache[key]; ok {
		c.ll.Remove(e)
		kz := e.Value.(*entry)
		valSize := int64(kz.value.Len())
		c.nBytes -= valSize
		delete(c.cache, key)
		return true
	}
	return false
}

func (c *LruCache) RemoveOldest() {
	ele := c.ll.Front()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(kv.value.Len())
	}
}

func (c *LruCache) Insert(key string, value Value) bool {
	// 先判断缓存中是否存在
	valSize := int64(value.Len())
	if valSize > c.maxBytes {
		return false // 单条数据过大，直接拒绝
	}

	if e, ok := c.cache[key]; ok {
		c.ll.MoveToBack(e)
		return true
	}
	// 确保有足够空间（可能需要循环移除多个旧元素）
	for c.nBytes+valSize > c.maxBytes && c.nBytes > 0 {
		c.RemoveOldest() // 内部已更新 nBytes
	}
	e := c.ll.PushBack(&entry{key, value}) // 通常新元素放头部
	c.cache[key] = e
	c.nBytes += valSize
	return true
}

func (c *LruCache) Len() int {
	return c.ll.Len()
}

func (c *LruCache) PrintElement() {
	for e := c.ll.Front(); e != nil; e = e.Next() {
		kv := e.Value.(*entry)
		print(kv.key, " ", kv.value.Len(), " ", c.nBytes, "\n")
	}
}
