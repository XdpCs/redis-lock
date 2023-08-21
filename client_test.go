package redislock

import (
	"context"
	"crypto/rc4"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestNewClient(t *testing.T) {
	// init redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})
	// close redis client
	defer rdb.Close()

	// init actualOne
	actualOne, err := NewClient(rdb)
	if err != nil {
		t.Fatalf("actualOne NewClient error:[%v]", err)
	}

	actualTwo, err := NewClient(rdb, WithCipherKey("11181114"))
	if err != nil {
		t.Fatalf("actualTwo NewClient error:[%v]", err)
	}

	// init cipherTwo
	cipherTwo, err := rc4.NewCipher([]byte("11181114"))
	if err != nil {
		t.Fatalf("init cipherTwo error:[%v]", err)
	}

	// init cipherThree
	cipherThree, err := rc4.NewCipher([]byte("1114"))
	if err != nil {
		t.Fatalf("init cipherThree error:[%v]", err)
	}

	actualThree, err := NewClient(rdb, WithCipher(cipherThree))
	if err != nil {
		t.Fatalf("actualThree NewClient error:[%v]", err)
	}

	// test cases
	cases := []struct {
		Name             string
		Actual, Expected *Client
	}{
		{
			"NewClientWithNothing",
			actualOne,
			&Client{redisClient: rdb, cipherKey: "-1", Cipher: actualOne.Cipher},
		},
		{
			"NewClientWithCipherKey",
			actualTwo,
			&Client{redisClient: rdb, cipherKey: "11181114", Cipher: cipherTwo},
		},
		{
			Name:     "NewClientWithCipher",
			Actual:   actualThree,
			Expected: &Client{redisClient: rdb, cipherKey: "-1", Cipher: cipherThree},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			compareClient(t, c.Expected, c.Actual)
		})
	}
}

func TestNewDefaultClient(t *testing.T) {
	// init redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})
	// close redis client
	defer rdb.Close()
	// init redislock client
	client, err := NewDefaultClient(rdb)
	if err != nil {
		t.Fatalf("NewDefaultClient error:[%v]", err)
	}

	cipherKey := "1118"
	cipher, err := rc4.NewCipher([]byte(cipherKey))
	if err != nil {
		t.Fatalf("init cipher error:[%v]", err)
	}

	// test cases
	compareClient(t, &Client{redisClient: rdb, cipherKey: cipherKey, Cipher: cipher}, client)
}

func TestClient_TryLock(t *testing.T) {
	// init redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})
	// close redis client
	defer rdb.Close()
	// init redislock client
	client, err := NewDefaultClient(rdb)
	if err != nil {
		t.Fatalf("NewDefaultClient error:[%v]", err)
	}
	keyOne := "testOne"
	keyTwo := "testTwo"
	defer teardown(t, rdb, []string{keyOne, keyTwo})

	ctxOne := context.Background()
	ctxTwo := context.Background()

	actualOne, err := client.TryLock(ctxOne, keyOne, -1)
	if err != nil {
		t.Fatalf("actualOne TryLock error:[%v]", err)
	}

	actualTwo, err := client.TryLock(ctxTwo, keyTwo, 10*time.Second)
	if err != nil {
		t.Fatalf("actualTwo TryLock error:[%v]", err)
	}

	// test cases
	cases := []struct {
		Name     string
		Actual   *Mutex
		Expected *Mutex
	}{
		{
			"TryLockWithWatchDog",
			actualOne,
			&Mutex{
				client:        client,
				key:           keyOne,
				expiration:    30 * time.Second,
				value:         actualOne.value,
				retryStrategy: NewNoRetry(),
				watchDog:      actualOne.watchDog,
			},
		},
		{
			"TryLockWithExpiration",
			actualTwo,
			&Mutex{
				client:        client,
				key:           keyTwo,
				expiration:    10 * time.Second,
				value:         actualTwo.value,
				retryStrategy: NewNoRetry(),
				watchDog:      nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			compareMutex(t, c.Expected, c.Actual)
		})
	}
}

func TestClient_TryLockWithRetryStrategy(t *testing.T) {
	// init redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})
	// close redis client
	defer rdb.Close()
	// init redislock client
	client, err := NewDefaultClient(rdb)
	if err != nil {
		t.Fatalf("NewDefaultClient error:[%v]", err)
	}
	keyOne := "testOne"
	keyTwo := "testTwo"
	defer teardown(t, rdb, []string{keyOne, keyTwo})

	ctx := context.Background()

	actualOne, err := client.TryLockWithRetryStrategy(ctx, keyOne, -1, NewNoRetry())
	if err != nil {
		t.Fatalf("actualOne TryLockWithRetryStrategy error:[%v]", err)
	}

	actualTwo, err := client.TryLockWithRetryStrategy(ctx, keyTwo, 30*time.Second, NewAverageRetry(3, 10*time.Second))
	if err != nil {
		t.Fatalf("actualTwo TryLockWithRetryStrategy error:[%v]", err)
	}

	// test cases
	cases := []struct {
		Name     string
		Actual   *Mutex
		Expected *Mutex
	}{
		{
			"TryLockWithNoRetryStrategy",
			actualOne,
			&Mutex{
				client:        client,
				key:           keyOne,
				expiration:    30 * time.Second,
				value:         actualOne.value,
				retryStrategy: NewNoRetry(),
				watchDog:      actualOne.watchDog,
			},
		},
		{
			"TryLockWithAverageRetryStrategy",
			actualTwo,
			&Mutex{
				client:        client,
				key:           keyTwo,
				expiration:    30 * time.Second,
				value:         actualTwo.value,
				retryStrategy: NewAverageRetry(3, 10*time.Second),
				watchDog:      nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			compareMutex(t, c.Expected, c.Actual)
		})
	}
}

func TestClient_TryLockWithWatchDog(t *testing.T) {
	// init redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})
	// close redis client
	defer rdb.Close()
	// init redislock client
	client, err := NewDefaultClient(rdb)
	if err != nil {
		t.Fatalf("NewDefaultClient error:[%v]", err)
	}
	keyOne := "testOne"
	keyTwo := "testTwo"
	keyThree := "testThree"
	defer teardown(t, rdb, []string{keyOne, keyTwo, keyThree})

	actualOne, err := client.TryLockWithWatchDog(context.Background(), keyOne, nil)
	if err != nil {
		t.Fatalf("actulOne TryLockWithWatchDog error:[%v]", err)
	}

	actualTwo, err := client.TryLockWithWatchDog(context.Background(), keyTwo, NewDefaultWatchDog())
	if err != nil {
		t.Fatalf("actualTwo TryLockWithWatchDog error:[%v]", err)
	}

	actualThree, err := client.TryLockWithWatchDog(context.Background(), keyThree, NewWatchDog(10*time.Second))
	if err != nil {
		t.Fatalf("actualThree TryLockWithWatchDog error:[%v]", err)
	}

	// test cases
	cases := []struct {
		Name     string
		Actual   *Mutex
		Expected *Mutex
	}{
		{
			"TryLockWithNilWatchDog",
			actualOne,
			&Mutex{
				client:        client,
				key:           keyOne,
				expiration:    30 * time.Second,
				value:         actualOne.value,
				retryStrategy: NewNoRetry(),
				watchDog:      actualOne.watchDog,
			},
		},
		{
			"TryLockWithDefaultWatchDog",
			actualTwo,
			&Mutex{
				client:        client,
				key:           keyTwo,
				expiration:    30 * time.Second,
				value:         actualTwo.value,
				retryStrategy: NewNoRetry(),
				watchDog:      actualTwo.watchDog,
			},
		},
		{
			"TryLockWithWatchDog",
			actualThree,
			&Mutex{
				client:        client,
				key:           keyThree,
				expiration:    10 * time.Second,
				value:         actualThree.value,
				retryStrategy: NewNoRetry(),
				watchDog:      actualThree.watchDog,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			compareMutex(t, c.Expected, c.Actual)
		})
	}
}

func compareMutex(t *testing.T, expected *Mutex, actual *Mutex) {
	t.Helper()
	if expected.client != actual.client {
		t.Errorf("client is not equal,expected %v, got %v", expected.client, actual.client)
	}

	if expected.key != actual.key {
		t.Errorf("key is not equal,expected %v, got %v", expected.key, actual.key)
	}

	if expected.expiration != actual.expiration {
		t.Errorf("expiration is not equal,expected %v, got %v", expected.expiration, actual.expiration)
	}

	if expected.value != actual.value {
		t.Errorf("value is not equal,expected %v, got %v", expected.value, actual.value)
	}

	if fmt.Sprintf("%T", expected.retryStrategy) != fmt.Sprintf("%T", actual.retryStrategy) {
		t.Errorf("retryStrategy is not equal,expected %v, got %v", expected.retryStrategy, actual.retryStrategy)
	}

	if expected.watchDog != actual.watchDog {
		t.Errorf("watchDog is not equal,expected %v, got %v", expected.watchDog, actual.watchDog)
	}
}

func compareClient(t *testing.T, expect, actual *Client) {
	t.Helper()
	if expect.redisClient != actual.redisClient {
		t.Errorf("redis client is not equal,expected %v, got %v", expect.redisClient, actual.redisClient)
	}

	if expect.cipherKey != actual.cipherKey {
		t.Errorf("cipher key is not equal,expected %v, got %v", expect.cipherKey, actual.cipherKey)
	}

	if expect.Cipher != actual.Cipher {
		now := time.Now().String()
		expectValue := make([]byte, len(now))
		actualValue := make([]byte, len(now))
		expect.Cipher.XORKeyStream(expectValue, []byte(now))
		actual.Cipher.XORKeyStream(actualValue, []byte(now))
		for i := 0; i < len(actualValue); i++ {
			if actualValue[i] != expectValue[i] {
				t.Errorf("cipher is not equal,expected %v, got %v", expect.Cipher, actual.Cipher)
			}
		}
	}
}

func teardown(t *testing.T, rc *redis.Client, lockKeys []string) {
	t.Helper()

	for _, key := range lockKeys {
		if err := rc.Del(context.Background(), key).Err(); err != nil {
			t.Fatal(err)
		}
	}

	if err := rc.Close(); err != nil {
		t.Fatal(err)
	}
}
