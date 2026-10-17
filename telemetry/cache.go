package main

import "sync/atomic"

type CacheStore struct {
	store atomic.Value
}

func NewCacheStore() *CacheStore {
	c := &CacheStore{}
	c.store.Store(make(map[string][]Node))
	return c
}

func (c *CacheStore) Update(machines []Node) {
	newMap := make(map[string][]Node)
	//for _, m := range machines {
	//	newMap[m.RuntimeInfo] = append(newMap[m], m)
	//}
	c.store.Store(newMap)
}
