package planet_8583

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/lanwenhong/lgobase/logger"
)

type Bitmap struct {
	Data []byte
}

func NewBitmap() *Bitmap {
	b := &Bitmap{}
	b.Data = make([]byte, 8)

	return b
}

func (b *Bitmap) Packbit(ctx context.Context, num int) error {
	if num > 64 {
		logger.Warnf(ctx, "num %d not support", num)
		return errors.New(fmt.Sprintf("num %d not support", num))
	}
	index, pos := num/8, num%8
	logger.Debugf(ctx, "bit: %d index: %d pos: %d", num, index, pos)
	if index != 0 {
		index = index
		//8， 16， 24, 32，40, 64bit,字节偏移减1
		if pos == 0 {
			index = index - 1
		}
	}

	if pos != 0 {
		pos = pos - 1
	} else if pos == 0 {
		pos = 7
	}
	b.Data[index] |= 0x80 >> pos
	return nil
}

func (b *Bitmap) SetBitMap(ctx context.Context, bitmap string) error {
	var err error
	b.Data, err = hex.DecodeString(bitmap)
	return err
}

func (b *Bitmap) SetBitMapByte(ctx context.Context, bitmap []byte) error {
	b.Data = bitmap
	return nil
}

func (b *Bitmap) HasDomain(ctx context.Context, domain int) bool {
	index, pos := domain/8, (8 - domain%8)
	if domain%8 == 0 && index != 0 {
		//index = 7
		index -= 1
		pos = 0
	}
	logger.Debugf(ctx, "check bit %d index %d pos %d", domain, index, pos)
	logger.Debugf(ctx, "check byte %02X", b.Data[index])

	bit := (b.Data[index] >> pos) & 0x01
	logger.Debugf(ctx, "domain: %d bit: %d", domain, bit)

	if bit == 1 {
		return true
	}
	return false
}
