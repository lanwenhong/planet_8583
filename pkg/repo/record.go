package repo

import (
	"context"
	"fmt"
	"time"
	
	"github.com/elliotchance/pie/v2"
	"github.com/lanwenhong/planet_8583/pkg/config"
	"github.com/lanwenhong/planet_8583/pkg/models"
	"github.com/lanwenhong/planet_8583/pkg/sruntime"
	"github.com/lanwenhong/planet_8583/pkg/utils"
)

var RecordRepo RecordDB

type RecordDB interface {
	GetRequireCaptureRecords(ctx context.Context, date time.Time) []*models.RecordPO
	UpdateLastEvent(ctx context.Context, syssns []string, lastEvent string)
}

type RecordDBImpl struct {
	Busicds []string
}

func (r *RecordDBImpl) GetRequireCaptureRecords(ctx context.Context, date time.Time) []*models.RecordPO {
	if len(r.Busicds) == 0 {
		return []*models.RecordPO{}
	}
	var ret []*models.RecordPO
	suffix := models.ParseTableSuffix(date.Format("20060102"))
	extTable := suffix.GetTableName("record_card_ext")
	recordTable := suffix.GetTableName("record")
	startTime := date.Format("2006-01-02")
	endTime := date.Add(24 * time.Hour).Format("2006-01-02")
	utils.MustNil(sruntime.Gsvr.Dbs.OrmPools["qf_trade"].
		WithContext(ctx).
		Table(fmt.Sprintf("%s as r", recordTable)).
		Select("ext.cardlastevent, r.userid, r.chnluserid,r.retcd, r.status,r.cancel, r.busicd, r.chnlid, r.chnlsn, r.clisn, r.txamt,r.chnltermid, r.txcurrcd, r.syssn, r.ext").
		Joins(fmt.Sprintf("LEFT JOIN %s as ext ON ext.syssn = r.syssn", extTable)).
		Where("busicd in ?", r.Busicds).
		Where("r.status = 1 AND r.retcd = '0000' AND cancel in (0, 3, 5)").
		Where("r.chnlid = ?", config.Conf.CaptureChannelID).
		Where("r.sysdtm >= ? AND r.sysdtm < ?", startTime, endTime).
		Find(&ret).Error)
	return ret
}

func (r *RecordDBImpl) UpdateLastEvent(ctx context.Context, syssns []string, lastEvent string) {
	if len(syssns) == 0 {
		return
	}
	groupMap := pie.GroupBy(syssns, func(syssn string) string { return syssn[:6] })
	for suffix, items := range groupMap {
		utils.MustNil(sruntime.Gsvr.Dbs.OrmPools["qf_trade"].
			WithContext(ctx).
			Table(fmt.Sprintf("record_card_ext_%s", suffix)).
			Where("syssn in ?", items).
			Update("cardlastevent", lastEvent).Error)
	}
}

func InitRecordRepo(_ context.Context) {
	var busicds []string
	busicds = append(busicds, config.Conf.PaymentBusicds...)
	busicds = append(busicds, config.Conf.RefundBusicds...)
	RecordRepo = &RecordDBImpl{Busicds: busicds}
}
