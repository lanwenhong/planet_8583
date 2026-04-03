package models

import (
	"fmt"
	"strings"
	
	"github.com/elliotchance/pie/v2"
)

type CaptureData struct {
	MchntId         string
	Tid             string
	TotalTxCnt      int
	TotalTxAmt      int
	TotalRefundCnt  int
	TotalRefundAmt  int
	TotalTimeoutCnt int
	TotalAdjustCnt  int
	Records         []*CaptureRecord
	BatchNumber     string
	BatchSyssn      string
}

type CaptureRecord struct {
	Syssn        string
	MchntId      string
	Clisn        string
	Tid          string
	TxAmt        int64
	CardNo       string
	TxTime       string
	TxDate       string
	ExpireDate   string
	PosEntryMode string
	AuthCode     string
	ICCData      string
	RRN          string
	Type         string
	Txcurrcd     string
}

func NewCaptureRecordByRecord(record *RecordPO) *CaptureRecord {
	return &CaptureRecord{
		Syssn:        record.Syssn,
		MchntId:      record.Chnluserid,
		Tid:          record.Chnltermid,
		TxAmt:        record.Txamt,
		CardNo:       record.Ext.CardNo,
		TxTime:       record.Ext.RespTime,
		TxDate:       record.Ext.RespDate,
		ExpireDate:   record.Ext.ExpireDate,
		PosEntryMode: record.Ext.EntryMode,
		AuthCode:     record.Ext.AuthCode,
		ICCData:      record.Ext.ICCData,
		RRN:          record.Chnlsn,
		Clisn:        record.Clisn,
		Type:         record.GetTradeType(),
		Txcurrcd:     record.Txcurrcd,
	}
}

func NewCaptureData(records []*RecordPO) *CaptureData {
	if len(records) == 0 {
		return nil
	}
	var (
		totalTxCnt     int
		totalTxAmt     int
		totalRefundCnt int
		totalRefundAmt int
		captureRecords []*CaptureRecord
	)
	
	for _, record := range records {
		tmp := NewCaptureRecordByRecord(record)
		if tmp.Type == "trade" {
			totalTxCnt++
			totalTxAmt += int(tmp.TxAmt)
		} else {
			totalRefundCnt++
			totalRefundAmt += int(tmp.TxAmt)
		}
		captureRecords = append(captureRecords, tmp)
	}
	
	return &CaptureData{
		MchntId:         records[0].Chnluserid,
		Tid:             records[0].Chnltermid,
		TotalTxCnt:      totalTxCnt,
		TotalTxAmt:      totalTxAmt,
		TotalRefundCnt:  totalRefundCnt,
		TotalRefundAmt:  totalRefundAmt,
		TotalTimeoutCnt: 0,
		TotalAdjustCnt:  0,
		Records:         captureRecords,
	}
}

func NewCaptureDatasByRecords(records []*RecordPO) []*CaptureData {
	groups := pie.GroupBy(records, func(record *RecordPO) string {
		return record.Chnluserid + ":" + record.Chnltermid
	})
	var captureDatas []*CaptureData
	for _, v := range groups {
		captureDatas = append(captureDatas, NewCaptureData(v))
	}
	return captureDatas
}

func (c *CaptureData) WithBatchIndex(batchNumber, batchSyssn string) *CaptureData {
	c.BatchNumber = batchNumber
	c.BatchSyssn = batchSyssn
	return c
}

func (c *CaptureData) GetBatchTotals() []byte {
	batchTotals := fmt.Sprintf(
		"%03d%012d%03d%012d%03d%03d",
		c.TotalTxCnt, c.TotalTxAmt, c.TotalRefundCnt, c.TotalRefundAmt, c.TotalTimeoutCnt, c.TotalAdjustCnt,
	)
	return []byte(batchTotals + strings.Repeat("0", 90-len(batchTotals)))
}

func (c *CaptureData) String() string {
	return fmt.Sprintf(
		"[%s:%s] batch_number=%s batch_syssn=%s txcnt=%d txamt=%d refund_cnt=%d refund_amt=%d",
		c.MchntId, c.Tid, c.BatchNumber, c.BatchSyssn,
		c.TotalTxCnt, c.TotalTxAmt, c.TotalRefundCnt, c.TotalRefundAmt,
	)
}
