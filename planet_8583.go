package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	
	"github.com/lanwenhong/planet_8583/pkg/config"
	"github.com/lanwenhong/planet_8583/pkg/repo"
	"github.com/lanwenhong/planet_8583/pkg/servers"
	"github.com/lanwenhong/planet_8583/pkg/sruntime"
	"github.com/lanwenhong/planet_8583/pkg/utils"
	"github.com/lanwenhong/planet_8583/planet_8583"
)

func Init(ctx context.Context) {
	sruntime.Init(ctx)
	utils.SetLoggerLevel(config.Conf.LogLevel)
	utils.InitLocker(sruntime.Gsvr.Rd)
	repo.InitRecordRepo(ctx)
}

func main() {
	ctx := context.Background()
	mode := flag.String("mode", "", "clear:清空数据 batch:单笔清算 batch_capture:执行清算")
	date := flag.String("date", "", "batch_capture, 指定日期,格式:20060102, 默认当前日期")
	tid := flag.String("tid", "87654321", "tid")
	txcnt := flag.Int("txcnt", 0, "交易笔数")
	txamt := flag.Int("txamt", 0, "交易金额, 单位是分")
	refundCnt := flag.Int("refundcnt", 0, "退款笔数")
	refundamt := flag.Int("refundamt", 0, "退款金额")
	timeoutcnt := flag.Int("timeoutcnt", 0, "超时笔数")
	adjustcnt := flag.Int("adjustcnt", 0, "调账笔数")
	_ = flag.CommandLine.Parse(os.Args[2:])
	Init(ctx)
	fmt.Printf("当前模式: %s", *mode)
	switch *mode {
	case "batch_capture":
		ctx := context.WithValue(ctx, "trace_id", "batch_capture")
		servers.RunBatchCapture(ctx, utils.PtrString(date))
	case "clear":
		planet_8583.Clear(ctx, *tid)
		break
	case "batch":
		planet_8583.Settle(ctx, *tid, *txcnt, *txamt, *refundCnt, *refundamt, *timeoutcnt, *adjustcnt)
		break
	default:
		fmt.Printf("未知模式: %s", *mode)
		os.Exit(1)
	}
}
