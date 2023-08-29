package redislock

import (
	"fmt"
	"testing"
	"time"
)

func TestNewNoRetry(t *testing.T) {
	noRetry := NewNoRetry()
	testCases := []struct {
		Name     string
		Actual   *NoRetry
		Expected *NoRetry
	}{
		{
			Name:     "NewNoRetry",
			Actual:   noRetry,
			Expected: &NoRetry{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if fmt.Sprintf("%T", testCase.Expected) != fmt.Sprintf("%T", testCase.Actual) {
				t.Errorf("retryStrategy is not equal,expected %v, got %v", testCase.Expected, testCase.Actual)
			}
		})
	}
}

func TestNewAverageRetry(t *testing.T) {
	averageRetry := NewAverageRetry(2, 1*time.Second)
	testCases := []struct {
		Name     string
		Actual   *AverageRetry
		Expected *AverageRetry
	}{
		{
			Name:     "NewAverageRetry",
			Actual:   averageRetry,
			Expected: &AverageRetry{maxRetryCount: 2, retryInterval: 1 * time.Second},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if fmt.Sprintf("%T", testCase.Expected) != fmt.Sprintf("%T", testCase.Actual) {
				t.Errorf("retryStrategy is not equal,expected %v, got %v", testCase.Expected, testCase.Actual)
			}
			if testCase.Expected.maxRetryCount != testCase.Actual.maxRetryCount {
				t.Errorf("maxRetryCount is not equal,expected %v, got %v", testCase.Expected.maxRetryCount, testCase.Actual.maxRetryCount)
			}
			if testCase.Expected.retryInterval != testCase.Actual.retryInterval {
				t.Errorf("retryInterval is not equal,expected %v, got %v", testCase.Expected.retryInterval, testCase.Actual.retryInterval)
			}
		})
	}
}

func TestNoRetry_NextRetryTime(t *testing.T) {
	noRetry := NewNoRetry()
	for i, exp := range []time.Duration{0, 0, 0} {
		if got := noRetry.NextRetryTime(); exp != got {
			t.Errorf("case %d: expected %v, got %v", i, exp, got)
		}
	}
}

func TestAverageRetry_NextRetryTime(t *testing.T) {
	averageRetry := NewAverageRetry(2, 1*time.Second)
	for i, exp := range []time.Duration{1 * time.Second, 1 * time.Second, 0} {
		if got := averageRetry.NextRetryTime(); exp != got {
			t.Errorf("case %d: expected %v, got %v", i, exp, got)
		}
	}
}
