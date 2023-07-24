package redislock

import "time"

// RetryStrategy is the interface used by redislock to retry.
type RetryStrategy interface {
	NextRetryTime() time.Duration
}

// NoRetry is a retry strategy that never retries.
type NoRetry struct{}

// NextRetryTime returns 0, which means no retry.
func (n *NoRetry) NextRetryTime() time.Duration {
	return 0
}

// NewNoRetry creates a new NoRetry.
func NewNoRetry() *NoRetry {
	return &NoRetry{}
}

// AverageRetry is a retry strategy that retries for a fixed number of times with a fixed interval.
type AverageRetry struct {
	maxRetryCount uint
	retryInterval time.Duration
}

// NextRetryTime returns retryInterval if maxRetryCount is greater than 0, otherwise returns 0.
func (a *AverageRetry) NextRetryTime() time.Duration {
	if a.maxRetryCount == 0 {
		return 0
	}
	a.maxRetryCount--

	return a.retryInterval
}

// NewAverageRetry creates a new AverageRetry.
func NewAverageRetry(maxRetryCount uint, retryInterval time.Duration) *AverageRetry {
	return &AverageRetry{
		maxRetryCount: maxRetryCount,
		retryInterval: retryInterval,
	}
}
