package planet_8583

import (
	"context"
	"reflect"

	"github.com/lanwenhong/lgobase/logger"
)

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
	Domain1                  string            `bit:"1" lentype:"1" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	CardNo                   string            `bit:"2" lentype:"1" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"` //主账号
	ProcessingCd             string            `bit:"3" lentype:"0" len:"6" paddingSrc:"N"  align:"N" padding:"0" dl_type:"n"` //交易处理码
	Txamt                    string            `bit:"4" lentype:"0" len:"12" paddingSrc:"Y" align:"L" padding:"0" dl_type:"n"` //交易金额
	Domain5                  string            `bit:"5" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain6                  string            `bit:"6" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain7                  string            `bit:"7" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain8                  string            `bit:"8" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain9                  string            `bit:"9" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain10                 string            `bit:"10" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Syssn                    string            `bit:"11" lentype:"0" len:"6" paddingSrc:"N" align:"N" padding:"0" dl_type:"n"`
	TimeLocalTransaction     string            `bit:"12" lentype:"0" len:"6" paddingSrc:"N" align:"N" padding:"0" dl_type:"n"`
	DateLocalTransaction     string            `bit:"13" lentype:"0" len:"4" paddingSrc:"N" align:"N" padding:"0" dl_type:"n"`
	CardDatetime             string            `bit:"14" lentype:"0" len:"6" paddingSrc:"N" align:"N" padding:"0" dl_type:"n"`
	Domain15                 string            `bit:"15" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain16                 string            `bit:"16" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain17                 string            `bit:"17" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain18                 string            `bit:"18" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain19                 string            `bit:"19" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain20                 string            `bit:"20" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain21                 string            `bit:"21" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	PosEntryMode             string            `bit:"22" lentype:"0" len:"3" paddingSrc:"N" align:"L" padding:"0" dl_type:"n"`
	Cardsequencenumber       string            `bit:"23" lentype:"0" len:"3" paddingSrc:"N" align:"L" padding:"0" dl_type:"n"`
	NetId                    string            `bit:"24" lentype:"0" len:"3" paddingSrc:"N" align:"L" padding:"0" dl_type:"n"`
	PosCondCd                string            `bit:"25" lentype:"0" len:"2" paddingSrc:"N" align:"N" padding:"0" dl_type:"n"`
	Domain26                 string            `bit:"26" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain27                 string            `bit:"27" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain28                 string            `bit:"28" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain29                 string            `bit:"29" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain30                 string            `bit:"30" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain31                 string            `bit:"31" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain32                 string            `bit:"32" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain33                 string            `bit:"33" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain34                 string            `bit:"34" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	TrackData2               string            `bit:"35" lentype:"1" len:"37" paddingSrc:"N" align:"R" padding:"F" dl_type:"n"` //2磁道
	Domain36                 string            `bit:"36" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	RetrievalReferenceNumber string            `bit:"37" lentype:"0" len:"12" paddingSrc:"N" align:"N"  padding:"0" dl_type:"an"`
	AuthorizationIDResponse  string            `bit:"38" lentype:"0" len:"6" paddingSrc:"N" align:"N" padding:"F" dl_type:"an"`
	ResponseCode             string            `bit:"39" lentype:"0" len:"2" paddingSrc:"N" align:"N"  padding:"0" dl_type:"an"`
	Domain40                 string            `bit:"40" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Tid                      string            `bit:"41" lentype:"0" len:"8" paddingSrc:"N" align:"N"  padding:"0" dl_type:"an"`
	MchntId                  string            `bit:"42" lentype:"0" len:"15" paddingSrc:"Y" align:"R" padding:" " dl_type:"an"`
	Domain43                 string            `bit:"43" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain44                 string            `bit:"44" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain45                 string            `bit:"45" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain46                 string            `bit:"46" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain47                 string            `bit:"47" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain48                 string            `bit:"48" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	CurrencyCd               string            `bit:"49" lentype:"0" len:"3" paddingSrc:"N" align:"L" padding:"0" dl_type:"n"`
	Domain50                 string            `bit:"50" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain51                 string            `bit:"51" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Pin                      string            `bit:"52" lentype:"0" len:"16" paddingSrc:"N" align:"N" padding:"0" dl_type:"n"`
	Domain53                 string            `bit:"53" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain54                 string            `bit:"54" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	ICCSystemRelatedData     []byte            `bit:"55" lentype:"2" len:"19" paddingSrc:"N" align:"R" padding:"F" dl_type:"b"` // iccdata
	Domain56                 string            `bit:"56" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain57                 string            `bit:"57" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain58                 string            `bit:"58" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain59                 string            `bit:"59" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain60                 string            `bit:"60" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain61                 string            `bit:"61" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain62                 string            `bit:"62" lentype:"0" len:"19" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain63                 []byte            `bit:"63" lentype:"2" len:"160" paddingSrc:"N" align:"N" padding:"0" dl_type:"b" tags:"12,IA,IB,IC,ID,IE,IF,IG,IH,IL"`
	Domain64                 string            `bit:"64" lentype:"0" len:"16" paddingSrc:"N" align:"N" padding:"F" dl_type:"n"`
	Domain63Tags             map[string][]byte //tags
	Domain64TagKey           []string
	MsgType                  string
}

type ProtoHandler struct {
	Tbuf    []byte //msg type + bitmap + data stream
	MsgType []byte
	Bit     *Bitmap
	Bdata   []byte //data stream
}

type ProtoStructConf struct {
	TagMap map[string]map[string]string
}

var PsConf *ProtoStructConf

func NewProtoStructConf(ctx context.Context) *ProtoStructConf {
	tags := []string{
		TAG_BIT,
		TAG_LENTYPE,
		TAG_LEN,
		TAG_PADDINGSRC,
		TAG_DLTYPE,
		TAG_ALIGN,
		TAG_PADDING_C,
	}
	psc := &ProtoStructConf{
		TagMap: make(map[string]map[string]string),
	}
	ps := &ProtoStruct{}
	v_stru := reflect.ValueOf(ps).Elem()
	count := v_stru.NumField()
	for i := 0; i < count; i++ {
		t_item := v_stru.Type().Field(i)
		xv := make(map[string]string)
		bit := ""
		for _, k := range tags {
			v := t_item.Tag.Get(k)
			logger.Debugf(ctx, "v: %s", v)
			if v != "" {
				xv[k] = v
				if k == "bit" {
					bit = v
				}
			}
		}
		if bit != "" {
			psc.TagMap[bit] = xv
		}
	}
	PsConf = psc
	return psc
}
