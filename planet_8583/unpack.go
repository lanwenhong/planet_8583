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

func (ph *ProtoHandler) UnpackNType(ctx context.Context, b []byte, v reflect.Value, t reflect.StructField, start *int, unparsed *int) error {
	var dlen int //当前域名的数据长度
	var rlen int
	var bStart, bEnd int
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
		rlen = clen
		dlen = clen/2 + clen%2
		logger.Debugf(ctx, "dlen: %d", dlen)
		bStart = clen % 2
		bEnd = clen%2 + clen
		logger.Debugf(ctx, "bStart: %d bEnd: %d", bStart, bEnd)
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
		bStart = 0
		bEnd = rlen

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
		bStart = 0
		bEnd = rlen
	}

	logger.Debugf(ctx, "unpack dlen: %d rlen: %d", dlen, rlen)
	if *unparsed < dlen {
		logger.Warnf(ctx, "unparsed len err: %d", *unparsed)
		return NewProtocolError(ERR_DATA_LEN)
	}
	//unpack data
	logger.Debugf(ctx, "start: %d end: %d", *start, *start+dlen)
	udata := b[*start : *start+dlen]
	bcdData := hex.EncodeToString(udata)
	//rdata := bcdData[0:rlen]
	logger.Debugf(ctx, "bStart: %d bEnd: %d", bStart, bEnd)
	rdata := bcdData[bStart:bEnd]
	logger.Debugf(ctx, "rdata: %s", rdata)
	v.SetString(strings.ToUpper(rdata))
	*start += dlen
	*unparsed -= dlen
	return nil
}

func (ph *ProtoHandler) unpackANType(ctx context.Context, b []byte, v reflect.Value, t reflect.StructField, start *int, unparsed *int) error {
	clen, err := ph.getTagInt(ctx, t, TAG_LEN)
	if err != nil {
		logger.Warnf(ctx, "err: %s", err.Error())
		return err
	}
	logger.Debugf(ctx, "clen: %d", clen)
	anData := string(b[*start : *start+clen])
	paddingSrc, _ := ph.getTagStr(ctx, t, TAG_PADDINGSRC)
	align, _ := ph.getTagStr(ctx, t, TAG_ALIGN)
	padding, _ := ph.getTagStr(ctx, t, TAG_PADDING_C)

	*start += clen
	*unparsed -= clen

	var trimData string = ""
	if paddingSrc != "" && align != "" && padding != "" && paddingSrc == "Y" {
		switch align {
		case "L":
			trimData = strings.TrimLeft(anData, padding)
		case "R":
			trimData = strings.TrimRight(anData, padding)
		default:
			trimData = anData
		}
	} else {
		trimData = anData
	}
	v.SetString(trimData)
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
		return ph.UnpackNType(ctx, b, v, t, start, unparsed)
	case "an":
		return ph.unpackANType(ctx, b, v, t, start, unparsed)

	}
	return nil
}

func (ph *ProtoHandler) unpackDomainSlice(ctx context.Context, b []byte, v reflect.Value, t reflect.StructField, start *int, unparsed *int) error {
	var dlen = 0

	lenType, err := ph.getTagInt(ctx, t, TAG_LENTYPE)
	if err != nil {
		logger.Warnf(ctx, "k: %s not found", TAG_LENTYPE)
		return err
	}
	switch lenType {
	case VARIABLELEN2:
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
		*start += 2
		*unparsed -= 2
	default:
		logger.Warnf(ctx, "len type: %d not support", lenType)
		return NewProtocolError(ERR_TAG)
	}
	logger.Debugf(ctx, "dlen: %d", dlen)
	if *unparsed < dlen {
		logger.Warnf(ctx, "data err unparsed < dlen: %d < %d", *unparsed, dlen)
		return NewProtocolError(ERR_DATA_LEN)
	}
	bData := b[*start : *start+dlen]
	v.SetBytes(bData)
	*start += dlen
	*unparsed -= dlen
	return nil
}

func (ph *ProtoHandler) UpackDomain63Tag(ctx context.Context, pData *ProtoStruct) error {
	start := 0
	if pData.Domain63Tags == nil {
		pData.Domain63Tags = make(map[string][]byte)
	}
	if pData.Domain63 != nil && len(pData.Domain63) > 2 {
		unparsed := len(pData.Domain63)
		for {
			//parse tag len 2byte
			blen := pData.Domain63[start : start+2]
			xlen := hex.EncodeToString(blen)
			dlen, err := strconv.Atoi(xlen)
			if err != nil {
				logger.Warnf(ctx, "parse err: %s", err.Error())
				return err
			}
			if dlen+2 > unparsed {
				logger.Warnf(ctx, "data format err: dlen: %d > unparsed: %d", dlen, unparsed)
				return NewProtocolError(ERR_DATA_LEN)
			}
			//parse tag name 2byte
			tagNameStart := start + 2
			tagNameEnd := start + 4
			tagName := string(pData.Domain63[tagNameStart:tagNameEnd])

			//parse tag data
			bStart := start
			bEnd := start + dlen + 2
			bData := pData.Domain63[bStart:bEnd]
			pData.Domain63Tags[tagName] = bData
			logger.Debugf(ctx, "dlen: %d tag: %s", dlen, tagName)

			start += 2
			start += dlen
			unparsed -= 2
			unparsed -= dlen

			if start == len(pData.Domain63) {
				logger.Debugf(ctx, "max len reach")
				break
			}
		}
	}
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
				err := ph.unpackDomainSlice(ctx, realData, item, t_item, &start, &unparsed)
				if err != nil {
					logger.Debugf(ctx, "err: %s", err.Error())
					return err
				}
			}
		}
	}
	return ph.UpackDomain63Tag(ctx, pData)
	//return nil
}
