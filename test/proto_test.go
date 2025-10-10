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
