package planet_8583

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/lanwenhong/planet_8583/planet_8583"

	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/util"
)

func TestTag12(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	tag12 := &planet_8583.Tag12{
		Len:       "0003",
		Tag:       "12",
		IndiCator: "X",
	}

	th := &planet_8583.TagHandler{}
	b, err := th.Pack(ctx, tag12)
	if err != nil {
		t.Fatal(err)
	}
	hexStr := hex.EncodeToString(b)
	logger.Debugf(ctx, "hex: %s", hexStr)
}

func TestTagIA(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	tagIA := &planet_8583.TagIA{
		Len:          "0004",
		Tag:          "IA",
		HostKeyIndex: "220",
	}

	th := &planet_8583.TagHandler{}
	b, err := th.Pack(ctx, tagIA)
	if err != nil {
		t.Fatal(err)
	}
	hexStr := hex.EncodeToString(b)
	logger.Debugf(ctx, "hex: %s", hexStr)
}

func TestTagIB(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	tagIB := &planet_8583.TagIB{
		Len:            "0006",
		Tag:            "IB",
		MacCheckDigits: "F9EA",
	}

	th := &planet_8583.TagHandler{}
	b, err := th.Pack(ctx, tagIB)
	if err != nil {
		t.Fatal(err)
	}
	hexStr := hex.EncodeToString(b)
	logger.Debugf(ctx, "hex: %s", strings.ToUpper(hexStr))
}

func TestTagIC(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	tagIC := &planet_8583.TagIC{
		Len:                  "0003",
		Tag:                  "IC",
		InteracTerminalClass: "03",
	}

	th := &planet_8583.TagHandler{}
	b, err := th.Pack(ctx, tagIC)
	if err != nil {
		t.Fatal(err)
	}
	hexStr := hex.EncodeToString(b)
	logger.Debugf(ctx, "hex: %s", strings.ToUpper(hexStr))
}

func TestTagID(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	tagID := &planet_8583.TagID{
		Len:                    "0003",
		Tag:                    "ID",
		InteracCustomerPresent: "1",
	}
	th := &planet_8583.TagHandler{}
	b, err := th.Pack(ctx, tagID)
	if err != nil {
		t.Fatal(err)
	}

	hexStr := hex.EncodeToString(b)
	logger.Debugf(ctx, "hex: %s", strings.ToUpper(hexStr))
}

func TestTagIE(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	tagIE := &planet_8583.TagIE{
		Len:                "0003",
		Tag:                "IE",
		InteracCardPresent: "0",
	}
	th := &planet_8583.TagHandler{}
	b, err := th.Pack(ctx, tagIE)
	if err != nil {
		t.Fatal(err)
	}

	hexStr := hex.EncodeToString(b)
	logger.Debugf(ctx, "hex: %s", strings.ToUpper(hexStr))
}

func TestTagIF(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	tagIF := &planet_8583.TagIF{
		Len:                          "0003",
		Tag:                          "IF",
		InteracCardCaptureCapability: "0",
	}
	th := &planet_8583.TagHandler{}
	b, err := th.Pack(ctx, tagIF)
	if err != nil {
		t.Fatal(err)
	}

	hexStr := hex.EncodeToString(b)
	logger.Debugf(ctx, "hex: %s", strings.ToUpper(hexStr))
}

func TestTagIG(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	tagIG := &planet_8583.TagIG{
		Len:               "0003",
		Tag:               "IG",
		BalanceinResponse: "0",
	}
	th := &planet_8583.TagHandler{}
	b, err := th.Pack(ctx, tagIG)
	if err != nil {
		t.Fatal(err)
	}

	hexStr := hex.EncodeToString(b)
	logger.Debugf(ctx, "hex: %s", strings.ToUpper(hexStr))
}

func TestTagIH(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	tagIH := &planet_8583.TagIH{
		Len:             "0003",
		Tag:             "IH",
		InteracSecurity: "0",
	}
	th := &planet_8583.TagHandler{}
	b, err := th.Pack(ctx, tagIH)
	if err != nil {
		t.Fatal(err)
	}

	hexStr := hex.EncodeToString(b)
	logger.Debugf(ctx, "hex: %s", strings.ToUpper(hexStr))
}

func TestTagIL(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	tagIL := &planet_8583.TagIL{
		Len:             "0010",
		Tag:             "IL",
		InteracSecurity: "0000702940000850",
	}
	th := &planet_8583.TagHandler{}
	b, err := th.Pack(ctx, tagIL)
	if err != nil {
		t.Fatal(err)
	}

	hexStr := hex.EncodeToString(b)
	logger.Debugf(ctx, "hex: %s", strings.ToUpper(hexStr))
}
