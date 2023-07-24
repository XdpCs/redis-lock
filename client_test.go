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
		t.Fatal("actualOne NewClient error")
	}

	actualTwo, err := NewClient(rdb, WithCipherKey("11181114"))
	if err != nil {
		t.Fatal("actualTwo NewClient error")
	}

	// init cipher
	cipher, err := rc4.NewCipher([]byte("1114"))
	if err != nil {
		t.Fatal("init cipher error")
	}

	actualThree, err := NewClient(rdb, WithCipher(cipher))
	if err != nil {
		t.Fatal("actualThree NewClient error")
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
			&Client{redisClient: rdb, cipherKey: "11181114", Cipher: actualTwo.Cipher},
		},
		{
			Name:     "NewClientWithCipher",
			Actual:   actualThree,
			Expected: &Client{redisClient: rdb, cipherKey: "-1", Cipher: cipher},
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
		t.Fatal("NewDefaultClient error")
	}

	// test cases
	compareClient(t, &Client{redisClient: rdb, cipherKey: "1118", Cipher: client.Cipher}, client)
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
		t.Fatal("NewDefaultClient error")
	}
	keyOne := "testOne"
	keyTwo := "testTwo"
	defer teardown(t, rdb, []string{keyOne, keyTwo})

	ctxOne := context.Background()
	ctxTwo := context.Background()

	actualOne, err := client.TryLock(ctxOne, keyOne, -1)
	if err != nil {
		t.Fatal("actualOne TryLock error")
	}

	actualTwo, err := client.TryLock(ctxTwo, keyTwo, 10*time.Second)
	if err != nil {
		t.Fatal("actualTwo TryLock error")
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
		t.Errorf("cipher is not equal,expected %v, got %v", expect.Cipher, actual.Cipher)
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
