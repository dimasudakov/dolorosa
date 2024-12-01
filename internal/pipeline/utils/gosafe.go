package utils

import (
	"context"
	"fmt"
	"runtime/debug"
)

// Go - safe version of 'go func' which recovers panics
func Go(ctx context.Context, fn func(ctx context.Context)) {
	if fn != nil {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[!!!] Recovered panic: %s\n%s", r, debug.Stack())
				}
			}()

			fn(ctx)
		}()
	}
}
