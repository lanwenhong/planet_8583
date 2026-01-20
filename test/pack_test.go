package planet_8583

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"
	
	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/util"
	"github.com/lanwenhong/planet_8583/planet_8583"
)

func TestPackTrack2(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ps := &planet_8583.ProtoStruct{
		MsgType:    "0200",
		TrackData2: "50100439999991007D0810120000000323701",
	}
	
	//ph := &planet_8583.ProtoHandler{}
	ph := planet_8583.NewProtoHandler()
	bdata, err := ph.PackStru(ctx, ps)
	if err != nil {
		t.Fatal(err)
	}
	//bcd := hex.EncodeToString(bdata)
	fs := planet_8583.FormatByte(ctx, bdata)
	logger.Debugf(ctx, "bcd: %s", fs)
}

func TestPackTxamt(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ps := &planet_8583.ProtoStruct{
		MsgType: "0200",
		Txamt:   "66",
	}
	//ph := &planet_8583.ProtoHandler{}
	ph := planet_8583.NewProtoHandler()
	bdata, err := ph.PackStru(ctx, ps)
	if err != nil {
		t.Fatal(err)
	}
	//bcd := hex.EncodeToString(bdata)
	fs := planet_8583.FormatByte(ctx, bdata)
	logger.Debugf(ctx, "bcd: %s", fs)
}

func TestRegisterD63Tag(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	//ph := &planet_8583.ProtoHandler{}
	ph := planet_8583.NewProtoHandler()
	pData := &planet_8583.ProtoStruct{
		MsgType: "0200",
	}
	//pData.Domain63Tags = make(map[string][]byte)
	
	tag12 := &planet_8583.Tag12{
		Len:       "0003",
		Tag:       "12",
		IndiCator: "X",
	}
	
	tagIA := &planet_8583.TagIA{
		Len:          "0004",
		Tag:          "IA",
		HostKeyIndex: "220",
	}
	
	tagIB := &planet_8583.TagIB{
		Len:            "0006",
		Tag:            "IB",
		MacCheckDigits: "F9EA",
	}
	
	tagIC := &planet_8583.TagIC{
		Len:                  "0003",
		Tag:                  "IC",
		InteracTerminalClass: "03",
	}
	
	tagID := &planet_8583.TagID{
		Len:                    "0003",
		Tag:                    "ID",
		InteracCustomerPresent: "1",
	}
	
	tagIE := &planet_8583.TagIE{
		Len:                "0003",
		Tag:                "IE",
		InteracCardPresent: "0",
	}
	
	tagIF := &planet_8583.TagIF{
		Len:                          "0003",
		Tag:                          "IF",
		InteracCardCaptureCapability: "0",
	}
	
	tagIG := &planet_8583.TagIG{
		Len:               "0003",
		Tag:               "IG",
		BalanceinResponse: "0",
	}
	
	tagIH := &planet_8583.TagIH{
		Len:             "0003",
		Tag:             "IH",
		InteracSecurity: "0",
	}
	
	tagIL := &planet_8583.TagIL{
		Len:             "0010",
		Tag:             "IL",
		InteracSecurity: "0000702940000850",
	}
	
	ph.RegisterD63Tag(ctx, "12", pData, tag12)
	ph.RegisterD63Tag(ctx, "IA", pData, tagIA)
	ph.RegisterD63Tag(ctx, "IB", pData, tagIB)
	ph.RegisterD63Tag(ctx, "IC", pData, tagIC)
	ph.RegisterD63Tag(ctx, "ID", pData, tagID)
	ph.RegisterD63Tag(ctx, "IE", pData, tagIE)
	ph.RegisterD63Tag(ctx, "IF", pData, tagIF)
	ph.RegisterD63Tag(ctx, "IG", pData, tagIG)
	ph.RegisterD63Tag(ctx, "IH", pData, tagIH)
	ph.RegisterD63Tag(ctx, "IL", pData, tagIL)
	
	for _, k := range pData.Domain63TagKey {
		logger.Debugf(ctx, "tag: %s", k)
	}
	
	bdata, err := ph.PackStru(ctx, pData)
	if err != nil {
		t.Fatal(err)
	}
	//bcd := hex.EncodeToString(bdata)
	
	fs := planet_8583.FormatByte(ctx, bdata)
	logger.Debugf(ctx, "bcd: %s", fs)
	
}

func TestAuthReq(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := planet_8583.NewProtoHandler()
	iccdata := "9F2608BB4FD0027DF6D4EC9F2701809F100706010A03A000009F37045AED67039F36020A31950580C00000009A032510279C01009F02060000000001005F2A02034482021C009F1A0203449F3303E028C89F34031E03009F3501228407A00000000310109F0902008C9F1E0843415357383332309F0306000000000000"
	biccdata, _ := hex.DecodeString(iccdata)
	pData := &planet_8583.ProtoStruct{
		MsgType:      "0100",
		CardNo:       "4336680006896670",
		ProcessingCd: "000000",
		Txamt:        "666",
		Syssn:        "100002",
		//PosEntryMode:         "072",
		PosEntryMode:         "021",
		Cardsequencenumber:   "001",
		NetId:                "226",
		PosCondCd:            "00",
		ICCSystemRelatedData: biccdata,
		TrackData2:           "4336680006896670D22022011193265100000",
		Tid:                  "11111111",
		MchntId:              "188000344333",
		CurrencyCd:           "344",
	}
	pData.Domain63Tags = make(map[string][]byte)
	
	// 12 O
	_ = ph.RegisterD63Tag(ctx, "12", pData, &planet_8583.Tag12{
		Len: "0003", Tag: "12", IndiCator: "X",
	})
	// IA O (Tag IA – Host Key Index)
	_ = ph.RegisterD63Tag(ctx, "IA", pData, &planet_8583.TagIA{
		Len: "0004", Tag: "IA", HostKeyIndex: "220",
	})
	// IB O (Tag IB – MAC Check Digits)
	_ = ph.RegisterD63Tag(ctx, "IB", pData, &planet_8583.TagIB{
		Len: "0006", Tag: "IB", MacCheckDigits: "F9EA",
	})
	// IL O (Tag IL – Pin Pad Serial Number)
	_ = ph.RegisterD63Tag(ctx, "IL", pData, &planet_8583.TagIL{
		Len: "0010", Tag: "IL", InteracSecurity: "0000702940000850",
	})
	// FA M
	_ = ph.RegisterD63Tag(ctx, "FA", pData, &planet_8583.TagFA{
		Len: "0003", Tag: "FA", FinalAuthIndicator: "U",
	})
	// TC M
	_ = ph.RegisterD63Tag(ctx, "TC", pData, &planet_8583.TagTC{
		Len: "0003", Tag: "TC", TerminalEntryCapabilities: "0",
	})
	
	if _, err := ph.PackStru(ctx, pData); err != nil {
		t.Fatal(err)
	}
	if err := ph.PackMac(ctx, "BBEFB74400000000"); err != nil {
		t.Fatal(err)
	}
	ph.Pack(ctx)
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
}

func TestAuthAdviceReq(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := planet_8583.NewProtoHandler()
	iccdata := "9F2608BB4FD0027DF6D4EC9F2701809F100706010A03A000009F37045AED67039F36020A31950580C00000009A032510279C01009F02060000000001005F2A02034482021C009F1A0203449F3303E028C89F34031E03009F3501228407A00000000310109F0902008C9F1E0843415357383332309F0306000000000000"
	biccdata, _ := hex.DecodeString(iccdata)
	pData := &planet_8583.ProtoStruct{
		MsgType:              "0120",
		CardNo:               "4336680006896670",
		ProcessingCd:         "020000",
		Txamt:                "551",
		CardholderBilling:    "551",
		Syssn:                "111889",
		TimeLocalTransaction: time.Now().Format("150405"),
		DateLocalTransaction: time.Now().Format("0102"),
		CardDatetime:         "9999",
		PosEntryMode:         "072",
		//PosEntryMode: "021",
		//Cardsequencenumber:   "001",
		NetId:                "226",
		PosCondCd:            "00",
		ICCSystemRelatedData: biccdata,
		// Must be present for Adjustments.
		RetrievalReferenceNumber: "531057747967",
		// Must be present for Offline Sales and Sales Completion.
		AuthorizationIDResponse: "576865",
		Tid:                     "12345678",
		MchntId:                 "188000344333",
		CurrencyCd:              "344",
	}
	pData.Domain63Tags = make(map[string][]byte)
	// 12 O
	_ = ph.RegisterD63Tag(ctx, "12", pData, &planet_8583.Tag12{
		Len: "0003", Tag: "12", IndiCator: "X",
	})
	// TC M
	_ = ph.RegisterD63Tag(ctx, "TC", pData, &planet_8583.TagTC{
		Len: "0003", Tag: "TC", TerminalEntryCapabilities: "0",
	})
	
	if _, err := ph.PackStru(ctx, pData); err != nil {
		t.Fatal(err)
	}
	//if err := ph.PackMac(ctx, "BBEFB74400000000"); err != nil {
	//	t.Fatal(err)
	//}
	ph.Pack(ctx)
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
}

func TestAuthAdviceRepeatReq(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := planet_8583.NewProtoHandler()
	iccdata := "9F2608BB4FD0027DF6D4EC9F2701809F100706010A03A000009F37045AED67039F36020A31950580C00000009A032510279C01009F02060000000001005F2A02034482021C009F1A0203449F3303E028C89F34031E03009F3501228407A00000000310109F0902008C9F1E0843415357383332309F0306000000000000"
	biccdata, _ := hex.DecodeString(iccdata)
	pData := &planet_8583.ProtoStruct{
		MsgType:              "0121",
		CardNo:               "4336680006896670",
		ProcessingCd:         "020000",
		Txamt:                "555",
		Syssn:                "100000",
		TimeLocalTransaction: time.Now().Format("150405"),
		DateLocalTransaction: time.Now().Format("0102"),
		CardDatetime:         "2601",
		//PosEntryMode:         "072",
		PosEntryMode:         "021",
		Cardsequencenumber:   "001",
		NetId:                "226",
		PosCondCd:            "00",
		ICCSystemRelatedData: biccdata,
		// Must be present for Adjustments.
		//RetrievalReferenceNumber: "531057742540",
		// Must be present for Offline Sales and Sales Completion.
		//AuthorizationIDResponse: "576821",
		Tid:        "10000000",
		MchntId:    "188000344333",
		CurrencyCd: "344",
		TrackData2: "4336680006896670D22022011193265100000",
	}
	pData.Domain63Tags = make(map[string][]byte)
	
	// 12 O
	_ = ph.RegisterD63Tag(ctx, "12", pData, &planet_8583.Tag12{
		Len: "0003", Tag: "12", IndiCator: "X",
	})
	// TC M
	_ = ph.RegisterD63Tag(ctx, "TC", pData, &planet_8583.TagTC{
		Len: "0003", Tag: "TC", TerminalEntryCapabilities: "0",
	})
	
	if _, err := ph.PackStru(ctx, pData); err != nil {
		t.Fatal(err)
	}
	if err := ph.PackMac(ctx, "BBEFB74400000000"); err != nil {
		t.Fatal(err)
	}
	ph.Pack(ctx)
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
}

func TestFinancialAdviceReq(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := planet_8583.NewProtoHandler()
	iccdata := "9F2608BB4FD0027DF6D4EC9F2701809F100706010A03A000009F37045AED67039F36020A31950580C00000009A032510279C01009F02060000000001005F2A02034482021C009F1A0203449F3303E028C89F34031E03009F3501228407A00000000310109F0902008C9F1E0843415357383332309F0306000000000000"
	biccdata, _ := hex.DecodeString(iccdata)
	pData := &planet_8583.ProtoStruct{
		MsgType:              "0220",
		CardNo:               "4336680006896670",
		ProcessingCd:         "000000",
		Txamt:                "111",
		Syssn:                "100001",
		TimeLocalTransaction: time.Now().Format("150405"),
		DateLocalTransaction: time.Now().Format("0102"),
		CardDatetime:         "2601",
		//PosEntryMode:         "072",
		PosEntryMode:         "021",
		Cardsequencenumber:   "001",
		NetId:                "226",
		PosCondCd:            "00",
		ICCSystemRelatedData: biccdata,
		// Must be present for Adjustments.
		RetrievalReferenceNumber: "531057742540",
		// Must be present for Offline Sales and Sales Completion.
		AuthorizationIDResponse: "576821",
		Tid:                     "10000000",
		MchntId:                 "188000344333",
		CurrencyCd:              "344",
		TrackData2:              "4336680006896670D22022011193265100000",
	}
	pData.Domain63Tags = make(map[string][]byte)
	
	// 12 O
	_ = ph.RegisterD63Tag(ctx, "12", pData, &planet_8583.Tag12{
		Len: "0003", Tag: "12", IndiCator: "X",
	})
	if _, err := ph.PackStru(ctx, pData); err != nil {
		t.Fatal(err)
	}
	if err := ph.PackMac(ctx, "BBEFB74400000000"); err != nil {
		t.Fatal(err)
	}
	ph.Pack(ctx)
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
}

func TestFinancialTrasReq(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := planet_8583.NewProtoHandler()
	iccdata := "9F2608BB4FD0027DF6D4EC9F2701809F100706010A03A000009F37045AED67039F36020A31950580C00000009A032510279C01009F02060000000001005F2A02034482021C009F1A0203449F3303E028C89F34031E03009F3501228407A00000000310109F0902008C9F1E0843415357383332309F0306000000000000"
	biccdata, _ := hex.DecodeString(iccdata)
	pData := &planet_8583.ProtoStruct{
		MsgType:      "0200",
		CardNo:       "4336680006896670",
		ProcessingCd: "000000",
		
		// 第一笔
		//Txamt: "200",
		//Syssn: "200001",
		
		// 第二笔
		//Txamt: "1200",
		//Syssn: "200002",
		
		// 第三笔
		//Txamt: "2400",
		//Syssn: "200003",
		
		// 第四笔
		//Txamt: "1900",
		//Syssn: "200004",
		//Tid:   "20000000",
		
		// 第六笔
		//Txamt: "1910",
		//Syssn: "200005",
		//Tid:   "20000001",
		
		//Txamt: "90",
		//Syssn: "200006",
		//Tid:   "20000001",
		
		// 第7笔
		//Txamt: "9910",
		//Syssn: "200005",
		//Tid:   "20000002",
		
		// 2025-11-16
		// 第8笔
		//Txamt: "110",
		//Syssn: "200007",
		//Tid:   "20000000",
		// 第9笔
		//Txamt: "110",
		//Syssn: "200008",
		//Tid:   "20000000",
		// 第10笔
		//Txamt: "100",
		//Syssn: "200009",
		//Tid:   "20000001",
		// 第11笔
		//Txamt: "100",
		//Syssn: "200009",
		//Tid:   "20000002",
		// 12
		//Txamt: "100",
		//Syssn: "200012",
		//Tid:   "20000003",
		// 13
		//Txamt: "100",
		//Syssn: "200013",
		//Tid:   "20000004",
		// 14
		//Txamt: "100",
		//Syssn: "200014",
		//Tid:   "20000005",
		
		// 2025-11-19
		//Txamt: "100",
		//Syssn: "200015",
		//Tid:   "20000000",
		
		//Txamt: "200",
		//Syssn: "200016",
		//Tid:   "20000000",
		
		Txamt: "300",
		Syssn: "200017",
		Tid:   "20000000",
		
		PosEntryMode:         "021",
		NetId:                "226",
		PosCondCd:            "00",
		TrackData2:           "4336680006896670D22022011193265100000",
		ICCSystemRelatedData: biccdata,
		MchntId:              "188000344333",
		CurrencyCd:           "344",
	}
	pData.Domain63Tags = make(map[string][]byte)
	
	// 12 O
	_ = ph.RegisterD63Tag(ctx, "12", pData, &planet_8583.Tag12{
		Len: "0003", Tag: "12", IndiCator: "X",
	})
	// IA O (Tag IA – Host Key Index)
	_ = ph.RegisterD63Tag(ctx, "IA", pData, &planet_8583.TagIA{
		Len: "0004", Tag: "IA", HostKeyIndex: "220",
	})
	// IB O (Tag IB – MAC Check Digits)
	_ = ph.RegisterD63Tag(ctx, "IB", pData, &planet_8583.TagIB{
		Len: "0006", Tag: "IB", MacCheckDigits: "F9EA",
	})
	// IC O (Tag IB – MAC Check Digits)
	_ = ph.RegisterD63Tag(ctx, "IC", pData, &planet_8583.TagIC{
		Len: "0003", Tag: "IC", InteracTerminalClass: "03",
	})
	// ID
	_ = ph.RegisterD63Tag(ctx, "ID", pData, &planet_8583.TagID{
		Len: "0003", Tag: "ID", InteracCustomerPresent: "1",
	})
	// IE
	_ = ph.RegisterD63Tag(ctx, "IE", pData, &planet_8583.TagIE{
		Len: "0003", Tag: "IE", InteracCardPresent: "0",
	})
	// IF
	_ = ph.RegisterD63Tag(ctx, "IF", pData, &planet_8583.TagIF{
		Len: "0003", Tag: "IF", InteracCardCaptureCapability: "0",
	})
	// IG
	_ = ph.RegisterD63Tag(ctx, "IG", pData, &planet_8583.TagIG{
		Len: "0003", Tag: "IG", BalanceinResponse: "0",
	})
	// IH
	_ = ph.RegisterD63Tag(ctx, "IH", pData, &planet_8583.TagIH{
		Len: "0003", Tag: "IH", InteracSecurity: "0",
	})
	// IL
	_ = ph.RegisterD63Tag(ctx, "IL", pData, &planet_8583.TagIL{
		Len: "0010", Tag: "IL", InteracSecurity: "0000702940000850",
	})
	// FA M
	//_ = ph.RegisterD63Tag(ctx, "FA", pData, &planet_8583.TagFA{
	//	Len: "0003", Tag: "FA", FinalAuthIndicator: "F",
	//})
	//// TC M
	//_ = ph.RegisterD63Tag(ctx, "TC", pData, &planet_8583.TagTC{
	//	Len: "0003", Tag: "TC", TerminalEntryCapabilities: "5",
	//})
	if _, err := ph.PackStru(ctx, pData); err != nil {
		t.Fatal(err)
	}
	//if err := ph.PackMac(ctx, "BBEFB74400000000"); err != nil {
	//	t.Fatal(err)
	//}
	ph.Pack(ctx)
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
}

func TestBatchUploadAdvice(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := planet_8583.NewProtoHandler()
	pData := &planet_8583.ProtoStruct{
		MsgType:      "0320",
		CardNo:       "4336680006896670",
		ProcessingCd: "000000",
		
		//Txamt:                    "100",
		//CardholderBilling:        "100",
		//Syssn:                    "200015",
		//RetrievalReferenceNumber: "532358676121",
		//Tid:                      "20000000",
		
		//Txamt:                    "200",
		//CardholderBilling:        "200",
		//Syssn:                    "200016",
		//RetrievalReferenceNumber: "532358676201",
		//Tid:                      "20000000",
		
		//Txamt:                    "100",
		//CardholderBilling:        "100",
		//Syssn:                    "200015",
		//RetrievalReferenceNumber: "532358676122",
		//Tid:                      "20000000",
		
		Txamt:                    "300",
		CardholderBilling:        "300",
		Syssn:                    "200017",
		RetrievalReferenceNumber: "532358676271",
		Tid:                      "20000000",
		
		TimeLocalTransaction: time.Now().Format("150405"),
		DateLocalTransaction: time.Now().Format("0102"),
		CardDatetime:         "2512",
		PosEntryMode:         "021",
		NetId:                "226",
		PosCondCd:            "00",
		ResponseCode:         "00",
		MchntId:              "188000344333",
		CurrencyCd:           "344",
		CurrencyCdCardholder: "344",
	}
	pData.Domain63Tags = make(map[string][]byte)
	// FA M
	_ = ph.RegisterD63Tag(ctx, "FA", pData, &planet_8583.TagFA{
		Len: "0003", Tag: "FA", FinalAuthIndicator: "U",
	})
	// TC M
	_ = ph.RegisterD63Tag(ctx, "TC", pData, &planet_8583.TagTC{
		Len: "0003", Tag: "TC", TerminalEntryCapabilities: "0",
	})
	if _, err := ph.PackStru(ctx, pData); err != nil {
		t.Fatal(err)
	}
	if err := ph.PackMac(ctx, "BBEFB74400000000"); err != nil {
		t.Fatal(err)
	}
	ph.Pack(ctx)
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
}

func NewD63BatchTotals(txCnt, txAmt, refundCnt, refundAmt, timeoutCnt, adjustCnt int) []byte {
	batchTotals := fmt.Sprintf("%03d%012d%03d%012d%03d%03d", txCnt, txAmt, refundCnt, refundAmt, timeoutCnt, adjustCnt)
	return []byte(batchTotals + strings.Repeat("0", 90-len(batchTotals)))
}

func TestBatchSettlement(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := planet_8583.NewProtoHandler()
	
	pData := &planet_8583.ProtoStruct{
		MsgType:      "0500",
		ProcessingCd: "920000",
		Syssn:        "400012",
		NetId:        "226",
		Tid:          "20000005",
		MchntId:      "188000344333",
		BatchNumber:  []byte("400012"),
		BatchTotals:  NewD63BatchTotals(1, 100, 1, 100, 0, 0),
	}
	pData.Domain63Tags = make(map[string][]byte)
	if _, err := ph.PackStru(ctx, pData); err != nil {
		t.Fatal(err)
	}
	if err := ph.PackMac(ctx, "BBEFB74400000000"); err != nil {
		t.Fatal(err)
	}
	ph.Pack(ctx)
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
}

func TestKeyExchangeReq(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := planet_8583.NewProtoHandler()
	pData := &planet_8583.ProtoStruct{
		MsgType:      "0800",
		ProcessingCd: "920000",
		Syssn:        "000900",
		Tid:          "11411111",
		MchntId:      "188000344333",
		NetId:        "226",
	}
	tagIL := &planet_8583.TagIL{
		Len:             "0010",
		Tag:             "IL",
		InteracSecurity: "0000702940000850",
	}
	
	tagPP := &planet_8583.TagPP{
		Len:                   "0018",
		Tag:                   "PP",
		PlanetPaymentPassword: "24504C414E4554245041594D454E5424",
	}
	
	ph.RegisterD63Tag(ctx, "IL", pData, tagIL)
	ph.RegisterD63Tag(ctx, "PP", pData, tagPP)
	
	_, err := ph.PackStru(ctx, pData)
	if err != nil {
		t.Fatal(err)
	}
	ph.Pack(ctx)
	
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
}

func TestKeyExchangeTr31Req(t *testing.T) {
	
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := planet_8583.NewProtoHandler()
	pData := &planet_8583.ProtoStruct{
		MsgType:              "0800",
		ProcessingCd:         "920000",
		Syssn:                "529323",
		TimeLocalTransaction: "133057",
		DateLocalTransaction: "0927",
		NetId:                "052",
	}
	tagIL := &planet_8583.TagIL{
		Len:             "0010",
		Tag:             "IL",
		InteracSecurity: "0000000000000002",
	}
	
	tagKT := &planet_8583.TagKT{
		Len:                      "0003",
		Tag:                      "KT",
		KeyExchangeMechanismType: "2",
	}
	
	tagPP := &planet_8583.TagPP{
		Len:                   "0018",
		Tag:                   "PP",
		PlanetPaymentPassword: "24504C414E4554245041594D454E5424",
	}
	
	err := ph.RegisterD63Tag(ctx, "IL", pData, tagIL)
	if err != nil {
		t.Fatal(err)
	}
	err = ph.RegisterD63Tag(ctx, "KT", pData, tagKT)
	if err != nil {
		t.Fatal(err)
	}
	err = ph.RegisterD63Tag(ctx, "PP", pData, tagPP)
	
	if err != nil {
		t.Fatal(err)
	}
	
	_, err = ph.PackStru(ctx, pData)
	if err != nil {
		t.Fatal(err)
	}
	ph.Pack(ctx)
	
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
	
}
