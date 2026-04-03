package servers

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
	
	"github.com/lanwenhong/planet_8583/pkg/config"
	"github.com/lanwenhong/planet_8583/pkg/constant"
	"github.com/lanwenhong/planet_8583/pkg/models"
	"github.com/lanwenhong/planet_8583/pkg/repo"
	"github.com/lanwenhong/planet_8583/pkg/service"
	"github.com/lanwenhong/planet_8583/pkg/sruntime"
	"github.com/lanwenhong/planet_8583/pkg/utils"
	"github.com/redis/go-redis/v9"
	
	"github.com/elliotchance/pie/v2"
	"github.com/lanwenhong/lgobase/logger"
)

type FailedCapture struct {
	Chnlsn string
	Syssn  string
	Error  error
}

type CaptureResult struct {
	Date       string
	Index      int64
	Total      int
	Success    int
	Failed     int
	FailedList []*FailedCapture
	mu         sync.Mutex
}

// RunBatchCapture
// 批量capture
func RunBatchCapture(ctx context.Context, dateStr string) {
	date := time.Now().AddDate(0, 0, config.Conf.BatchCaptureDay)
	captureYesterday := config.Conf.CaptureYesterday
	if dateStr != "" {
		d, err := time.Parse("20060102", dateStr)
		utils.MustExitOnError(err)
		date = d
		captureYesterday = false
	}
	
	if captureYesterday {
		yesterday := date.AddDate(0, 0, -1)
		key := fmt.Sprintf("pp:capture:%s", yesterday.Format("20060102"))
		utils.LockerSrv.TryLock(ctx, key, time.Hour, func() {
			BatchCapture(ctx, yesterday)
		})
	}
	
	key := fmt.Sprintf("pp:capture:%s", date.Format("20060102"))
	utils.LockerSrv.TryLock(ctx, key, time.Hour, func() {
		BatchCapture(ctx, date)
	})
}

func BatchMGet(ctx context.Context, rdb redis.Cmdable, s ...string) map[string]string {
	batchNum := 1000
	var res map[string]string
	for i := 0; i < len(s); i += batchNum {
		j := i + batchNum
		if j > len(s) {
			j = len(s)
		}
		values := rdb.MGet(ctx, s[i:j]...).Val()
		for k, v := range values {
			if vv, ok := v.(string); ok {
				res[s[i+k]] = vv
			}
		}
	}
	return res
}

func FilterCapturedRecord(ctx context.Context, records []*models.RecordPO) []*models.RecordPO {
	var flags []string
	for _, record := range records {
		flags = append(flags, "pp:cpauted:"+record.Syssn)
	}
	flagMap := BatchMGet(ctx, sruntime.Gsvr.Rd, flags...)
	return pie.Filter(records, func(record *models.RecordPO) bool {
		_, ok := flagMap["pp:cpauted:"+record.Syssn]
		return !ok
	})
}

func BatchCapture(ctx context.Context, date time.Time) {
	dateStr := date.Format("2006-01-02")
	records := FilterCapturedRecord(ctx, repo.RecordRepo.GetRequireCaptureRecords(ctx, date))
	result := &CaptureResult{
		Total: len(records), Date: dateStr,
		Index: sruntime.Gsvr.Rd.Incr(ctx, "pp:capture:cnt:"+dateStr).Val(),
	}
	batchNumber := fmt.Sprintf("%06d", result.Index*2)
	batchSyssn := fmt.Sprintf("%06d", result.Index*2+1)
	captureDatas := models.NewCaptureDatasByRecords(records)
	for _, data := range captureDatas {
		if err := processCapture(ctx, data.WithBatchIndex(batchNumber, batchSyssn)); err == nil {
			result.Success = result.Success + data.TotalTxCnt + data.TotalRefundCnt
		} else {
			result.Failed = result.Failed + data.TotalTxCnt + data.TotalRefundCnt
			for _, record := range data.Records {
				result.FailedList = append(result.FailedList, &FailedCapture{
					Syssn: record.Syssn, Chnlsn: record.RRN, Error: err,
				})
			}
		}
	}
	utils.SafeWithLog(ctx, func() {
		LogCaptureResult(ctx, result)
		EmailCaptureResult(ctx, result)
	})
}

func processCapture(ctx context.Context, capture *models.CaptureData) error {
	return utils.Safe2(func() {
		service.Settle(ctx, capture)
		
		utils.SafeWithLog(ctx, func() {
			// 更新记录状态为已清算
			var syssns []string
			for _, record := range capture.Records {
				syssns = append(syssns, record.Syssn)
				sruntime.Gsvr.Rd.Set(ctx, "pp:cpauted:"+record.Syssn, "1", time.Hour*48)
			}
			repo.RecordRepo.UpdateLastEvent(ctx, syssns, constant.ExtEventCaptured)
		})
	})
}

func LogCaptureResult(ctx context.Context, result *CaptureResult) {
	logger.Infof(
		ctx, "Batch capture completed: Date=%s, Total=%d, Success=%d, Failed=%d",
		result.Date, result.Total, result.Success, result.Failed,
	)
	if len(result.FailedList) > 0 {
		logger.Errorf(ctx, "Failed captures:")
		for _, failed := range result.FailedList {
			logger.Errorf(ctx, "syssn: %s, error: %v", failed.Syssn, failed.Error)
		}
	}
}

func EmailCaptureResult(ctx context.Context, result *CaptureResult) {
	if result.Total == 0 && result.Index != 1 {
		logger.Infof(ctx, "No captures to send email")
		return
	}
	
	if len(config.Conf.MailTo) == 0 || config.Conf.MailUserName == "" {
		return
	}
	
	subject := "Finish capture Airwallex transactions"
	if config.Conf.Env == "qa" {
		subject = "[Testing] " + subject
	}
	
	body := fmt.Sprintf(
		"Capture Date:%s, Capture Records Total=%d, Success=%d, Failed=%d",
		result.Date, result.Total, result.Success, result.Failed,
	)
	
	host := utils.OR(os.Getenv("MAIL_HOST"), "smtp.exmail.qq.com")
	port := utils.SafeToInt(utils.OR(os.Getenv("MAIL_PORT"), "465"))
	username := utils.OR(config.Conf.MailUserName, os.Getenv("MAIL_USERNAME"))
	password := utils.OR(config.Conf.MailPassword, os.Getenv("MAIL_PASSWORD"))
	logger.Infof(ctx, "host:%s port:%d username:%s", host, port, username)
	
	srv := utils.NewMailService(host, port, username, password)
	srv.Send(config.Conf.MailTo, subject, body, false)
	logger.Infof(ctx, "SendMail[to:%v subject:%s] success.", config.Conf.MailTo, subject)
}
