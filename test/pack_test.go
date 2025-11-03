package planet_8583

import (
	"context"
	"encoding/hex"
	"testing"

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

func TestFinancialTrasReq(t *testing.T) {

	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	//ph := &planet_8583.ProtoHandler{}
	ph := planet_8583.NewProtoHandler()
	iccdata := "9F2608E59FCA338AE60ADD9F2701809F100706010103A030029F37044957D0F79F3602000C950500000000009A032510209C01009F02060000000001005F2A020344820219809F1A0203449F03060000000000009F3303E038C89F3501219F1E0808690710522059879F090200018408A0000000250105029F4104000000019F34031F0202"
	biccdata, _ := hex.DecodeString(iccdata)
	pData := &planet_8583.ProtoStruct{
		MsgType: "0200",
		CardNo:  "4336680006896670",
		//ProcessingCd: "002000",
		ProcessingCd: "000000",
		Txamt:        "555",
		Syssn:        "000888",
		//PosEntryMode:         "021",
		PosEntryMode:       "072",
		Cardsequencenumber: "001",
		NetId:              "226",
		PosCondCd:          "00",
		TrackData2:         "4336680006896670D22022011193265100000",
		//ICCSystemRelatedData: "9F260825906395A2142A119F2701809F100706010A03A000009F37048F0956B49F360209AA950500000000009A032510169C01009F02060000000000015F2A020156820220009F1A0201569F3303E0F8C89F3501228407A00000000310109F0902008C9F34031F03029F1E0843415357383332309F0306000000000000",
		//ICCSystemRelatedData: "9F2608CD4C0BC6E5D0F9DA9F2701809F100706010A03A000009F3704896ADAE09F360209FF950580C00808009A032510239C01009F02060000000000015F2A02015682021C009F1A0201569F3303E0F8C89F34031F00009F3501228407A00000000310109F0902008C9F1E0843415357383332309F0306000000000000",
		ICCSystemRelatedData: biccdata,
		Tid:                  "12345678",
		MchntId:              "188000344333",
		CurrencyCd:           "344",
		//Pin:          "AA17EAB7BF18034B",
	}
	pData.Domain63Tags = make(map[string][]byte)

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

	//3bit

	_, err := ph.PackStru(ctx, pData)
	if err != nil {
		t.Fatal(err)
	}
	/*ph.PackMac(ctx, "BBEFB74400000000")
	if err != nil {
		t.Fatal(err)
	}*/
	ph.Pack(ctx)

	//bcd := hex.EncodeToString(ph.Tbuf)
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

func TestRefundReq(t *testing.T) {

	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	//ph := &planet_8583.ProtoHandler{}
	ph := planet_8583.NewProtoHandler()
	iccdata := "9F2608E59FCA338AE60ADD9F2701809F100706010103A030029F37044957D0F79F3602000C950500000000009A032510209C01009F02060000000001005F2A020344820219809F1A0203449F03060000000000009F3303E038C89F3501219F1E0808690710522059879F090200018408A0000000250105029F4104000000019F34031F0202"
	biccdata, _ := hex.DecodeString(iccdata)
	pData := &planet_8583.ProtoStruct{
		MsgType:      "0200",
		CardNo:       "4336680006896670",
		ProcessingCd: "220000",
		Txamt:        "555",
		Syssn:        "000888",
		//PosEntryMode:         "021",
		PosEntryMode:       "072",
		Cardsequencenumber: "001",
		NetId:              "226",
		PosCondCd:          "00",
		TrackData2:         "4336680006896670D22022011193265100000",
		//ICCSystemRelatedData: "9F260825906395A2142A119F2701809F100706010A03A000009F37048F0956B49F360209AA950500000000009A032510169C01009F02060000000000015F2A020156820220009F1A0201569F3303E0F8C89F3501228407A00000000310109F0902008C9F34031F03029F1E0843415357383332309F0306000000000000",
		//ICCSystemRelatedData: "9F2608CD4C0BC6E5D0F9DA9F2701809F100706010A03A000009F3704896ADAE09F360209FF950580C00808009A032510239C01009F02060000000000015F2A02015682021C009F1A0201569F3303E0F8C89F34031F00009F3501228407A00000000310109F0902008C9F1E0843415357383332309F0306000000000000",
		ICCSystemRelatedData: biccdata,
		Tid:                  "12345678",
		MchntId:              "188000344333",
		CurrencyCd:           "344",
		//Pin:          "AA17EAB7BF18034B",
	}
	pData.Domain63Tags = make(map[string][]byte)

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

	//3bit

	_, err := ph.PackStru(ctx, pData)
	if err != nil {
		t.Fatal(err)
	}
	/*ph.PackMac(ctx, "BBEFB74400000000")
	if err != nil {
		t.Fatal(err)
	}*/
	ph.Pack(ctx)

	//bcd := hex.EncodeToString(ph.Tbuf)
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
}
