package main

import (
	"fmt"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function.
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	getter    Getter
	name      string
	mainCache ByteCache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, maxSize int64, getter Getter) *Group {
	if getter == nil {
		panic("getter is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	group := &Group{
		name:   name,
		getter: getter,
		mainCache: ByteCache{
			cacheBytes: maxSize,
		},
	}
	groups[name] = group
	return group
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

func (group *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, nil
	}
	if val, ok := group.mainCache.get(key); ok {
		fmt.Println("cache hit")
		return val, nil
	}
	return group.loadLocally(key)
}

func (group *Group) loadLocally(key string) (ByteView, error) {
	bytes, err := group.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	// 获取成功 丢到缓存中并返回
	b := ByteView{
		b: bytes,
	}
	group.mainCache.add(key, b)
	fmt.Println("loadLocally hit")
	return b, nil
}
