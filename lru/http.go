package main

import (
	"log"
	"net/http"
	"strings"
)

type HttpPool struct {
	baseUrl string
	cache   *LruCache
}

func NewHttpPool(baseUrl string, cache *LruCache) *HttpPool {
	return &HttpPool{
		baseUrl: baseUrl,
		cache:   cache,
	}
}

//func (p *HttpPool) Log(prefix string, format string, v ...interface{}) {
//	log.Printf("[Server %s] %s", prefix, fmt.Sprintf(format, v...))
//}

func (h *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, h.baseUrl) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	parts := strings.SplitN(r.URL.Path[len(h.baseUrl):], "/", 1)
	key := parts[0]
	key = strings.Replace(key, "/", "", -1)
	log.Printf("[Server %s] %s", r.URL.Path, key)
	w.Header().Set("Content-Type", "octet-stream")
	data, ok := h.cache.Get(key)
	if ok {
		w.Write([]byte(data.(String)))
		return
	} else {
		//h.Log("%s %s",  r.URL.Path, key)
		log.Printf("cache miss for key %s", r.URL.Path)
	}
	w.Write([]byte("not found"))
}

type String string

func (s String) Len() int {
	return len(s)
}

func main() {
	value1, value2, value3 := String("value1"), String("value2"), String("value3")
	//value4 := String("value4")
	cacheCap := value1.Len() + value2.Len() + value3.Len()
	lru := NewCache(int64(cacheCap))
	lru.Insert("key1", value1)
	lru.Insert("key2", value2)
	httpPool := NewHttpPool("/lru", lru)
	http.ListenAndServe("localhost:9999", httpPool)
}
