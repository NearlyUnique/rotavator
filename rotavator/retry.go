package rotavator

import (
	"context"
	"errors"
	"time"
)

type MaxAttemptsErr struct {
	Err error
}

// Error satisfy error interface
func (e MaxAttemptsErr) Error() string {
	return e.Err.Error()
}

// Unwrap satisfy error interface
func (e MaxAttemptsErr) Unwrap() error {
	return e.Err
}

// ExitRetryErr enclose your error in this to terminate the Retry loop immediately
type ExitRetryErr struct {
	Err error
}

// Error satisfy error interface
func (e ExitRetryErr) Error() string {
	return e.Err.Error()
}

// Unwrap satisfy error interface
func (e ExitRetryErr) Unwrap() error {
	return e.Err
}

// Retry work. use context.WithDeadline() or context.WithTimeout() to s et time limit
// use ExitRetryErr{Err:err} to terminate without further retries
// MaxAttemptsErr{} is returned if attempts == max
func Retry(ctx context.Context, max int, fn func() error) error {
	d := 100 * time.Millisecond
	var err error
	for attempt := 0; attempt < max; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}
		var exitErr ExitRetryErr
		if errors.As(err, &exitErr) {
			return exitErr.Err
		}
		if attempt+1 == max {
			// don't wait if we're not going round again
			break
		}
		t := time.NewTicker(d)
		select {
		case <-ctx.Done():
			// ran out of time?
			t.Stop()
			return ctx.Err()
		case <-t.C:
			// backing off
			break
		}
		d <<= 2
	}
	return MaxAttemptsErr{err}
}
