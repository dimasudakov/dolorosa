package utils

import (
	"context"
	"errors"
	"time"

	"github.com/samber/lo"
)

var ErrTimeoutExceeded = errors.New("check was failed")

// WithTimeout оборачивает функцию f и добавляет таймаут.
func WithTimeout[T any](ctx context.Context, fn func() (T, error), timeout time.Duration) (T, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), timeout)
	defer cancel()

	done := make(chan struct {
		result T
		err    error
	}, 1)

	Go(ctx, func(_ context.Context) {
		defer close(done)

		result, err := fn()

		done <- struct {
			result T
			err    error
		}{result, err}
	})

	select {
	case res := <-done:
		return res.result, res.err
	case <-ctx.Done():
		return lo.Empty[T](), ErrTimeoutExceeded
	}
}
