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
		// maps the node name to (node name, available runtimes, active pods)
		Nodes: make(map[string]Node),
		// maps the runtime names to (pod name, pod ip, runtime)
		WarmPods: make(map[string][]ActivePod),
		// maps the runtime names to node names
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
