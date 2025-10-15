package planet_8583

import (
	"context"
	"encoding/hex"
	"errors"
	"reflect"
	"strconv"

	"github.com/lanwenhong/lgobase/logger"
)

type Tag12 struct {
	Len       string `len:"4" idl_type:"n"` //0003
	Tag       string `len:"2" idl_type:"an"`
	IndiCator string `len:"1" idl_type:"an"`
}

type TagIA struct {
	Len          string `len:"4" idl_type:"n"` //0004
	Tag          string `len:"2" idl_type:"an"`
	HostKeyIndex string `len:"3" idl_type:"n" padding:"0"`
}

type TagIB struct {
	Len            string `len:"4" idl_type:"n"` //0006
	Tag            string `len:"2" idl_type:"an"`
	MacCheckDigits string `len:"4" idl_type:"an"`
}

type TagIC struct {
	Len                  string `len:"4" idl_type:"n"` //0003
	Tag                  string `len:"2" idl_type:"an"`
	InteracTerminalClass string `len:"2" idl_type:"n"`
}

type TagID struct {
	Len                    string `len:"4" idl_type:"n"` //0003
	Tag                    string `len:"2" idl_type:"an"`
	InteracCustomerPresent string `len:"1" idl_type:"n"  padding:"0"`
}

type TagIE struct {
	Len                string `len:"4" idl_type:"n"` //0003
	Tag                string `len:"2" idl_type:"an"`
	InteracCardPresent string `len:"1" idl_type:"n"`
}

type TagIF struct {
	Len                          string `len:"4" idl_type:"n"` //0003
	Tag                          string `len:"2" idl_type:"an"`
	InteracCardCaptureCapability string `len:"1" idl_type:"n"`
}

type TagIG struct {
	Len               string `len:"4" idl_type:"n"` //0003
	Tag               string `len:"2" idl_type:"an"`
	BalanceinResponse string `len:"1" idl_type:"n"`
}

type TagIH struct {
	Len             string `len:"4" idl_type:"n"` //0003
	Tag             string `len:"2" idl_type:"an"`
	InteracSecurity string `len:"1" idl_type:"n"`
}

type TagIL struct {
	Len             string `len:"4" idl_type:"n"` //0010
	Tag             string `len:"2" idl_type:"an"`
	InteracSecurity string `len:"16" idl_type:"n"`
}

type TagHandler struct {
}

func (th *TagHandler) getTagLen(tv reflect.StructField) int {
	slen := tv.Tag.Get("len")
	dlen, _ := strconv.Atoi(slen)
	return dlen

}

func (th *TagHandler) getTagLenType(ctx context.Context, tv reflect.StructField) string {
	lenType := tv.Tag.Get("idl_type")
	logger.Debugf(ctx, "lenType: %v", lenType)
	return lenType
}

func (th *TagHandler) packString(ctx context.Context, s string, dlen int, dlenType string) ([]byte, error) {
	b := []byte{}
	slen := len(s)
	logger.Debugf(ctx, "slen: %d dlen: %d dlenType: %s s: %s", slen, dlen, dlenType, s)
	if slen > dlen {
		logger.Warnf(ctx, "slen > dlen")
		return b, errors.New("slen > dlen")
	}
	switch dlenType {
	case "n":
		bcd := []byte{}
		num := slen % 2
		for i := 0; i < num; i++ {
			bcd = append(bcd, "0"...)
		}
		bcd = append(bcd, s...)
		b, _ = hex.DecodeString(string(bcd))
	case "an":
		b = append(b, s...)
	default:
		logger.Warnf(ctx, "not support")
		return b, errors.New("not support")
	}
	return b, nil
}

func (th *TagHandler) Pack(ctx context.Context, tagStru interface{}) ([]byte, error) {
	v_stru := reflect.ValueOf(tagStru).Elem()
	count := v_stru.NumField()
	logger.Debugf(ctx, "count: %d", count)

	var tagBuf = []byte{}

	for i := 0; i < count; i++ {
		item := v_stru.Field(i)
		t_item := v_stru.Type().Field(i)
		dlen := th.getTagLen(t_item)
		lenType := th.getTagLenType(ctx, t_item)
		switch item.Kind() {
		case reflect.String:
			b, err := th.packString(ctx, item.Interface().(string), dlen, lenType)
			if err == nil {
				tagBuf = append(tagBuf, b...)
			} else {
				return b, errors.New("not support")
			}
		default:
			logger.Warnf(ctx, "not support")
			return tagBuf, errors.New("not support")
		}
	}
	return tagBuf, nil
}

func (th *TagHandler) unPackTagData(ctx context.Context, tagData []byte, dlen int, dlenType string,
	v reflect.Value, start *int, unparsed *int) error {
	b := []byte{}
	//slen := len(s)
	switch dlenType {
	case "n":
		pdlen := dlen + dlen%2
		pdlen = pdlen / 2

		if *unparsed < pdlen {
			logger.Warnf(ctx, "unparsed: %d < pdlen: %d", *unparsed, pdlen)
			return NewProtocolError(ERR_DATA_LEN)
		}
		b = tagData[*start : *start+pdlen]
		sb := hex.EncodeToString(b)
		v.SetString(sb)

		*start += pdlen
		*unparsed -= pdlen
	case "an":
		if *unparsed < dlen {
			logger.Warnf(ctx, "unparsed: %d < dlen: %d", *unparsed, dlen)
			return NewProtocolError(ERR_DATA_LEN)
		}
		b = tagData[*start : *start+dlen]
		v.SetString(string(b))
		*start += dlen
		*unparsed -= dlen
	default:
		logger.Warnf(ctx, "not support")
		return errors.New("not support")
	}
	return nil
}

func (th *TagHandler) Unpack(ctx context.Context, tagName string, tagStru interface{}, tagData []byte) error {
	var err error = nil
	v_stru := reflect.ValueOf(tagStru).Elem()
	count := v_stru.NumField()
	logger.Debugf(ctx, "count: %d", count)

	start := 0
	unparsed := len(tagData)
	if unparsed < 4 {
		logger.Warnf(ctx, "data format err unparsed: %d < 4", unparsed)
		return NewProtocolError(ERR_DATA_LEN)
	}
	tagLen := 0
	for i := 0; i < count; i++ {
		item := v_stru.Field(i)
		if i == 0 {
			blen := tagData[start : start+2]
			xlen := hex.EncodeToString(blen)
			tagLen, err = strconv.Atoi(xlen)
			if err != nil {
				logger.Warnf(ctx, "parse err: %s", err.Error())
				return err
			}
			logger.Debugf(ctx, "tagLen: %d", tagLen)
			item.SetString(xlen)
			start += 2
			unparsed -= 2
		} else if i == 1 {
			tagName := string(tagData[start : start+2])
			item.SetString(tagName)
			start += 2
			unparsed -= 2
		} else {
			t_item := v_stru.Type().Field(i)
			dlen := th.getTagLen(t_item)
			lenType := th.getTagLenType(ctx, t_item)
			err = th.unPackTagData(ctx, tagData, dlen, lenType, item, &start, &unparsed)
			if err != nil {
				return err
			}
		}
		if start == len(tagData) {
			logger.Debugf(ctx, "max len reach")
			break
		}
	}
	return nil
}
