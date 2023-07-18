package redislock

import "crypto/rc4"

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

type mutexOption struct {
	retryStrategy RetryStrategy
	watchDog      *WatchDog
}
