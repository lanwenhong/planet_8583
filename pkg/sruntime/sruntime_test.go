package sruntime

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSruntime(t *testing.T) {
	ctx := context.Background()
	InitRuntime(ctx)

	t.Run("graceful shutdown", func(t *testing.T) {
		var cnt int
		Gsvr.Go(func() {
			Gsvr.AddShutdown("go1", func() {
				time.Sleep(time.Millisecond * 100)
				cnt += 1
			})
		})
		Gsvr.Go(func() {
			Gsvr.AddShutdown("go2", func() {
				time.Sleep(time.Millisecond * 100)
				cnt += 1
			})
		})
		go func() {
			time.Sleep(time.Millisecond * 50)
			Gsvr.SignalCh <- syscall.SIGTERM
		}()
		Gsvr.GracefulShutdown(ctx, time.Millisecond*150)
		// 150ms，只会有一个goroutine执行完关闭
		assert.Equal(t, 1, cnt)
	})
}
