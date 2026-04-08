package models

import (
	"github.com/elliotchance/pie/v2"
	"github.com/lanwenhong/planet_8583/pkg/config"
)

type RecordPO struct {
	ID         int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Busicd     string     `json:"busicd" gorm:"column:busicd;type:varchar(6);"`
	Userid     int        `json:"userid" gorm:"column:userid;"`
	Syssn      string     `json:"syssn" gorm:"column:syssn;type:varchar(40);"`
	Clisn      string     `json:"clisn" gorm:"column:clisn;type:varchar(6);"`
	Ext        *RecordExt `json:"ext" gorm:"column:ext;type:json;serializer:json"`
	Txamt      int64      `json:"txamt" gorm:"column:txamt;not null"`
	Chnlsn     string     `json:"chnlsn" gorm:"column:chnlsn;"`
	Chnltermid string     `json:"chnltermid" gorm:"column:chnltermid;type:varchar(64)"`
	Chnluserid string     `json:"chnluserid" gorm:"column:chnluserid;type:varchar(64)"`
	Txcurrcd   string     `json:"txcurrcd" gorm:"column:txcurrcd;type:varchar(3);"`
	Retcd      string     `json:"retcd" gorm:"column:retcd;type:varchar(4);"`
	Status     int16      `json:"status" gorm:"column:status;"`
	Cancel     int16      `json:"cancel" gorm:"column:cancel;"`
	
	CardLastEvent string `json:"cardlastevent" gorm:"column:cardlastevent;"`
}

type RecordExt struct {
	EntryMode  string `json:"entry_mode,omitempty"`
	CardNo     string `json:"cardNo,omitempty"`
	ICCData    string `json:"iccdata,omitempty"`
	RespTime   string `json:"time_lt,omitempty"`
	RespDate   string `json:"date_lt,omitempty"`
	ExpireDate string `json:"expired_date,omitempty"`
	AuthCode   string `json:"authCode,omitempty"`
}

func (r *RecordPO) GetTradeType() string {
	if pie.Contains(config.Conf.PaymentBusicds, r.Busicd) {
		return "trade"
	}
	return "refund"
}
