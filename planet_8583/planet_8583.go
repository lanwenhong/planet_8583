package planet_8583

const (
	TAG_BIT         = "bit"
	TAG_LENTYPE     = "lentype"
	TAG_len         = "len"
	TAG_DTYPE       = "dtype"
	TAG_DLTYPE      = "dl_type"
	TAG_LEFT_ALIGN  = "l_align"
	TAG_RIGHT_ALIGN = "r_align"
	TAG_PADDING_C   = "padding"
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
	CardNo       string            `bit:"2" lentype:"1" len:"19" dtype:"0" l_align:"n", r_align:"n",padding:"", dl_type:"n"` //主账号
	ProcessingCd string            `bit:"3" lentype:"0" len:"6" dtype:"0" l_align:"n", r_align:"n",padding:"" dl_type:"n"`   //交易处理码
	Txamt        string            `bit:"4" lentype:"0" len:"12" dtype:"0" l_align:"n", r_align:"n",padding:"0" dl_type:"n"` //交易金额
	Syssn        string            `bit:"11" lentype:"0" len:"6" dtype:"0" l_align:"n", r_align:"n",padding:"" dl_type:"n"`
	CardDatetime string            `bit:"14" lentype:"0" len:"6" dtype:"0" l_align:"n", r_align:"n",padding:"" dl_type:"n"`
	PosEntryMode string            `bit:"12" lentype:"0" len:"3" dtype:"0" l_align:"n", r_align:"n",padding:"" dl_type:"n"`
	NetId        string            `bit:"24" lentype:"0" len:"3" dtype:"0" l_align:"n", r_align:"n",padding:"" dl_type:"n"`
	PosCondCd    string            `bit:"25" lentype:"0" len:"2" dtype:"0" l_align:"n", r_align:"n",padding:"" dl_type:"n"`
	TrackData2   string            `bit:"35" lentype:"1" len:"37" dtype:"0" l_align:"n", r_align:"n",padding:"" dl_type:"n"` //2磁道
	Tid          string            `bit:"41" lentype:"0" len:"8" dtype:"1" l_align:"n", r_align:"n",padding:"" dl_type:"n"`
	MchntId      string            `bit:"42" lentype:"0" len:"15" dtype:"1" l_align:"n", r_align:"n",padding:"" dl_type:"n"`
	CurrencyCd   string            `bit:"49" lentype:"0" len:"3" dtype:"1" l_align:"n", r_align:"n",padding:"0" dl_type:"n"`
	pin          string            `bit:"52" lentype:"0" len:"16" dtype:"1" l_align:"n", r_align:"n",padding:"" dl_type:"n"`
	Domain63     []byte            `bit:"52" lentype:"0" len:"16" dtype:"1" l_align:"n", r_align:"n",padding:"" dl_type:"" tags:"12,IA,IB,IC,ID,IE,IF,IG,IH,IL"`
	Mac          []byte            `bit:"64" lentype:"0" len:"8" dtype:"1" l_align:"n", r_align:"n",padding:"" dl_type:"n"`
	Tags         map[string][]byte //tags
}
