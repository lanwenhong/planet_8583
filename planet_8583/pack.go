package planet_8583

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/lanwenhong/lgobase/logger"
)

func (ph *ProtoHandler) getTagInt(ctx context.Context, tv reflect.StructField, tagName string) (int, error) {
	v := tv.Tag.Get(tagName)
	if v == "" {
		//return 0, NewProtocolError(ERR_TAG_NOTFOUND)
		return 0, NewProtocolTagNotFoundErr()
	}
	nv, err := strconv.Atoi(v)
	if err != nil {
		logger.Warnf(ctx, "err: %s", err.Error())
	}
	return nv, err

}

func (ph *ProtoHandler) getTagStr(ctx context.Context, tv reflect.StructField, tagName string) (string, error) {
	v := tv.Tag.Get(tagName)
	if v == "" {
		//logger.Warnf(ctx, "k: %s not found", tagName)
		logger.Debugf(ctx, "k: %s not found", tagName)
		return "", NewProtocolError(ERR_TAG_NOTFOUND)
	}
	return v, nil
}

func (ph *ProtoHandler) packLen(ctx context.Context, slen int, tv reflect.StructField) ([]byte, error) {
	bcd := []byte{}
	var err error = nil
	lenType, err := ph.getTagInt(ctx, tv, TAG_LENTYPE)
	if err != nil {
		logger.Warnf(ctx, "k: %s not found", TAG_LENTYPE)
		return bcd, err
	}
	switch lenType {
	case FIXEDLENGTH:
		return bcd, nil
	case VARIABLELEN2:
		bcdLen := fmt.Sprintf("%02d", slen)
		bcd, err = hex.DecodeString(bcdLen)
	case VARIABLELEN3:
		bcdLen := fmt.Sprintf("%04d", slen/2)
		bcd, err = hex.DecodeString(bcdLen)
	default:
		logger.Warnf(ctx, "lenType: %d not support", lenType)
		err = NewProtocolError(ERR_TAG)
	}
	return bcd, err
}

func (ph *ProtoHandler) packANType(ctx context.Context, s string, tv reflect.StructField) ([]byte, error) {
	return []byte(s), nil
}

func (ph *ProtoHandler) packNType(ctx context.Context, s string, tv reflect.StructField) ([]byte, error) {
	dataBuf := []byte{}
	totalBuf := []byte{}
	tagAlign, err := ph.getTagStr(ctx, tv, TAG_ALIGN)
	if err != nil {
		logger.Warnf(ctx, "err: %s", err.Error())
		return dataBuf, err
	}
	slen := len(s)
	num := slen % 2
	logger.Debugf(ctx, "tag: %s slen: %d num: %d", tagAlign, slen, num)

	var tagPadding string
	switch tagAlign {
	case "N":
		dataBuf, err = hex.DecodeString(s)
	case "L":
		bcd := []byte{}
		tagPadding, err = ph.getTagStr(ctx, tv, TAG_PADDING_C)
		if err == nil {
			for i := 0; i < num; i++ {
				bcd = append(bcd, tagPadding...)
			}
			bcd = append(bcd, s...)
			dataBuf, err = hex.DecodeString(string(bcd))
		}
	case "R":
		logger.Debugf(ctx, "padding R")
		bcd := []byte{}
		tagPadding, err = ph.getTagStr(ctx, tv, TAG_PADDING_C)
		logger.Debugf(ctx, "tagPadding: %s", tagPadding)
		if err == nil {
			bcd = append(bcd, s...)
			for i := 0; i < num; i++ {
				logger.Debugf(ctx, "appending")
				bcd = append(bcd, tagPadding...)
			}
			logger.Debugf(ctx, "bcd: %s", bcd)
			dataBuf, err = hex.DecodeString(string(bcd))
		}
	default:
		logger.Warnf(ctx, "tagAlign: %s not support")
		return dataBuf, NewProtocolError(ERR_TAG)
	}

	if err != nil {
		logger.Warnf(ctx, "err: %s", err.Error())
		return totalBuf, err
	}
	bcdlen, _ := ph.packLen(ctx, slen, tv)
	if len(bcdlen) > 0 {
		totalBuf = append(totalBuf, bcdlen...)
	}
	totalBuf = append(totalBuf, dataBuf...)
	return totalBuf, err
}

func (ph *ProtoHandler) needPaddingSrc(ctx context.Context, s string, tv reflect.StructField) (string, error) {
	lt, err := ph.getTagInt(ctx, tv, TAG_LENTYPE)
	if err != nil {
		logger.Warnf(ctx, "err: %s", err.Error())
		return "", err

	}

	//变长域不参与pandding src
	if lt != FIXEDLENGTH {
		return s, nil
	}

	paddingSrc, err := ph.getTagStr(ctx, tv, TAG_PADDINGSRC)
	if err != nil {
		logger.Warnf(ctx, "err: %s", err.Error())
		return "", err
	}
	//padding
	//var err error
	dlen, err := ph.getTagInt(ctx, tv, TAG_LEN)
	logger.Debugf(ctx, "dlen: %d", dlen)
	tagAlign, err := ph.getTagStr(ctx, tv, TAG_ALIGN)
	logger.Debugf(ctx, "tagAlign: %s", tagAlign)
	paddingC, err := ph.getTagStr(ctx, tv, TAG_PADDING_C)
	logger.Debugf(ctx, "paddingC: %s", paddingC)

	if err != nil {
		logger.Warnf(ctx, "err: %s", err.Error())
		return "", nil
	}
	slen := len(s)
	if paddingSrc == "N" {
		//校验长度
		if slen != dlen {
			logger.Warnf(ctx, "slen %d dlen %d not equal", slen, dlen)
			return "", NewProtocolError(ERR_DATA_LEN)
		}
		return s, nil
	}

	if slen > dlen {
		logger.Warnf(ctx, "slen %d > dlen %d ", slen, dlen)
		return "", NewProtocolError(ERR_DATA_LEN)
	}

	var ts string = ""
	var builder strings.Builder
	switch tagAlign {
	case "L":
		for i := 0; i < dlen-slen; i++ {
			builder.WriteString(paddingC)
		}
		builder.WriteString(s)
		ts = builder.String()
	case "R":
		builder.WriteString(s)
		for i := 0; i < dlen-slen; i++ {
			builder.WriteString(paddingC)
		}
		ts = builder.String()
	default:
		logger.Warnf(ctx, "not found tagAlign")
		return "", NewProtocolError(ERR_TAG)
	}

	logger.Debugf(ctx, "ts: %s", ts)
	return ts, nil

}

func (ph *ProtoHandler) packDomainStr(ctx context.Context, s string, tv reflect.StructField) ([]byte, error) {
	//dlt, err := ph.getDlType(ctx, tv)
	dBuf := []byte{}
	dlt, err := ph.getTagStr(ctx, tv, TAG_DLTYPE)
	if err != nil {
		logger.Warnf(ctx, "err: %s", err.Error())
		return dBuf, err
	}

	vs, err := ph.needPaddingSrc(ctx, s, tv)
	if err != nil {
		logger.Warnf(ctx, "err: %s", err.Error())
		return dBuf, err
	}

	switch dlt {
	case "n":
		dBuf, err = ph.packNType(ctx, vs, tv)
	case "an":
		dBuf, _ = ph.packANType(ctx, vs, tv)
	}

	fs := FormatByte(ctx, dBuf)
	logger.Debugf(ctx, "format pack str: %s", fs)
	return dBuf, nil
}

func (ph *ProtoHandler) packDomainSlice(ctx context.Context, s []byte, tv reflect.StructField) ([]byte, error) {
	b := []byte{}
	slen := len(s)
	logger.Debugf(ctx, "slen: %d", slen)
	zlen, _ := ph.packLen(ctx, slen, tv)
	b = append(b, zlen...)
	b = append(b, s...)

	fs := FormatByte(ctx, b)
	logger.Debugf(ctx, "format pack slice: %s", fs)
	return b, nil
}

func (ph *ProtoHandler) RegisterD63Tag(ctx context.Context, tag string, pData *ProtoStruct, tagStru interface{}) error {
	th := &TagHandler{}
	b, err := th.Pack(ctx, tagStru)
	if err != nil {
		logger.Warnf(ctx, "err: %s", err.Error())
		return err
	}
	pData.Domain64TagKey = append(pData.Domain64TagKey, tag)
	if pData.Domain63Tags == nil {
		pData.Domain63Tags = make(map[string][]byte)
	}
	pData.Domain63Tags[tag] = b
	pData.Domain63 = append(pData.Domain63, b...)
	return nil

}

func (ph *ProtoHandler) PackMac(ctx context.Context, mac string) error {
	mdata, err := hex.DecodeString(mac)
	if err != nil {
		logger.Warnf(ctx, "mac format err: %s", err.Error())
		return NewProtocolError(ERR_DEFAULT)
	}
	ph.Bit.Packbit(ctx, 64)
	ph.Bdata = append(ph.Bdata, mdata...)
	return nil
}

func (ph *ProtoHandler) PackStru(ctx context.Context, pData *ProtoStruct) ([]byte, error) {
	bitmap := ph.Bit
	bdata := []byte{}

	var err error
	//pack msg type
	if len(pData.MsgType) != 4 {
		return bdata, NewProtocolError(ERR_DATA_LEN)
	}
	ph.MsgType, err = hex.DecodeString(pData.MsgType)
	if err != nil {
		logger.Warnf(ctx, "err: %s", err.Error())
		return bdata, err
	}
	v_stru := reflect.ValueOf(pData).Elem()
	count := v_stru.NumField()
	logger.Debugf(ctx, "count: %d", count)
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
		logger.Debugf(ctx, "try pack domain nbit: %d", nbit)
		switch item.Kind() {
		case reflect.String:
			s := item.Interface().(string)
			if s == "" {
				continue
			}
			logger.Debugf(ctx, "pack %s", s)
			b, err := ph.packDomainStr(ctx, s, t_item)
			if err != nil {
				logger.Warnf(ctx, "err: %s", err.Error())
				return bdata, err
			}
			/*if nbit == 23 {
				testBcd := "DE55"
				b, _ = hex.DecodeString(testBcd)
			}*/
			bdata = append(bdata, b...)
			bitmap.Packbit(ctx, nbit)
		case reflect.Slice:
			s := item.Interface().([]byte)
			if len(s) == 0 {
				continue
			}
			b, err := ph.packDomainSlice(ctx, s, t_item)
			if err != nil {
				logger.Warnf(ctx, "err: %s", err.Error())
				return bdata, err
			}
			bdata = append(bdata, b...)
			bitmap.Packbit(ctx, nbit)
		}
	}
	ph.Bdata = bdata
	return bdata, nil
}

func (ph *ProtoHandler) Pack(ctx context.Context) {
	ph.Tbuf = append(ph.Tbuf, ph.MsgType...)
	ph.Tbuf = append(ph.Tbuf, ph.Bit.Data...)
	ph.Tbuf = append(ph.Tbuf, ph.Bdata...)
}

func NewProtoHandler() *ProtoHandler {
	ph := &ProtoHandler{
		Bit: NewBitmap(),
	}
	return ph
}
