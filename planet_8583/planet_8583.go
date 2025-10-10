package planet_8583

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/lanwenhong/lgobase/logger"
)

const (
	TAG_BIT       = "bit"
	TAG_LENTYPE   = "lentype"
	TAG_len       = "len"
	TAG_DTYPE     = "dtype"
	TAG_DLTYPE    = "dl_type"
	TAG_ALIGN     = "align"
	TAG_PADDING_C = "padding"
)

const (
	DATA_LEN_TYPE_N      = "n"
	DATA_LEN_TYPE_AN     = "an"
	DATA_LEN_TYPE_Z      = "z"
	DATA_LEN_TYPE_ANS    = "ans"
	DATA_LEN_TYpE_SHADED = "shaded"
)

const (
	FIXEDLENGTH  = iota //固定长度
	VARIABLELEN2        //2位变长
	VARIABLELEN3        //3位变长
)

type Bitmap struct {
	Data []byte
}

type ProtoStruct struct {
	CardNo       string            `bit:"2" lentype:"1" len:"19" dtype:"0" align:"N" padding:"F" dl_type:"n"` //主账号
	ProcessingCd string            `bit:"3" lentype:"0" len:"6" dtype:"0"  align:"N" padding:"" dl_type:"n"`  //交易处理码
	Txamt        string            `bit:"4" lentype:"0" len:"12" dtype:"0" align:"L" padding:"0" dl_type:"n"` //交易金额
	Syssn        string            `bit:"11" lentype:"0" len:"6" dtype:"0" align:"N" padding:"" dl_type:"n"`
	CardDatetime string            `bit:"14" lentype:"0" len:"6" dtype:"0" align:"N" padding:"" dl_type:"n"`
	PosEntryMode string            `bit:"12" lentype:"0" len:"3" dtype:"0" align:"N" padding:"" dl_type:"n"`
	NetId        string            `bit:"24" lentype:"0" len:"3" dtype:"0" align:"N" padding:"" dl_type:"n"`
	PosCondCd    string            `bit:"25" lentype:"0" len:"2" dtype:"0" align:"N" padding:"" dl_type:"n"`
	TrackData2   string            `bit:"35" lentype:"1" len:"37" dtype:"0" align:"R" padding:"F" dl_type:"n"` //2磁道
	Tid          string            `bit:"41" lentype:"0" len:"8" dtype:"1" align:"N"  padding:"" dl_type:"n"`
	MchntId      string            `bit:"42" lentype:"0" len:"15" dtype:"1" align:"N" padding:"" dl_type:"n"`
	CurrencyCd   string            `bit:"49" lentype:"0" len:"3" dtype:"1" align:"L" padding:"0" dl_type:"n"`
	Pin          string            `bit:"52" lentype:"0" len:"16" dtype:"1" align:"N" padding:"" dl_type:"n"`
	Domain63     []byte            `bit:"63" lentype:"0" len:"16" dtype:"1" align:"N" padding:"" dl_type:"complex" tags:"12,IA,IB,IC,ID,IE,IF,IG,IH,IL"`
	Tags         map[string][]byte //tags
}

func NewBitmap() *Bitmap {
	b := &Bitmap{}
	b.Data = make([]byte, 8)

	return b
}

func (b *Bitmap) Packbit(ctx context.Context, num uint) error {
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
	logger.Debugf(ctx, "bit: %d", bit)

	if bit == 1 {
		return true
	}
	return false
}

type ProtoHandler struct {
}

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
		bcdLen := fmt.Sprintf("%04d", slen)
		bcd, err = hex.DecodeString(bcdLen)
	default:
		logger.Warnf(ctx, "lenType: %d not support", lenType)
		err = NewProtocolError(ERR_TAG)
	}
	return bcd, err
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

func (ph *ProtoHandler) packDomain(ctx context.Context, s string, tv reflect.StructField) ([]byte, error) {
	//dlt, err := ph.getDlType(ctx, tv)
	dBuf := []byte{}
	dlt, err := ph.getTagStr(ctx, tv, TAG_DLTYPE)
	if err != nil {
		logger.Warnf(ctx, "err: %s", err.Error())
		return dBuf, err
	}
	switch dlt {
	case "n":
		dBuf, err = ph.packNType(ctx, s, tv)
	}
	return dBuf, nil
}

func (ph *ProtoHandler) Pack(ctx context.Context, pData *ProtoStruct) ([]byte, error) {
	bdata := []byte{}
	//bitmap := NewBitmap()

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
		logger.Debugf(ctx, "nbit: %d", nbit)
		switch item.Kind() {
		case reflect.String:
			logger.Debugf(ctx, "i: %d", i)
			s := item.Interface().(string)
			if s == "" {
				continue
			}
			logger.Debugf(ctx, "pack %s", s)
			bdata, err = ph.packDomain(ctx, s, t_item)
			if err != nil {
				logger.Warnf(ctx, "err: %s", err.Error())
			}
		}
	}
	return bdata, nil
}
