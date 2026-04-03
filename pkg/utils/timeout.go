package utils

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"
)

func DoWithTimeout(ctx context.Context, fn func() error, timeout time.Duration) error {
	nextCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// create channel with buffer size 1 to avoid goroutine leak
	done := make(chan error, 1)
	panicChan := make(chan any, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- fmt.Sprintf("%+v\n\n%s", p, strings.TrimSpace(string(debug.Stack())))
			}
		}()
		done <- fn()
	}()

	select {
	case p := <-panicChan:
		panic(p)
	case err := <-done:
		return err
	case <-nextCtx.Done():
		return nextCtx.Err()
	}
}
