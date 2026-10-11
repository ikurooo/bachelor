package main

import "sync/atomic"

type Cache struct {
	store atomic.Value
}

func (c *Cache) Get() map[string]MachineInfo {
	val := c.store.Load()
	if val == nil {
		return nil
	}

	m, ok := val.(map[string]MachineInfo)
	if !ok {
		return nil
	}

	return m
}

func (c *Cache) Update(newData map[string]MachineInfo) {
	c.store.Store(newData)
}
