package utils

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	
	"github.com/lanwenhong/lgobase/logger"
)

func Must(fn func() error) {
	if err := fn(); err != nil {
		panic(err)
	}
}

func MustTrue(flag bool, err error) {
	if !flag {
		panic(err)
	}
}

func MustNil(err error) {
	if err != nil {
		panic(err)
	}
}

func MustExitOnError(err error) {
	if err != nil {
		logger.Infof(context.Background(), "exit error: %v", err)
		os.Exit(1)
	}
}

func SafeWithResp[T any](fn func() T) (resp T) {
	defer func() {
		if p := recover(); p != nil {
		}
	}()
	return fn()
}

func Safe(fn func() error) (err error) {
	defer func() {
		if p := recover(); p != nil {
			if errP, ok := p.(error); ok {
				err = errP
			} else {
				err = fmt.Errorf("panic: %v\n%s", p, string(debug.Stack()))
			}
		}
	}()
	
	err = fn()
	return
}

func Safe2(fn func()) (err error) {
	defer func() {
		if p := recover(); p != nil {
			if errP, ok := p.(error); ok {
				err = errP
			} else {
				err = fmt.Errorf("panic: %v\n%s", p, string(debug.Stack()))
			}
		}
	}()
	fn()
	return
}

func SafeWithLog(ctx context.Context, fn func()) {
	defer func() {
		if p := recover(); p != nil {
			logger.Infof(ctx, "panic: %v", p)
		}
	}()
	fn()
	return
}

func SafeGo(ctx context.Context, fn func()) {
	go func() {
		defer func() {
			if p := recover(); p != nil {
				logger.Errorf(ctx, "panic: %v\n%s", p, string(debug.Stack()))
			}
		}()
		fn()
	}()
}

func SafeGoWithInfoLog(ctx context.Context, fn func()) {
	go func() {
		defer func() {
			if p := recover(); p != nil {
				logger.Infof(ctx, "panic: %v", p)
			}
		}()
		fn()
	}()
}
