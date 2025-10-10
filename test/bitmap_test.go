package planet_8583

import (
	"context"
	"testing"

	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/util"
	"github.com/lanwenhong/planet_8583/planet_8583"
)

func TestHasBit(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	bm := planet_8583.NewBitmap()
	bm.SetBitMap(ctx, "3020058020C09003")
	bm.HasDomain(ctx, 3)

	for i := 1; i <= 64; i++ {
		if bm.HasDomain(ctx, i) {
			logger.Debugf(ctx, "===================has %d", i)
		}
	}

}
