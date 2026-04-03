package sruntime

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	
	"github.com/lanwenhong/planet_8583/pkg/config"
	"github.com/lanwenhong/planet_8583/pkg/utils"
	
	"github.com/lanwenhong/lgobase/dbenc"
	"github.com/lanwenhong/lgobase/dbpool"
	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/redispool"
	"github.com/redis/go-redis/v9"
	"github.com/sourcegraph/conc"
)

var (
	Gsvr *Sruntime
	once sync.Once
)

type Sruntime struct {
	// redis
	Rop        *redis.Client
	RopCluster *redis.ClusterClient
	Rd         redis.Cmdable
	
	// mysql
	Dbs *dbpool.Dbpool
	
	// shutdown
	ShutDownServers sync.Map
	SignalCh        chan os.Signal
	WorkWG          conc.WaitGroup
}

func (s *Sruntime) Go(fn func()) { s.WorkWG.Go(fn) }

func (s *Sruntime) WithRedis(ctx context.Context) *Sruntime {
	if config.Conf.Cluster == false {
		rdb := redispool.NewGrPool(
			ctx, config.Conf.RedisUser, config.Conf.RedisPasswd, config.Conf.Db, config.Conf.RedisAddr, config.Conf.PoolSize, config.Conf.MinIdle,
			time.Duration(config.Conf.ConnTimeout)*time.Second,
			time.Duration(config.Conf.ReadTimeout)*time.Second,
			time.Duration(config.Conf.WriteTimeout)*time.Second,
		)
		if config.Conf.RedisLog {
			rdb.AddHook(utils.NewRedisLoggerHook())
		}
		s.Rop = rdb
		s.Rd = rdb
	} else {
		rdb := redispool.NewClusterPool(
			ctx, config.Conf.RedisUser, config.Conf.RedisPasswd, config.Conf.RedisAddrs, config.Conf.PoolSize, config.Conf.MinIdle,
			time.Duration(config.Conf.ConnTimeout)*time.Second,
			time.Duration(config.Conf.ReadTimeout)*time.Second,
			time.Duration(config.Conf.WriteTimeout)*time.Second,
		)
		if config.Conf.RedisLog {
			rdb.AddHook(utils.NewRedisLoggerHook())
		}
		s.RopCluster = rdb
		s.Rd = rdb
	}
	return s
}
func (s *Sruntime) WithDbs(ctx context.Context) *Sruntime {
	if config.Conf.TokenFile == "" {
		return s
	}
	dbConf := dbenc.DbConfNew(ctx, config.Conf.TokenFile)
	s.Dbs = dbpool.DbpoolNew(dbConf)
	_ = s.Dbs.Add(ctx, "qf_trade", config.Conf.TradeDB, dbpool.USE_GORM)
	return s
}

func (s *Sruntime) GracefulShutdown(ctx context.Context, timeout time.Duration) {
	signal.Notify(s.SignalCh, os.Interrupt, syscall.SIGTERM)
	
	<-s.SignalCh
	logger.Info(ctx, "Received shutdown signal. Initiating graceful shutdown...")
	_ = utils.DoWithTimeout(ctx, func() error {
		s.ShutdownAll()
		return nil
	}, timeout)
}

func (s *Sruntime) AddShutdown(name string, fn func()) {
	s.ShutDownServers.Store(name, fn)
}

func (s *Sruntime) ShutdownAll() {
	s.ShutDownServers.Range(func(key, value interface{}) bool {
		value.(func())()
		return true
	})
}

func (s *Sruntime) Shutdown(name string) {
	if fn, ok := s.ShutDownServers.Load(name); ok {
		fn.(func())()
	}
}

func NewSruntime() *Sruntime {
	return &Sruntime{ShutDownServers: sync.Map{}, SignalCh: make(chan os.Signal, 1)}
}

func InitRuntime(ctx context.Context) {
	once.Do(func() {
		Gsvr = NewSruntime().WithDbs(ctx).WithRedis(ctx)
	})
}
