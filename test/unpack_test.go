package planet_8583

import (
	"context"
	"encoding/hex"
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

func TestTransactionRsp(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	bcd := "0210203801000A800003002000000027110313052602263931343630303036313339324B4D3131343131313131004600044941011000064942363245300006494D3338384100104B4D7195A765341CB9FC00104B50C9F82E5C9A3051C6D0F9634E00000000"

	b, _ := hex.DecodeString(bcd)
	uph := planet_8583.NewProtoHandler()
	ups := &planet_8583.ProtoStruct{}
	err := uph.Unpack(ctx, b, ups)
	if err != nil {
		t.Fatal(err)
	}
	d63 := planet_8583.FormatByte(ctx, ups.Domain63)
	logger.Debugf(ctx, "d63: %s", d63)
	logger.Debugf(ctx, "ProcessingCd: %s", ups.ProcessingCd)
	logger.Debugf(ctx, "syssn: %s", ups.Syssn)
	logger.Debugf(ctx, "TimeLocalTransaction: %s", ups.TimeLocalTransaction)
	logger.Debugf(ctx, "DateLocalTransaction: %s", ups.DateLocalTransaction)
	logger.Debugf(ctx, "NetId: %s", ups.NetId)
	logger.Debugf(ctx, "RetrievalReferenceNumber: %s", ups.RetrievalReferenceNumber)
	logger.Debugf(ctx, "ResponseCode: %s", ups.ResponseCode)
	logger.Debugf(ctx, "Tid: %s", ups.Tid)

	uTagIA := &planet_8583.TagIA{}
	uTagIB := &planet_8583.TagIB{}
	uTagIM := &planet_8583.TagIM{}
	uTagKM8 := &planet_8583.TagKM8{}
	uTagKP8 := &planet_8583.TagKP8{}

	uth := &planet_8583.TagHandler{}
	uth.UnpackFromPStru(ctx, "IA", uTagIA, ups)
	logger.Debugf(ctx, "IA: %v", uTagIA)
	uth.UnpackFromPStru(ctx, "IB", uTagIB, ups)
	logger.Debugf(ctx, "IB: %v", uTagIB)
	uth.UnpackFromPStru(ctx, "IM", uTagIM, ups)
	logger.Debugf(ctx, "IM: %v", uTagIM)
	uth.UnpackFromPStru(ctx, "KM", uTagKM8, ups)
	logger.Debugf(ctx, "KM: %v", uTagKM8)
	uth.UnpackFromPStru(ctx, "KP", uTagKP8, ups)
	logger.Debugf(ctx, "KP: %v", uTagKP8)

	logger.Debugf(ctx, "mac: %s", ups.Domain64)

}

func TestKeyExchangeRsp(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	bcd := "08102038010002800002A000000004811623080513022630303131343131313131004600044941011000064942453343440006494D4630324100104B4D13A148F84F8B0D0D00104B50ECBC3257CEC7639B"
	//bcd := "0810203801000a8000009200000004811325091022022635323935353634353238303133303131343131313131"
	b, _ := hex.DecodeString(bcd)
	uph := planet_8583.NewProtoHandler()
	ups := &planet_8583.ProtoStruct{}
	err := uph.Unpack(ctx, b, ups)
	if err != nil {
		t.Fatal(err)
	}
	d63 := planet_8583.FormatByte(ctx, ups.Domain63)
	logger.Debugf(ctx, "d63: %s", d63)
	logger.Debugf(ctx, "ProcessingCd: %s", ups.ProcessingCd)
	logger.Debugf(ctx, "syssn: %s", ups.Syssn)
	logger.Debugf(ctx, "TimeLocalTransaction: %s", ups.TimeLocalTransaction)
	logger.Debugf(ctx, "DateLocalTransaction: %s", ups.DateLocalTransaction)
	logger.Debugf(ctx, "NetId: %s", ups.NetId)
	logger.Debugf(ctx, "ResponseCode: %s", ups.ResponseCode)

	uTagIA := &planet_8583.TagIA{}
	uTagIB := &planet_8583.TagIB{}
	uTagIM := &planet_8583.TagIM{}
	uTagKM8 := &planet_8583.TagKM8{}
	uTagKP8 := &planet_8583.TagKP8{}

	uth := &planet_8583.TagHandler{}
	uth.UnpackFromPStru(ctx, "IA", uTagIA, ups)
	logger.Debugf(ctx, "IA: %v", uTagIA)
	uth.UnpackFromPStru(ctx, "IB", uTagIB, ups)
	logger.Debugf(ctx, "IB: %v", uTagIB)
	uth.UnpackFromPStru(ctx, "IM", uTagIM, ups)
	logger.Debugf(ctx, "IM: %v", uTagIM)
	uth.UnpackFromPStru(ctx, "KM", uTagKM8, ups)
	logger.Debugf(ctx, "KM: %v", uTagKM8)
	uth.UnpackFromPStru(ctx, "KP", uTagKP8, ups)
	logger.Debugf(ctx, "KP: %v", uTagKP8)

}

func TestKeyExchangeTr31(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	bcd := "08102038010002000002920000529323130025121300523030028700044941022000064942413842450006494D413037310006494E413138300082504B42303038305030544E30304E303030304332363935423136414233383338363342384130433145433435314534313731303931453534333339434233444342463035434444434541393446443135304200824D4B42303038304D31544E30304E30303030413238304437463337353033453436364345303336383246384441453432434338303143394631453830413531463538334244304230343842333446373245350082444B42303038304430544E30304E303030303231343637463943433938323031453741373635344437364630464630414436383345443030453036314141303638413439323046394632463643303942384300034B5432"

	b, _ := hex.DecodeString(bcd)
	uph := planet_8583.NewProtoHandler()
	ups := &planet_8583.ProtoStruct{}
	err := uph.Unpack(ctx, b, ups)
	if err != nil {
		t.Fatal(err)
	}
	d63 := planet_8583.FormatByte(ctx, ups.Domain63)
	logger.Debugf(ctx, "d63: %s", d63)

	logger.Debugf(ctx, "ProcessingCd: %s", ups.ProcessingCd)
	logger.Debugf(ctx, "syssn: %s", ups.Syssn)
	logger.Debugf(ctx, "TimeLocalTransaction: %s", ups.TimeLocalTransaction)
	logger.Debugf(ctx, "DateLocalTransaction: %s", ups.DateLocalTransaction)
	logger.Debugf(ctx, "NetId: %s", ups.NetId)
	logger.Debugf(ctx, "ResponseCode: %s", ups.ResponseCode)

	uTagIA := &planet_8583.TagIA{}
	uTagIB := &planet_8583.TagIB{}
	uTagIM := &planet_8583.TagIM{}
	uTagIN := &planet_8583.TagIN{}
	uTagPK := &planet_8583.TagPK{}
	uTagMK := &planet_8583.TagMK{}
	uTagDK := &planet_8583.TagDK{}
	uTagKT := &planet_8583.TagKT{}

	uth := &planet_8583.TagHandler{}
	uth.UnpackFromPStru(ctx, "IA", uTagIA, ups)
	logger.Debugf(ctx, "IA: %v", uTagIA)
	uth.UnpackFromPStru(ctx, "IB", uTagIB, ups)
	logger.Debugf(ctx, "IB: %v", uTagIB)
	uth.UnpackFromPStru(ctx, "IM", uTagIM, ups)
	logger.Debugf(ctx, "IM: %v", uTagIM)
	uth.UnpackFromPStru(ctx, "IN", uTagIN, ups)
	logger.Debugf(ctx, "IN: %v", uTagIN)

	uth.UnpackFromPStru(ctx, "PK", uTagPK, ups)
	logger.Debugf(ctx, "PK: %v", uTagPK)

	uth.UnpackFromPStru(ctx, "MK", uTagMK, ups)
	logger.Debugf(ctx, "MK: %v", uTagMK)

	uth.UnpackFromPStru(ctx, "DK", uTagDK, ups)
	logger.Debugf(ctx, "DK: %v", uTagDK)

	uth.UnpackFromPStru(ctx, "KT", uTagKT, ups)
	logger.Debugf(ctx, "KT: %v", uTagKT)

}

func TestUnpackFromBCD(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	//bcd := "0210303801000e80800000200000000000055500093715242910220226353239353536343837393030353634333631303031323334353637380344"
	//bcd := "0210202001000a800000002000000937000020202020202020202020202039363132333435363738"
	//bcd := "0210202001000a800000002000000997000020202020202020202020202039363132333435363738"
	//bcd := "0210202001000a800000002000000888000020202020202020202020202039363132333435363738"
	// bcd := "0210303801000e80800000200000000000055500093716061910220226353239353536343933303435353634333831303031323334353637380344"
	//bcd := "0210202001000a800000002000000888000020202020202020202020202039363132333435363738"
	//bcd := "0210202001000a800000002000000888000020202020202020202020202039363132333435363738"
	//bcd := "0210202001000a800000002000000888000020202020202020202020202039363132333435363738"
	//bcd := "0210202001000a800000002000000888000020202020202020202020202039363132333435363738"
	//bcd := "0122002001000a00000000000000002020202020202020202020203936"
	//bcd := "0210202001000a00000000200000088800002020202020202020202020203936"
	//bcd := "0210202001000a800000002000000888000020202020202020202020202039363132333435363738"
	//bcd := "0210202001000a800000002000000888000020202020202020202020202039363132333435363738"
	//bcd := "0210202001000a800000002000000888000020202020202020202020202039363132333435363738"
	//bcd := "0210202001000a800000002000000888000020202020202020202020202039363132333435363738"
	//bcd := "0210202001000a00000000200000000000002020202020202020202020203936"
	//bcd := "0210203801000a8000000000000008881543161103022635333037353735303939363639363132333435363738"
	//bcd := "0210303801000e8082000000000000000001000001061611201103022635333037353735323132353235373334383830303132333435363738034400059f36020a31"
	//bcd := "0210303801000e8082002200000000000001000001071619451103022635333037353735323132353220202020202030303132333435363738034400059f36020a31"
	//bcd := "0210303801000e8082000000000000000001000001111719271103022635333037353735333439363635373336373230303132333435363738034400059f36020a31"
	//bcd := "0210203801000a8000002000000001141830181103022635333037353735333834393739363132333435363738" // 无可track data, icc
	//bcd := "0210303801000e8082002000000000000001000001141833081103022635333037353735333836333136353433323130303132333435363738034400059f36020a31"
	//bcd := "0210303801000e8082000000000000000001000001111837341103022635333037353735333838333935373337393430303132333435363738034400059f36020a31"
	//bcd := "0210303801000e8082002000000000000001000001141840021103022635333037353735333839363736353433323130303132333435363738034400059f36020a31"
	bcd := "0210303801000e8082000000000000000001000001111521561104022635333038353736303030383435373436333030303132333435363738034400059f36020a31"
	b, _ := hex.DecodeString(bcd)
	uph := planet_8583.NewProtoHandler()
	ups := &planet_8583.ProtoStruct{}
	err := uph.Unpack(ctx, b, ups)
	if err != nil {
		t.Fatal(err)
	}
	logger.Debugf(ctx, "ups: %+v", ups)

}
