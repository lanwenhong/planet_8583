package planet_8583

import (
	"context"
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

	for _, k := range pData.Domain64TagKey {
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
	pData := &planet_8583.ProtoStruct{
		MsgType:      "0200",
		ProcessingCd: "002000",
		Txamt:        "555",
		Syssn:        "000027",
		PosEntryMode: "021",
		NetId:        "226",
		PosCondCd:    "00",
		TrackData2:   "50100439999991007D0810120000000323701",
		Tid:          "11411111",
		MchntId:      "99988802",
		CurrencyCd:   "124",
		Pin:          "AA17EAB7BF18034B",
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

	for _, k := range pData.Domain64TagKey {
		logger.Debugf(ctx, "tag: %s", k)
	}

	//3bit

	_, err := ph.PackStru(ctx, pData)
	if err != nil {
		t.Fatal(err)
	}
	ph.PackMac(ctx, "BBEFB74400000000")
	if err != nil {
		t.Fatal(err)
	}
	ph.Pack(ctx)

	//bcd := hex.EncodeToString(ph.Tbuf)
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
}
