package redis

import (
	"sync"
	"time"

	rds "github.com/go-redis/redis"
)

var (
	instance     InMemoryDB
	syncInstance = new(sync.Mutex)
)

// InMemoryDB declares Set and Get operations for redis.
type InMemoryDB interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) ([]byte, error)
}

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

//BuildNewInMemoryConnection returns an in-memory database connection
func BuildNewInMemoryConnection(address string, password string) InMemoryDB {
	if instance == nil {

		syncInstance.Lock()
		defer syncInstance.Unlock()

		if instance == nil {
			instance = &Redis{
				Address:  address,
				Password: password,
			}
		}

	}
	return instance
}

// ResetInMemoryDb is used to erase instance.
func ResetInMemoryDb() {
	syncInstance.Lock()
	defer syncInstance.Unlock()
	instance = nil
}
