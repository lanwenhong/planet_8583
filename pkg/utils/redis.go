package utils

import (
	"context"
	"time"
	
	"github.com/lanwenhong/lgobase/logger"
	"github.com/redis/go-redis/v9"
)

// RedisLoggerHook redis 日志hook
type RedisLoggerHook struct {
}

func NewRedisLoggerHook() *RedisLoggerHook {
	return &RedisLoggerHook{}
}

func (r *RedisLoggerHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (r *RedisLoggerHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		if ctx.Value("redis_log") != "true" {
			return next(ctx, cmd)
		}
		
		start := time.Now()
		err := next(ctx, cmd)
		duration := time.Since(start)
		logger.Infof(ctx, "Redis cmd=%s args=%v result=%v err=%v cost=%s",
			cmd.Name(), cmd.Args(), cmd.String(), err, duration)
		return err
	}
}

func (r *RedisLoggerHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		if ctx.Value("redis_log") != "true" {
			return next(ctx, cmds)
		}
		
		start := time.Now()
		err := next(ctx, cmds)
		duration := time.Since(start)
		for i, cmd := range cmds {
			logger.Infof(ctx, "Redis Pipeline[%d] cmd=%s args=%v result=%v err=%v",
				i, cmd.Name(), cmd.Args(), cmd.String(), cmd.Err())
		}
		logger.Infof(ctx, "Redis Pipeline: count=%d, cost=%s, err=%v",
			len(cmds), duration, err)
		
		return err
	}
}

func WithRedisLog(ctx context.Context) context.Context {
	return context.WithValue(ctx, "redis_log", "true")
}
