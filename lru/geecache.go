package main

import (
	"fmt"
	"github.com/zhou-Qingyang/hrm/ZetaCache/lru/pb"
	"github.com/zhou-Qingyang/hrm/ZetaCache/lru/singleflight"
	"log"
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
	peers     PeerPicker
	loader    *singleflight.Group
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
		loader: &singleflight.Group{},
	}
	groups[name] = group
	return group
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (group *Group) Get(key string) (value ByteView, err error) {
	//if key == "" {
	//	return ByteView{}, nil
	//}
	//if val, ok := group.mainCache.get(key); ok {
	//	fmt.Println("cache hit")
	//	return val, nil
	//}
	//return group.loadLocally(key)
	viewi, err := group.loader.Do(key, func() (interface{}, error) {
		if group.peers != nil {
			if peer, ok := group.peers.PickPeer(key); ok {
				if value, err := group.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer")
			}
		}
		return group.loadLocally(key)
	})
	if err == nil {
		return viewi.(ByteView), nil
	}
	return
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	//bytes, err := peer.Get(g.name, key)
	//if err != nil {
	//	return ByteView{}, err
	//}
	//return ByteView{b: bytes}, nil
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
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
