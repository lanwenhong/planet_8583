package sruntime

import (
	"context"
	"errors"
	"os"
	"runtime"
	
	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/planet_8583/pkg/config"
	"github.com/lanwenhong/planet_8583/pkg/utils"
)

// InitLog 初始化日志
func InitLog() {
	loglevel, f := logger.LoggerLevelIndex(config.Conf.LogLevel)
	utils.MustTrue(f, errors.New("log level err"))
	logger.Newglog(config.Conf.LogDir, config.Conf.LogFile, config.Conf.LogFileErr, &logger.Glogconf{
		RotateMethod: logger.ROTATE_FILE_DAILY,
		Loglevel:     loglevel,
		Stdout:       config.Conf.LogStdOut,
		Colorful:     config.Conf.Colorfull,
	})
}

// InitConf 初始化配置项
func InitConf() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	utils.MustTrue(len(os.Args) >= 2, errors.New("input param error"))
	config.ParseConfig(os.Args[1])
}

func Init(ctx context.Context) {
	InitConf()
	InitLog()
	InitRuntime(ctx)
}
