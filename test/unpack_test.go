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
