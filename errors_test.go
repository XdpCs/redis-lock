package redislock

import (
	"fmt"
	"testing"
)

func TestIsWatchDogExpiredNotLessThanZero(t *testing.T) {
	type args struct {
		err error
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"IsWatchDogExpiredNotLessThanZero", args{ErrWatchDogExpiredNotLessThanZero}, true},
		{"IsWatchDogExpiredNotLessThanZeroWithWrap", args{fmt.Errorf("errors.Wrap %w", ErrWatchDogExpiredNotLessThanZero)}, true},
		{"NotIsWatchDogExpiredNotLessThanZero", args{ErrWatchDogIsNil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsWatchDogExpiredNotLessThanZero(tt.args.err); got != tt.want {
				t.Errorf("IsWatchDogExpiredNotLessThanZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsWatchDogNotStarted(t *testing.T) {
	type args struct {
		err error
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"IsWatchDogNotStarted", args{ErrWatchDogNotStarted}, true},
		{"IsWatchDogNotStartedWithWrap", args{fmt.Errorf("errors.Wrap %w", ErrWatchDogNotStarted)}, true},
		{"NotIsWatchDogNotStarted", args{ErrWatchDogIsNil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsWatchDogNotStarted(tt.args.err); got != tt.want {
				t.Errorf("IsWatchDogNotStarted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsWatchDogIsNil(t *testing.T) {
	type args struct {
		err error
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"IsWatchDogIsNil", args{ErrWatchDogIsNil}, true},
		{"IsWatchDogIsNilWithWrap", args{fmt.Errorf("errors.Wrap %w", ErrWatchDogIsNil)}, true},
		{"NotIsWatchDogIsNil", args{ErrWatchDogNotStarted}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsWatchDogIsNil(tt.args.err); got != tt.want {
				t.Errorf("IsWatchDogIsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMutexLockFailed(t *testing.T) {
	type args struct {
		err error
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"IsMutexLockFailed", args{ErrMutexLockFailed}, true},
		{"IsMutexLockFailedWithWrap", args{fmt.Errorf("errors.Wrap %w", ErrMutexLockFailed)}, true},
		{"NotIsMutexLockFailed", args{ErrMutexNotHeld}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMutexLockFailed(tt.args.err); got != tt.want {
				t.Errorf("IsMutexLockFailed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMutexNotHeld(t *testing.T) {
	type args struct {
		err error
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"IsMutexNotHeld", args{ErrMutexNotHeld}, true},
		{"IsMutexNotHeldWithWrap", args{fmt.Errorf("errors.Wrap %w", ErrMutexNotHeld)}, true},
		{"NotIsMutexNotHeld", args{ErrMutexNotInitialized}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMutexNotHeld(tt.args.err); got != tt.want {
				t.Errorf("IsMutexNotHeld() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMutexNotInitialized(t *testing.T) {
	type args struct {
		err error
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"IsMutexNotInitialized", args{ErrMutexNotInitialized}, true},
		{"IsMutexNotInitializedWithWrap", args{fmt.Errorf("errors.Wrap %w", ErrMutexNotInitialized)}, true},
		{"NotIsMutexNotInitialized", args{ErrWatchDogNotStarted}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMutexNotInitialized(tt.args.err); got != tt.want {
				t.Errorf("IsMutexNotInitialized() = %v, want %v", got, tt.want)
			}
		})
	}
}
