package planet_8583

import (
	"context"
	"testing"

	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/util"
	"github.com/lanwenhong/planet_8583/planet_8583"
)

func TestUnapckTrack(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ps := &planet_8583.ProtoStruct{
		MsgType:    "0200",
		TrackData2: "50100439999991007D0810120000000323701",
	}

	//ph := &planet_8583.ProtoHandler{}
	ph := planet_8583.NewProtoHandler()
	_, err := ph.PackStru(ctx, ps)
	if err != nil {
		t.Fatal(err)
	}

	ph.Pack(ctx)
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)

	uph := planet_8583.NewProtoHandler()
	ups := planet_8583.ProtoStruct{}
	err = uph.Unpack(ctx, ph.Tbuf, &ups)
	if err != nil {
		t.Fatal(err)
	}
	logger.Debugf(ctx, "track2: %s", ups.TrackData2)
	//bcd := hex.EncodeToString(bdata)
	//fs := planet_8583.FormatByte(ctx, bdata)

}

func TestUnpackMchntIdTid(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ps := &planet_8583.ProtoStruct{
		MsgType: "0200",
		Tid:     "11411111",
		MchntId: "99988802",
	}
	ph := planet_8583.NewProtoHandler()
	_, err := ph.PackStru(ctx, ps)
	if err != nil {
		t.Fatal(err)
	}

	ph.Pack(ctx)
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)

	uph := planet_8583.NewProtoHandler()
	ups := planet_8583.ProtoStruct{}
	err = uph.Unpack(ctx, ph.Tbuf, &ups)
	if err != nil {
		t.Fatal(err)
	}
	logger.Debugf(ctx, "MchntId: %s len: %d", ups.MchntId, len(ups.MchntId))
	logger.Debugf(ctx, "Tid: %s len: %d", ups.Tid, len(ups.Tid))

}

func TestUnpackDomain63(t *testing.T) {
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
	_, err := ph.PackStru(ctx, pData)
	if err != nil {
		t.Fatal(err)
	}
	ph.Pack(ctx)
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)

	//unpack
	uph := planet_8583.NewProtoHandler()
	ups := planet_8583.ProtoStruct{}
	err = uph.Unpack(ctx, ph.Tbuf, &ups)
	if err != nil {
		t.Fatal(err)
	}
	d63 := planet_8583.FormatByte(ctx, ups.Domain63)
	logger.Debugf(ctx, "d63: %s", d63)

	for k, v := range ups.Domain63Tags {
		fv := planet_8583.FormatByte(ctx, v)
		logger.Debugf(ctx, "tag: %s data: %s", k, fv)
	}

	//unparse tag12
	uTag12 := &planet_8583.Tag12{}
	uth := &planet_8583.TagHandler{}
	tagData := ups.Domain63Tags["12"]
	err = uth.Unpack(ctx, "12", uTag12, tagData)
	if err != nil {
		t.Fatal(err)
	}
	logger.Debugf(ctx, "%v", uTag12)
}
