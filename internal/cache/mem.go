package cache

import (
	"sync"
	"time"
)

type InMemoryCache struct {
	data sync.Map
}

func NewInMemoryCache() Cache {
	return &InMemoryCache{}
}

func (c *InMemoryCache) Set(key string, value interface{}, expiration time.Duration) {
	c.data.Store(key, value)
	if expiration > 0 {
		go func() {
			time.Sleep(expiration)
			c.data.Delete(key)
		}()
	}
}

func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	return c.data.Load(key)
}

func (c *InMemoryCache) Delete(key string) {
	c.data.Delete(key)
}
