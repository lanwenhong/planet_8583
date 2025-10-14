package planet_8583

import (
	"context"
	"testing"

	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/util"
	"github.com/lanwenhong/planet_8583/planet_8583"
)

func TestConfLoad(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	psc := planet_8583.NewProtoStructConf(ctx)
	logger.Debugf(ctx, "psc: %v", psc)
}
