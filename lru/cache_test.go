package main

import (
	"fmt"
	"testing"
)

// type String string
//
//	func (s String) Len() int {
//		return len(s)
//	}
func TestLRUCacheGet(t *testing.T) {
	lru := NewCache(1000)
	lru.Insert("key1", String("value1"))
	if val, _ := lru.Get("key1"); string(val.(String)) != "value1" {
		fmt.Println("Error: value1 not found in LRU cache")
	}
	fmt.Println("LRU Cache Size:", lru.Len())
}

func TestLRUCacheRemoveOldest(t *testing.T) {
	value1, value2, value3 := String("value1"), String("value2"), String("value3")
	value4 := String("value4")
	cacheCap := value1.Len() + value2.Len() + value3.Len()
	lru := NewCache(int64(cacheCap))
	lru.Insert("key1", value1)
	lru.Insert("key1", value1)
	lru.Insert("key2", value2)
	lru.Insert("key3", value3)
	lru.Insert("key2", value2)
	lru.Insert("key4", value4)
	lru.Insert("key5", String("value5"))
	//if val, ok := lru.Get("key1"); ok {
	//	//fmt.Println("Error: key1 should have been removed from LRU cache")
	//	fmt.Printf("Error: key1 should have been removed from LRU cache, but value is %s\n", string(val.(String)))
	//}
	fmt.Println("LRU Cache Size:", lru.Len())
	lru.PrintElement()
}
