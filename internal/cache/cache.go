package cache

import "time"

type Cache interface {
	Set(key string, value interface{}, expiration time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}
