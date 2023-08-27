package redislock

import "errors"

var (
	ErrWatchDogExpiredNotLessThanZero = errors.New("watch dog expired not less than zero")
	ErrWatchDogNotStarted             = errors.New("watch dog not started")
	ErrWatchDogIsNil                  = errors.New("watch dog is nil")
	ErrMutexLockFailed                = errors.New("mutex locks failed")
	ErrMutexNotHeld                   = errors.New("mutex not held")
	ErrMutexNotInitialized            = errors.New("mutex not initialized")
)

// IsWatchDogExpiredNotLessThanZero returns true if err is ErrWatchDogExpiredNotLessThanZero.
func IsWatchDogExpiredNotLessThanZero(err error) bool {
	return errors.Is(err, ErrWatchDogExpiredNotLessThanZero)
}

// IsWatchDogNotStarted returns true if err is ErrWatchDogNotStarted.
func IsWatchDogNotStarted(err error) bool {
	return errors.Is(err, ErrWatchDogNotStarted)
}

// IsWatchDogIsNil returns true if err is ErrWatchDogIsNil.
func IsWatchDogIsNil(err error) bool {
	return errors.Is(err, ErrWatchDogIsNil)
}

// IsMutexLockFailed returns true if err is ErrMutexLockFailed.
func IsMutexLockFailed(err error) bool {
	return errors.Is(err, ErrMutexLockFailed)
}

// IsMutexNotHeld returns true if err is ErrMutexNotHeld.
func IsMutexNotHeld(err error) bool {
	return errors.Is(err, ErrMutexNotHeld)
}

// IsMutexNotInitialized returns true if err is ErrMutexNotInitialized.
func IsMutexNotInitialized(err error) bool {
	return errors.Is(err, ErrMutexNotInitialized)
}
