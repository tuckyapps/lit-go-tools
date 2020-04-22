package datasource

import "time"

// InMemoryDB declares Set and Get operations for redis.
type InMemoryDB interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) ([]byte, error)
}
