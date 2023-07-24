package redislock

import (
	"context"
	"crypto/rc4"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient is the interface used by redislock to interact with redis.
type RedisClient interface {
	redis.Scripter
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
}

// Client is the redislock client, wraps RedisClient.
type Client struct {
	redisClient RedisClient
	cipherKey   string // if cipherKey is "-1", use unix timestamp as Cipher cipherKey.
	*rc4.Cipher        // customize cipher, default is rc4.NewCipher with cipherKey.
}

// NewClient creates a new redislock client.
func NewClient(redisClient RedisClient, options ...ClientOption) (*Client, error) {
	c := &Client{redisClient: redisClient, cipherKey: "-1"}

	for _, option := range options {
		option(c)
	}

	if c.Cipher == nil {
		key := c.cipherKey
		if c.cipherKey == "-1" {
			key = time.Now().String()
		}

		cipher, err := rc4.NewCipher([]byte(key))
		if err != nil {
			return nil, fmt.Errorf("rc4.NewCipher error: %w", err)
		}
		c.Cipher = cipher
	}

	return c, nil
}

// NewDefaultClient creates a new default redislock client.
func NewDefaultClient(redisClient RedisClient) (*Client, error) {
	return NewClient(redisClient, WithCipherKey("1118"))
}

type ClientOption func(client *Client)

// WithCipherKey sets the cipherKey of the client.
func WithCipherKey(cipherKey string) ClientOption {
	return func(client *Client) {
		client.cipherKey = cipherKey
	}
}

// WithCipher sets the cipher of the client,
// if you set WithCipherKey and WithCipher at the same time,
// WithCipherKey will be ignored.
func WithCipher(cipher *rc4.Cipher) ClientOption {
	return func(client *Client) {
		client.Cipher = cipher
	}
}

// TryLock tries to acquire a lock with default parameter.
func (c *Client) TryLock(ctx context.Context, key string, expiration time.Duration) (*Mutex, error) {
	option := &mutexOption{}

	// expiration == -1 means no expiration, so start watch dog.
	if expiration == -1 {
		option.watchDog = NewDefaultWatchDog()
		expiration = option.watchDog.expiration
	}
	option.retryStrategy = NewNoRetry()

	return c.tryLock(ctx, key, expiration, option)
}

// TryLockWithRetry tries to acquire a lock with retry strategy.
func (c *Client) TryLockWithRetry(ctx context.Context, key string, expiration time.Duration, retryStrategy RetryStrategy) (*Mutex, error) {
	option := &mutexOption{}

	// expiration == -1 means no expiration, so start watch dog.
	if expiration == -1 {
		option.watchDog = NewDefaultWatchDog()
		expiration = option.watchDog.expiration
	}
	option.retryStrategy = retryStrategy

	return c.tryLock(ctx, key, expiration, option)
}

// TryLockWithWatchDog tries to acquire a lock with watch dog.
func (c *Client) TryLockWithWatchDog(ctx context.Context, key string, watchDog *WatchDog) (*Mutex, error) {
	option := &mutexOption{}
	var err error

	watchDog, err = checkWatchDogReturnWatchDog(watchDog)
	if err != nil {
		return nil, fmt.Errorf("checkWatchDogReturnWatchDog error: %w", err)
	}

	option.watchDog = watchDog
	option.retryStrategy = NewNoRetry()

	return c.tryLock(ctx, key, watchDog.expiration, option)
}

// TryLockWithRetryAndWatchDog tries to acquire a lock with retry strategy and watch dog.
func (c *Client) TryLockWithRetryAndWatchDog(ctx context.Context, key string, retryStrategy RetryStrategy, watchDog *WatchDog) (*Mutex, error) {
	option := &mutexOption{}
	var err error

	watchDog, err = checkWatchDogReturnWatchDog(watchDog)
	if err != nil {
		return nil, fmt.Errorf("checkWatchDogReturnWatchDog error: %w", err)
	}

	option.watchDog = watchDog
	option.retryStrategy = retryStrategy

	return c.tryLock(ctx, key, watchDog.expiration, option)
}

func (c *Client) tryLock(ctx context.Context, key string, expiration time.Duration, option *mutexOption) (*Mutex, error) {
	value, err := c.getValue()
	if err != nil {
		return nil, fmt.Errorf("c.getValue error: %w", err)
	}

	parentCtx := ctx
	childCtx := ctx

	mutex := newMutex(c, key, value, expiration, option.retryStrategy)

	if option.watchDog != nil {
		mutex.setWatchDog(option.watchDog)
	}

	if _, ok := childCtx.Deadline(); !ok {
		var cancel context.CancelFunc
		childCtx, cancel = context.WithDeadline(childCtx, time.Now().Add(expiration))
		defer cancel()
	}

	var ticker *time.Ticker
	for {
		ok, err := c.lock(childCtx, key, value, expiration)
		if err != nil {
			return nil, fmt.Errorf("c.lock error: %w", err)
		}

		if ok {
			if option.watchDog != nil {
				mutex.runWatchDog(parentCtx)
			}
			return mutex, nil
		}

		retryTime := option.retryStrategy.NextRetryTime()
		if retryTime == 0 {
			return nil, ErrMutexLockFailed
		}

		if ticker == nil {
			ticker = time.NewTicker(retryTime)
			defer ticker.Stop()
		} else {
			ticker.Reset(retryTime)
		}

		select {
		case <-childCtx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
	}
}

func (c *Client) lock(ctx context.Context, key, value string, expiration time.Duration) (bool, error) {
	return c.redisClient.SetNX(ctx, key, value, expiration).Result()
}

// getValue returns a value that is unique to this client.
func (c *Client) getValue() (string, error) {
	cipher := c.Cipher
	nowString := time.Now().String()
	value := make([]byte, len(nowString))
	cipher.XORKeyStream(value, []byte(nowString))

	return string(value), nil
}
