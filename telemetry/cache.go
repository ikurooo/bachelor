package main

import (
	"sync/atomic"
)

// CacheStore holds thread-safe cluster topology
type CacheStore struct {
	store atomic.Value // Stores Topology
}

func NewCacheStore() *CacheStore {
	c := &CacheStore{}
	c.store.Store(Topology{
		Nodes:           make(map[string]Node),
		WarmPods:        make(map[string][]ActivePod),
		CompatibleNodes: make(map[string][]string),
	})
	return c
}

func (c *CacheStore) Get() Topology {
	return c.store.Load().(Topology)
}

func (c *CacheStore) Update(newTopo Topology) {
	c.store.Store(newTopo)
}
