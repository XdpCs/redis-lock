package redislock

import (
	"context"
	"fmt"
	"time"
)

const DefaultExpiration = 30 * time.Second

// WatchDog is a watch dog for redis lock.
type WatchDog struct {
	expiration time.Duration      // lock expiration.
	cancelFunc context.CancelFunc // cancel function.
	isStart    uint32             // watch dog is start or not start.
}

// NewWatchDog creates a new WatchDog.
func NewWatchDog(expiration time.Duration) *WatchDog {
	return newWatchDog(expiration)
}

// NewDefaultWatchDog creates a new WatchDog with default expiration.
func NewDefaultWatchDog() *WatchDog {
	return newWatchDog(DefaultExpiration)
}

func newWatchDog(expiration time.Duration) *WatchDog {
	return &WatchDog{
		expiration: expiration,
	}
}

func checkWatchDog(watchDog *WatchDog) error {
	if watchDog == nil {
		return ErrWatchDogIsNil
	}

	if watchDog.expiration <= 0 {
		return ErrWatchDogExpiredNotLessThanZero
	}
	return nil
}

func checkWatchDogReturnWatchDog(watchDog *WatchDog) (*WatchDog, error) {
	if err := checkWatchDog(watchDog); err != nil {
		if IsWatchDogIsNil(err) {
			watchDog = NewDefaultWatchDog()
		} else {
			return nil, fmt.Errorf("checkWatchDog error: %w", err)
		}
	}
	return watchDog, nil
}
