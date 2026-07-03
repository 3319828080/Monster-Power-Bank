package data

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// RedisLock implements a simple distributed lock via Redis SETNX.
type RedisLock struct {
	rdb    *redis.Client
	prefix string
}

func NewRedisLock(rdb *redis.Client) *RedisLock {
	return &RedisLock{rdb: rdb, prefix: "lock:"}
}

func (l *RedisLock) Lock(ctx context.Context, key string, ttl time.Duration) (string, bool, error) {
	token := uuid.New().String()
	ok, err := l.rdb.SetNX(ctx, l.prefix+key, token, ttl).Result()
	if err != nil {
		return "", false, err
	}
	return token, ok, nil
}

func (l *RedisLock) Unlock(ctx context.Context, key, token string) error {
	script := `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	end
	return 0
	`
	return l.rdb.Eval(ctx, script, []string{l.prefix + key}, token).Err()
}
