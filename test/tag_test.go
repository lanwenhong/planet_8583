package planet_8583

import (
	"context"
	"encoding/hex"
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
		IndiCator: "D",
	}

	th := &planet_8583.TagHandler{}
	b, err := th.Pack(ctx, tag12)
	if err != nil {
		t.Fatal(err)
	}
	hexStr := hex.EncodeToString(b)
	logger.Debugf(ctx, "hex: %s", hexStr)
}
