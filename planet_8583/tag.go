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
	IndiCator string `len:"2" idl_type:"an"`
}

type TagIA struct {
	Len          string `len:"4" idl_type:"n"` //0004
	Tag          string `len:"2" idl_type:"an"`
	HostKeyIndex string `len:"3" idl_type:"n" padding:"0"`
}

type TagIB struct {
	Len            string `len:"4" idl_type:"n"` //0006
	Tag            string `len:"2" idl_type:"an"`
	MacCheckDigits string `len:"4" idl_type:"n"`
}

type TagIC struct {
	Len                  string `len:"4" idl_type:"n"` //0003
	Tag                  string `len:"2" idl_type:"an"`
	InteracTerminalClass string `len:"2" idl_type:"n"`
}

type TagID struct {
	Len                    string `len:"4" idl_type:"n"` //0003
	Tag                    string `len:"2" idl_type:"an"`
	InteracCustomerPresent string `len:"1" idl_type:"n"`
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
