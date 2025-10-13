package planet_8583

const (
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
)

type ProtoStruct struct {
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
}
