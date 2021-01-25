package cache

import (
	"github.com/shoulai/go-auth/cache/memory"
	"github.com/shoulai/go-auth/cache/redis"
	"github.com/shoulai/go-auth/config"
)

type IStore interface {
	Set(key string, val interface{}, expiration int) error
	Get(key string, rel interface{}) error
	Del(key string) error
	Close() error
}

func NewStore(cfg config.Config) (IStore, func(), error) {
	var store IStore
	switch cfg.CacheName {
	case "redis":
		store_, err := redis.NewStore(cfg.Redis)
		if err != nil {
			return nil, nil, err
		}
		store = store_
		break
	default:
		store = memory.NewStore()
		break
	}
	closeStore := func() {
		store.Close()
	}
	return store, closeStore, nil
}
