package planet_8583

import "fmt"

const (
	ERR_BIT          = 1
	ERR_TAG          = 2
	ERR_TAG_NOTFOUND = 3
	ERR_DATA_TYPE    = 4
	ERR_DATA_LEN     = 6
	ERR_MAC_FORMAT   = 7
	ERR_TAG63        = 8
	ERR_DEFAULT      = 9
)

var ErrMap map[int]string = map[int]string{
	ERR_BIT:          "位图设置错误",
	ERR_TAG:          "结构体tag错误",
	ERR_TAG_NOTFOUND: "tag未找到",
	ERR_DATA_TYPE:    "数据类型错误",
	ERR_DATA_LEN:     "数据长度错误",
	ERR_MAC_FORMAT:   "mac格式错误",
	ERR_TAG63:        "63域tag错误",
}

type ProtocolError struct {
	Code int
	Msg  string
}

type ProtocolTagNotFoundErr struct {
	ProtocolError
}

func (pe *ProtocolError) SetCode(code int) {
	if msg, ok := ErrMap[code]; ok {
		pe.Msg = msg
	}
}

func (pe *ProtocolError) GetCode() int {
	return pe.Code
}

func (pe *ProtocolError) Error() string {
	return fmt.Sprintf("code: %d err: %s", pe.Code, pe.Msg)
}

func NewProtocolError(code int) *ProtocolError {
	pe := &ProtocolError{
		Code: code,
	}
	if msg, ok := ErrMap[code]; ok {
		pe.Msg = msg
	}
	return pe
}

func NewProtocolErrorDefault(msg string) *ProtocolError {
	pe := &ProtocolError{
		Code: ERR_DEFAULT,
		Msg:  msg,
	}
	return pe
}

func NewProtocolTagNotFoundErr() *ProtocolTagNotFoundErr {
	pe := &ProtocolTagNotFoundErr{
		//Code: ERR_TAG_NOTFOUND,
	}
	pe.Code = ERR_TAG_NOTFOUND
	if msg, ok := ErrMap[ERR_TAG_NOTFOUND]; ok {
		pe.Msg = msg
	}
	return pe

}
