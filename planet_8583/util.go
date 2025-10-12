package planet_8583

import (
	"context"
	"encoding/hex"
	"strings"
)

func FormatByte(ctx context.Context, data []byte) string {
	fb := []byte{}
	bcd := hex.EncodeToString(data)
	for i := 0; i < len(bcd); i++ {
		fb = append(fb, byte(bcd[i]))
		if (i+1)%2 == 0 {
			fb = append(fb, " "...)
		}
	}
	//logger.Debugf(ctx, "bcd: %s", strings.ToUpper(string(fb)))
	return strings.ToUpper(string(fb))
}
