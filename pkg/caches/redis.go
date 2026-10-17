package caches

import (
	"chanlearning/pkg/service/common/envvar"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	mu  sync.Mutex
	rd  *redis.Client
	err error
)

// singleton pattern
func Cacheinit() *redis.Client {
	if rd == nil {
		mu.Lock()
		defer mu.Unlock()
		if rd == nil {
			url := os.Getenv(envvar.Redis)
			opts, err := redis.ParseURL(url)
			if err != nil {
				panic(err)
			}
			rd = redis.NewClient(opts)
			return rd
		} else {
			return rd
		}
	} else {
		return rd
	}
}
