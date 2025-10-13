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

/*const (
	TAG_BIT        = "bit"
	TAG_LENTYPE    = "lentype"
	TAG_LEN        = "len"
	TAG_PADDINGSRC = "paddingSrc"
	TAG_DLTYPE     = "dl_type"
	TAG_ALIGN      = "align"
	TAG_PADDING_C  = "padding"
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
)*/

/*type Bitmap struct {
	Data []byte
}*/

/*type ProtoStruct struct {
	CardNo                   string            `bit:"2" lentype:"1" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"` //主账号
	ProcessingCd             string            `bit:"3" lentype:"0" len:"6" paddingSrc:"N"  align:"N" padding:"0" dl_type:"n"` //交易处理码
	Txamt                    string            `bit:"4" lentype:"0" len:"12" paddingSrc:"Y" align:"L" padding:"0" dl_type:"n"` //交易金额
	Syssn                    string            `bit:"11" lentype:"0" len:"6" paddingSrc:"N" align:"N" padding:"0" dl_type:"n"`
	TimeLocalTransaction     string            `bit:"12" lentype:"0" len:"6" paddingSrc:"N" align:"N" padding:"0" dl_type:"n"`
	DateLocalTransaction     string            `bit:"13" lentype:"0" len:"4" paddingSrc:"N" align:"N" padding:"0" dl_type:"n"`
	CardDatetime             string            `bit:"14" lentype:"0" len:"6" paddingSrc:"N" align:"N" padding:"0" dl_type:"n"`
	PosEntryMode             string            `bit:"22" lentype:"0" len:"3" paddingSrc:"N" align:"L" padding:"0" dl_type:"n"`
	NetId                    string            `bit:"24" lentype:"0" len:"3" paddingSrc:"N" align:"L" padding:"0" dl_type:"n"`
	PosCondCd                string            `bit:"25" lentype:"0" len:"2" paddingSrc:"N" align:"N" padding:"0" dl_type:"n"`
	TrackData2               string            `bit:"35" lentype:"1" len:"37" paddingSrc:"N" align:"R" padding:"F" dl_type:"n"` //2磁道
	RetrievalReferenceNumber string            `bit:"37" lentype:"0" len:"12" paddingSrc:"N" align:"N"  padding:"0" dl_type:"an"`
	ResponseCode             string            `bit:"39" lentype:"0" len:"2" paddingSrc:"N" align:"N"  padding:"0" dl_type:"an"`
	Tid                      string            `bit:"41" lentype:"0" len:"8" paddingSrc:"N" align:"N"  padding:"0" dl_type:"an"`
	MchntId                  string            `bit:"42" lentype:"0" len:"15" paddingSrc:"Y" align:"R" padding:" " dl_type:"an"`
	CurrencyCd               string            `bit:"49" lentype:"0" len:"3" paddingSrc:"N" align:"L" padding:"0" dl_type:"n"`
	Pin                      string            `bit:"52" lentype:"0" len:"16" paddingSrc:"N" align:"N" padding:"0" dl_type:"n"`
	Domain63                 []byte            `bit:"63" lentype:"2" len:"160" paddingSrc:"N" align:"N" padding:"0" dl_type:"complex" tags:"12,IA,IB,IC,ID,IE,IF,IG,IH,IL"`
	Domain63Tags             map[string][]byte //tags
	Domain64TagKey           []string
	MsgType                  string
}*/

/*func NewBitmap() *Bitmap {
	b := &Bitmap{}
	b.Data = make([]byte, 8)

	return b
}*/

/*func (b *Bitmap) Packbit(ctx context.Context, num int) error {
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
}*/

type ProtoHandler struct {
	Tbuf    []byte //msg type + bitmap + data stream
	MsgType []byte
	Bit     *Bitmap
	Bdata   []byte //data stream
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
