package redis

import (
	"time"

	rds "github.com/go-redis/redis"
)

// Redis wrapper.
type Redis struct {
	Address  string
	Password string
}

func (r Redis) newClient() *rds.Client {
	return rds.NewClient(&rds.Options{
		Addr:     r.Address,
		Password: r.Password,
	})
}

// Get returns a value associated with the provided key.
func (r Redis) Get(key string) ([]byte, error) {
	client := r.newClient()
	return client.Get(key).Bytes()
}

// Set sets a value for the provided key.
func (r Redis) Set(key string, value interface{}, expiration time.Duration) error {
	client := r.newClient()
	return client.Set(key, value, expiration).Err()
}
