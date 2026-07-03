package lock

import "errors"

var (
	ErrLockAcquireFailed = errors.New("lock acquire failed")
)
