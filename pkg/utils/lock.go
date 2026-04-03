package utils

import (
	"context"
	"time"
	
	"github.com/google/uuid"
	"github.com/lanwenhong/lgobase/logger"
	"github.com/redis/go-redis/v9"
)

var LockerSrv *Locker

type Locker struct {
	Redis redis.Cmdable
}

func NewLocker(redis redis.Cmdable) *Locker {
	return &Locker{Redis: redis}
}

func (k *Locker) TryLock(ctx context.Context, key string, expiration time.Duration, fn func()) {
	// 获取锁，如果没有获取到直接报错
	uniqueValue := uuid.NewString()
	logger.Infof(ctx, "TryLock key: %s, uniqueValue: %s", key, uniqueValue)
	if ok, err := k.Redis.SetNX(ctx, key, uniqueValue, expiration).Result(); err != nil || !ok {
		panic("get lock error")
	}
	cancelRenew := make(chan struct{})
	
	// 自动释放锁
	defer func() {
		close(cancelRenew)
		luaScript := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end`
		k.Redis.Eval(ctx, luaScript, []string{key}, uniqueValue)
	}()
	
	// 自动续期
	go func() {
		ticker := time.NewTicker(expiration / 2) // 每半个过期时间续期
		defer ticker.Stop()
		
		for {
			select {
			case <-cancelRenew:
				return
			case <-ticker.C:
				val, err := k.Redis.Get(ctx, key).Result()
				if err != nil || val != uniqueValue {
					return
				}
				k.Redis.Expire(ctx, key, expiration)
			}
		}
	}()
	
	// 执行方法
	fn()
}

func InitLocker(r redis.Cmdable) {
	LockerSrv = NewLocker(r)
}
