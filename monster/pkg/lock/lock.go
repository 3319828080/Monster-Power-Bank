package lock

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	defaultExpiration = 30 * time.Second
	defaultRetryCount = 3
	defaultRetryDelay = 200 * time.Millisecond
	renewalInterval   = 10 * time.Second
)

type Lock struct {
	client     *redis.Client
	key        string
	value      string
	expiration time.Duration
	mu         sync.Mutex
	renewalCh  chan struct{}
}

func NewLock(client *redis.Client, key string, opts ...Option) *Lock {
	l := &Lock{
		client:     client,
		key:        key,
		value:      uuid.New().String(),
		expiration: defaultExpiration,
		renewalCh:  make(chan struct{}, 1),
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

type Option func(*Lock)

func WithExpiration(exp time.Duration) Option {
	return func(l *Lock) { l.expiration = exp }
}

func (l *Lock) Lock(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	var ok bool
	for i := 0; i < defaultRetryCount; i++ {
		ok, _ = l.client.SetNX(ctx, l.key, l.value, l.expiration).Result()
		if ok {
			go l.renewal(ctx)
			return nil
		}
		time.Sleep(defaultRetryDelay)
	}
	return ErrLockAcquireFailed
}

func (l *Lock) Unlock(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	close(l.renewalCh)

	script := `if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("DEL", KEYS[1]) else return 0 end`
	return l.client.Eval(ctx, script, []string{l.key}, l.value).Err()
}

func (l *Lock) renewal(ctx context.Context) {
	ticker := time.NewTicker(renewalInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			script := `if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("PEXPIRE", KEYS[1], ARGV[2]) else return 0 end`
			l.client.Eval(ctx, script, []string{l.key}, l.value, l.expiration.Milliseconds()).Err()
		case <-l.renewalCh:
			return
		case <-ctx.Done():
			return
		}
	}
}
