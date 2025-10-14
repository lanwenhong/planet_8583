package planet_8583

import (
	"context"
	"encoding/hex"
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/lanwenhong/lgobase/logger"
)

func (ph *ProtoHandler) unpackNType(ctx context.Context, b []byte, v reflect.Value, t reflect.StructField, start *int, unparsed *int) error {
	var dlen int //当前域名的数据长度
	var rlen int
	lenType, err := ph.getTagInt(ctx, t, TAG_LENTYPE)
	if err != nil {
		logger.Warnf(ctx, "k: %s not found", TAG_LENTYPE)
		return err
	}
	switch lenType {
	case FIXEDLENGTH:
		clen, err := ph.getTagInt(ctx, t, TAG_LEN)
		if err != nil {
			logger.Warnf(ctx, "err: %s", err.Error())
			return err
		}
		logger.Debugf(ctx, "dlen: %d", clen)
		dlen = clen/2 + clen%2
		rlen = dlen
	case VARIABLELEN2:
		//2 byte
		if *unparsed < 1 {
			logger.Warnf(ctx, "unparsed len err: %d", *unparsed)
			return NewProtocolError(ERR_DATA_LEN)
		}
		s := *start
		e := *start + 1
		slen := b[s:e]
		xlen := hex.EncodeToString(slen)
		dlen, err = strconv.Atoi(xlen)
		if err != nil {
			logger.Warnf(ctx, "parse err: %s", err.Error())
			return err
		}
		rlen = dlen
		//n格式的变长数据，不够2字节，会补f，长度标记的是没有补f的长度
		dlen = dlen/2 + dlen%2
		*start += 1
		*unparsed -= 1
	case VARIABLELEN3:
		if *unparsed < 2 {
			logger.Warnf(ctx, "unparsed len err: %d", *unparsed)
			return NewProtocolError(ERR_DATA_LEN)
		}
		s := *start
		e := *start + 2
		slen := b[s:e]
		xlen := hex.EncodeToString(slen)
		dlen, err = strconv.Atoi(xlen)
		if err != nil {
			logger.Warnf(ctx, "parse err: %s", err.Error())
			return err
		}
		rlen = dlen
		dlen = dlen/2 + dlen%2
		*start += 2
		*unparsed -= 2
	}

	logger.Debugf(ctx, "unpack dlen: %d rlen: %d", dlen, rlen)
	//unpack data
	udata := b[*start : *start+dlen]
	bcdData := hex.EncodeToString(udata)
	rdata := bcdData[0:rlen]
	logger.Debugf(ctx, "rdata: %s", rdata)
	v.SetString(strings.ToUpper(rdata))
	return nil
}

func (ph *ProtoHandler) unpackDomainStr(ctx context.Context, b []byte, v reflect.Value, t reflect.StructField, start *int, unparsed *int) error {
	dlt, err := ph.getTagStr(ctx, t, TAG_DLTYPE)
	if err != nil {
		logger.Warnf(ctx, "err: %s", err.Error())
		return err
	}
	switch dlt {
	case "n":
		return ph.unpackNType(ctx, b, v, t, start, unparsed)
	case "an":
	}
	return nil
}

func (ph *ProtoHandler) unpackDomainSlice(ctx context.Context, b []byte, v reflect.Value, t reflect.StructField, start *int, unparsed *int) error {
	return nil
}

func (ph *ProtoHandler) Unpack(ctx context.Context, bData []byte, pData *ProtoStruct) error {

	var msgType []byte
	var bitMap []byte
	//2 byte msg type
	msgType = bData[:2]
	//8 byte bitmap
	bitMap = bData[2 : 2+8]
	//data
	realData := bData[10:]

	logger.Debugf(ctx, "msgType: %v", msgType)
	//bitmap := NewBitmap()
	ph.Bit.SetBitMapByte(ctx, bitMap)

	start := 0
	unparsed := len(bData)

	logger.Debugf(ctx, "start: %d unparsed: %d", start, unparsed)

	v_stru := reflect.ValueOf(pData).Elem()
	count := v_stru.NumField()
	for i := 0; i < count; i++ {
		item := v_stru.Field(i)
		t_item := v_stru.Type().Field(i)
		nbit, err := ph.getTagInt(ctx, t_item, TAG_BIT)
		var pErr *ProtocolTagNotFoundErr
		if errors.As(err, &pErr) {
			continue
		}
		if err != nil {
			logger.Warnf(ctx, "err: %s", err.Error())
			continue
		}
		logger.Debugf(ctx, "check bit %d", nbit)
		if ph.Bit.HasDomain(ctx, nbit) {
			switch item.Kind() {
			case reflect.String:
				err := ph.unpackDomainStr(ctx, realData, item, t_item, &start, &unparsed)
				if err != nil {
					logger.Debugf(ctx, "err: %s", err.Error())
					return err
				}
			case reflect.Slice:
			}
		}
	}
	return nil
}
