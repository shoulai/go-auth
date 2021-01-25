package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/shoulai/go-auth/config"
	"time"
)

type Store struct {
	client redisClient
	prefix string
}

type redisClient interface {
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Expire(key string, expiration time.Duration) *redis.BoolCmd
	Exists(keys ...string) *redis.IntCmd
	TxPipeline() redis.Pipeliner
	Del(keys ...string) *redis.IntCmd
	Close() error
}

// NewStore 创建基于redis存储实例
func NewStore(cfg *config.Redis) (*Store, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		DB:       cfg.DB,
		Password: cfg.Password,
	})
	cmdPing := cli.Ping()
	err := cmdPing.Err()
	if err != nil {
		return nil, err
	}
	return &Store{
		client: cli,
		prefix: cfg.KeyPrefix,
	}, nil
}

func (s *Store) Set(key string, val interface{}, expiration int) error {
	key = fmt.Sprintf(s.prefix+"_%v", key)
	cmd := s.client.Set(key, val, time.Duration(expiration)*time.Second)
	return cmd.Err()
}

// Delete ...
func (s *Store) Del(key string) error {
	key = fmt.Sprintf(s.prefix+"_%v", key)
	cmd := s.client.Del(key)
	if err := cmd.Err(); err != nil {
		return err
	}
	return nil
}

func (s *Store) Get(key string, rel interface{}) error {
	key = fmt.Sprintf(s.prefix+"_%v", key)
	cmd := s.client.Get(key)
	return cmd.Scan(rel)
}

func (s *Store) Close() error {
	return s.client.Close()
}
