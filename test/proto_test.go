package planet_8583

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/util"
	"github.com/lanwenhong/planet_8583/planet_8583"
)

func TestPackTrack2(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ps := &planet_8583.ProtoStruct{
		TrackData2: "50100439999991007D0810120000000323701",
	}

	ph := &planet_8583.ProtoHandler{}
	bdata, err := ph.Pack(ctx, ps)
	if err != nil {
		t.Fatal(err)
	}
	bcd := hex.EncodeToString(bdata)
	logger.Debugf(ctx, "bcd: %s", bcd)
}

func TestPackTxamt(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ps := &planet_8583.ProtoStruct{
		Txamt: "66",
	}
	ph := &planet_8583.ProtoHandler{}
	bdata, err := ph.Pack(ctx, ps)
	if err != nil {
		t.Fatal(err)
	}
	bcd := hex.EncodeToString(bdata)
	logger.Debugf(ctx, "bcd: %s", strings.ToUpper(bcd))
}

func TestRegisterD63Tag(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := &planet_8583.ProtoHandler{}
	pData := &planet_8583.ProtoStruct{}
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

	bdata, err := ph.Pack(ctx, pData)
	if err != nil {
		t.Fatal(err)
	}
	bcd := hex.EncodeToString(bdata)
	logger.Debugf(ctx, "bcd: %s", strings.ToUpper(bcd))

}
