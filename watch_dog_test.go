package redislock

import (
	"testing"
	"time"
)

func TestNewDefaultWatchDog(t *testing.T) {
	watchDog := NewDefaultWatchDog()

	// test cases
	cases := []struct {
		Name             string
		Actual, Expected *WatchDog
	}{
		{
			"NewDefaultWatchDog",
			watchDog,
			&WatchDog{expiration: DefaultExpiration},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			compareWatchDog(t, c.Expected, c.Actual)
		})
	}
}

func TestNewWatchDog(t *testing.T) {
	watchDogOne := NewWatchDog(DefaultExpiration)
	watchDogTwo := NewWatchDog(-1 * time.Second)

	// test cases
	cases := []struct {
		Name             string
		Actual, Expected *WatchDog
	}{
		{
			"NewWatchDogWithPositiveTime",
			watchDogOne,
			&WatchDog{expiration: DefaultExpiration},
		},
		{
			"NewWatchDogWithNegativeTime",
			watchDogTwo,
			&WatchDog{expiration: -1 * time.Second},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			compareWatchDog(t, c.Expected, c.Actual)
		})
	}
}

func compareWatchDog(t *testing.T, actual, expected *WatchDog) {
	if actual.expiration != expected.expiration {
		t.Errorf("actual expiration:[%v], expected expiration:[%v]", actual.expiration, expected.expiration)
	}
	if actual.isStart != expected.isStart {
		t.Errorf("actual isStart:[%v], expected isStart:[%v]", actual.isStart, expected.isStart)
	}
}
