package service

import (
	"context"
	"testing"
	
	"github.com/lanwenhong/planet_8583/pkg/config"
	"github.com/lanwenhong/planet_8583/pkg/models"
)

func TestSettle(t *testing.T) {
	config.Conf.PlanetAddr = "terminal.uat.planetpayment.com"
	config.Conf.PlanetPort = 40860
	config.Conf.PlanetCertPath = "/home/qfpay/qfconf/keys/planetd_root.cert"
	
	ctx := context.Background()
	Settle(ctx, &models.CaptureData{
		MchntId:         "188000344333",
		Tid:             "99998888",
		TotalTxCnt:      0,
		TotalTxAmt:      0,
		TotalRefundCnt:  1,
		TotalRefundAmt:  100000,
		TotalTimeoutCnt: 0,
		TotalAdjustCnt:  0,
		BatchNumber:     "000014",
		BatchSyssn:      "000015",
		Records: []*models.CaptureRecord{
			{
				Syssn:        "20260408180500020001967732",
				MchntId:      "188000344333",
				Tid:          "99998888",
				Clisn:        "000021",
				TxAmt:        100000,
				CardNo:       "4514617557672096",
				TxTime:       "141831",
				TxDate:       "0408",
				ExpireDate:   "2803",
				PosEntryMode: "072",
				AuthCode:     "781476",
				ICCData:      "9F260828B9FE3E82B782C49F2701809F100706011203A000009F3704D73881A69F360202E7950500000000009A032604039C01009F02060000001000005F2A020344820220209F1A0203449F3303E0B8C89F3501228407A00000000310109F0902008C9F6E04207000",
				RRN:          "609869733899",
				Type:         "refund",
				Txcurrcd:     "344",
			},
		},
	})
	
}
