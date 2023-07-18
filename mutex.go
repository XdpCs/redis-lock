package redislock

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

// Mutex is a distributed mutex implementation based on redis.
type Mutex struct {
	client        *Client
	key           string
	value         string
	expiration    time.Duration
	retryStrategy RetryStrategy
	watchDog      *WatchDog
}

// Unlock releases the lock.
func (m *Mutex) Unlock(ctx context.Context) error {
	defer func() {
		// stop watch dog
		if m.watchDog != nil {
			if m.watchDog.cancelFunc != nil {
				m.watchDog.cancelFunc()
			}
		}
	}()

	if m == nil {
		return ErrMutexNotHeld
	}

	status, err := luaUnlock.Run(ctx, m.client.redisClient, []string{m.key}, m.value).Result()
	if err == redis.Nil {
		return ErrMutexNotHeld
	} else if err != nil {
		return err
	}

	if status != 1 {
		return ErrMutexNotHeld
	}
	return nil
}

func (m *Mutex) Refresh(ctx context.Context) error {
	if m == nil {
		return ErrMutexNotHeld
	}

	status, err := luaRefresh.Run(ctx, m.client.redisClient, []string{m.key}, m.value, m.expiration.Milliseconds()).Result()
	if err != nil {
		return err
	}

	if status == 1 {
		return nil
	}
	return nil
}

func (m *Mutex) runWatchDog(ctx context.Context) {
	for !atomic.CompareAndSwapUint32(&m.watchDog.isStart, 0, 1) {
	}

	ctx, m.watchDog.cancelFunc = context.WithCancel(ctx)
	go func() {
		ticker := time.NewTicker(m.watchDog.expiration / 2)
		defer ticker.Stop()

		for range ticker.C {
			select {
			case <-ctx.Done():
				return
			default:
			}

			if err := m.Refresh(ctx); err != nil {
				m.watchDog.cancelFunc()
				return
			}
		}
	}()
}

func (m *Mutex) stopWatchDog() {
	if m.watchDog.cancelFunc != nil {
		m.watchDog.cancelFunc()
	}
}

// newMutex creates a new Mutex.
func newMutex(client *Client, key, value string, strategy RetryStrategy) *Mutex {
	return &Mutex{
		client:        client,
		key:           key,
		value:         value,
		retryStrategy: strategy,
	}
}

func (m *Mutex) setWatchDog(watchDog *WatchDog) {
	m.watchDog = watchDog
}
