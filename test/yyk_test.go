package planet_8583

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
	
	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/util"
	"github.com/lanwenhong/planet_8583/planet_8583"
)

func GetNonZeroFields(v interface{}) map[string]string {
	ret := map[string]string{}
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		fmt.Printf("值（类型：%s）: %v\n", reflect.TypeOf(v).Kind(), v)
		return ret
	}
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)
		if !fieldVal.CanInterface() || isZero(fieldVal) {
			continue
		}
		// 字段如果是 []byte，尝试转换为 BCD
		if fieldVal.Kind() == reflect.Slice && fieldVal.Type().Elem().Kind() == reflect.Uint8 {
			b := fieldVal.Interface().([]byte)
			ret[fieldType.Name] = hex.EncodeToString(b)
		} else {
			// 非 []byte 字段直接打印
			fmt.Printf("  %s: %v\n", fieldType.Name, fieldVal.Interface())
			ret[fieldType.Name] = fmt.Sprintf("%v", fieldVal.Interface())
		}
		
	}
	return ret
}

func Unpack(ctx context.Context, bcd string) *planet_8583.ProtoStruct {
	b, _ := hex.DecodeString(bcd)
	uph := planet_8583.NewProtoHandler()
	ups := &planet_8583.ProtoStruct{}
	if err := uph.Unpack(ctx, b, ups); err != nil {
		panic(err)
	}
	return ups
}

func UnpackWithMap(ctx context.Context, bcd string) map[string]string {
	return GetNonZeroFields(Unpack(ctx, bcd))
}

func PrintData(ctx context.Context, datas ...*planet_8583.ProtoStruct) {
	fmt.Println("=============================")
	for _, data := range datas {
		GetNonZeroFields(data)
		fmt.Println("=============================")
	}
}

// 获取SSL加密的TCP连接（对应原get_socket函数）
func getSocket() (net.Conn, error) {
	host := "terminal.uat.planetpayment.com"
	port := 40860
	caChainPath := "/home/yyk/data/pp/pp.crt"
	
	// 加载CA证书池（用于验证服务端证书）
	caCert, err := ioutil.ReadFile(caChainPath)
	if err != nil {
		return nil, fmt.Errorf("读取CA证书失败: %v", err)
	}
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("添加CA证书到证书池失败")
	}
	
	// 配置TLS客户端参数（对应原SSL上下文）
	tlsConfig := &tls.Config{
		RootCAs:            caCertPool, // 服务端证书验证的CA池
		InsecureSkipVerify: true,       // 对应原verify_mode=CERT_NONE（禁用服务端证书验证）
		ServerName:         host,       // SNI（Server Name Indication）
		// 如需启用双向认证，添加以下配置（对应原load_cert_chain）：
		// Certificates: []tls.Certificate{clientCert},
	}
	
	// 连接地址拼接
	addr := fmt.Sprintf("%s:%d", host, port)
	st := time.Now()
	
	// 建立TLS连接（替代原wrap_socket+connect）
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("TLS连接失败: %v", err)
	}
	
	// 统计连接耗时
	connectTime := time.Since(st).Seconds()
	fmt.Printf("成功连接到 %s（SSL双向认证通过） cost=%.6f\n", addr, connectTime)
	
	return conn, nil
}

// 拼接请求报文头（对应原add_req_header函数）
func addReqHeader(req string) []byte {
	// 1. 去除所有空白字符（同Python的strip+split逻辑）
	cleanReq := strings.Join(strings.Fields(req), "")
	
	// 2. 拼接TPDU头
	TPDU := "6000100000"
	raw := TPDU + cleanReq
	
	// 3. 计算报文长度（十六进制，4位补零）
	length := len(raw) / 2 // 每个十六进制字符占4位，2个字符=1字节
	lenHex := fmt.Sprintf("%04x", length)
	
	// 4. 拼接最终报文头+内容
	finalHex := lenHex + raw
	
	// 5. 十六进制字符串转字节数组
	data, err := hex.DecodeString(finalHex)
	if err != nil {
		panic(err)
	}
	
	return data
}

func SendData(ctx context.Context, datas []string) []*planet_8583.ProtoStruct {
	sock, err := getSocket()
	if err != nil {
		fmt.Printf("获取连接失败: %v\n", err)
		os.Exit(1)
	}
	defer sock.Close()
	
	var ret []*planet_8583.ProtoStruct
	for _, data := range datas {
		_, err = sock.Write(addReqHeader(data))
		if err != nil {
			fmt.Printf("发送数据失败: %v\n", err)
			os.Exit(1)
		}
		
		buf := make([]byte, 9016)
		n, err := sock.Read(buf)
		if err != nil {
			panic(err)
		}
		
		recvHex := hex.EncodeToString(buf[:n])
		ret = append(ret, Unpack(ctx, recvHex[14:]))
	}
	return ret
}

func GetTradeData(txamt, syssn, tid string) string {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := planet_8583.NewProtoHandler()
	iccdata := "9F2608BB4FD0027DF6D4EC9F2701809F100706010A03A000009F37045AED67039F36020A31950580C00000009A032510279C01009F02060000000001005F2A02034482021C009F1A0203449F3303E028C89F34031E03009F3501228407A00000000310109F0902008C9F1E0843415357383332309F0306000000000000"
	biccdata, _ := hex.DecodeString(iccdata)
	pData := &planet_8583.ProtoStruct{
		MsgType:      "0200",
		CardNo:       "4336680006896670",
		ProcessingCd: "000000",
		Txamt:        txamt,
		Syssn:        syssn,
		Tid:          tid,
		
		PosEntryMode:         "021",
		NetId:                "226",
		PosCondCd:            "00",
		TrackData2:           "4336680006896670D22022011193265100000",
		ICCSystemRelatedData: biccdata,
		MchntId:              "188000344333",
		CurrencyCd:           "344",
	}
	pData.Domain63Tags = make(map[string][]byte)
	
	// 12 O
	_ = ph.RegisterD63Tag(ctx, "12", pData, &planet_8583.Tag12{
		Len: "0003", Tag: "12", IndiCator: "X",
	})
	// IA O (Tag IA – Host Key Index)
	_ = ph.RegisterD63Tag(ctx, "IA", pData, &planet_8583.TagIA{
		Len: "0004", Tag: "IA", HostKeyIndex: "220",
	})
	// IB O (Tag IB – MAC Check Digits)
	_ = ph.RegisterD63Tag(ctx, "IB", pData, &planet_8583.TagIB{
		Len: "0006", Tag: "IB", MacCheckDigits: "F9EA",
	})
	// IC O (Tag IB – MAC Check Digits)
	_ = ph.RegisterD63Tag(ctx, "IC", pData, &planet_8583.TagIC{
		Len: "0003", Tag: "IC", InteracTerminalClass: "03",
	})
	// ID
	_ = ph.RegisterD63Tag(ctx, "ID", pData, &planet_8583.TagID{
		Len: "0003", Tag: "ID", InteracCustomerPresent: "1",
	})
	// IE
	_ = ph.RegisterD63Tag(ctx, "IE", pData, &planet_8583.TagIE{
		Len: "0003", Tag: "IE", InteracCardPresent: "0",
	})
	// IF
	_ = ph.RegisterD63Tag(ctx, "IF", pData, &planet_8583.TagIF{
		Len: "0003", Tag: "IF", InteracCardCaptureCapability: "0",
	})
	// IG
	_ = ph.RegisterD63Tag(ctx, "IG", pData, &planet_8583.TagIG{
		Len: "0003", Tag: "IG", BalanceinResponse: "0",
	})
	// IH
	_ = ph.RegisterD63Tag(ctx, "IH", pData, &planet_8583.TagIH{
		Len: "0003", Tag: "IH", InteracSecurity: "0",
	})
	// IL
	_ = ph.RegisterD63Tag(ctx, "IL", pData, &planet_8583.TagIL{
		Len: "0010", Tag: "IL", InteracSecurity: "0000702940000850",
	})
	// FA M
	//_ = ph.RegisterD63Tag(ctx, "FA", pData, &planet_8583.TagFA{
	//	Len: "0003", Tag: "FA", FinalAuthIndicator: "F",
	//})
	//// TC M
	//_ = ph.RegisterD63Tag(ctx, "TC", pData, &planet_8583.TagTC{
	//	Len: "0003", Tag: "TC", TerminalEntryCapabilities: "5",
	//})
	if _, err := ph.PackStru(ctx, pData); err != nil {
		panic(err)
	}
	ph.Pack(ctx)
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
	return fs
}

func GetSettleData(syssn, tid string, batchTotals []byte) string {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := planet_8583.NewProtoHandler()
	pData := &planet_8583.ProtoStruct{
		MsgType:      "0500",
		ProcessingCd: "920000",
		Syssn:        syssn,
		NetId:        "226",
		Tid:          tid,
		MchntId:      "188000344333",
		BatchNumber:  []byte(syssn),
		//BatchTotals:  NewD63BatchTotals(1, 100, 1, 100, 0, 0),
		BatchTotals: batchTotals,
	}
	pData.Domain63Tags = make(map[string][]byte)
	if _, err := ph.PackStru(ctx, pData); err != nil {
		panic(err)
	}
	if err := ph.PackMac(ctx, "BBEFB74400000000"); err != nil {
		panic(err)
	}
	ph.Pack(ctx)
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
	return fs
}

func TestTrade(t *testing.T) {
	ctx := context.Background()
	datas := []string{
		//GetTradeData("50", "200028", "20000012"),
		//GetTradeData("51", "200029", "20000012"),
		//GetTradeData("52", "200030", "20000012"),
		//GetTradeData("61", "200031", "20000012"),
		//GetTradeData("62", "200032", "20000012"),
		//GetTradeData("63", "200033", "20000012"),
		GetTradeData("60", "200034", "20000012"),
		GetTradeData("61", "200035", "20000012"),
		GetTradeData("62", "200036", "20000012"),
		GetTradeData("63", "200037", "20000012"),
	}
	PrintData(ctx, SendData(ctx, datas)...)
}

func TestAuthTrade(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := planet_8583.NewProtoHandler()
	iccdata := "9F2608BB4FD0027DF6D4EC9F2701809F100706010A03A000009F37045AED67039F36020A31950580C00000009A032510279C01009F02060000000001005F2A02034482021C009F1A0203449F3303E028C89F34031E03009F3501228407A00000000310109F0902008C9F1E0843415357383332309F0306000000000000"
	biccdata, _ := hex.DecodeString(iccdata)
	pData := &planet_8583.ProtoStruct{
		MsgType:              "0100",
		CardNo:               "4336680006896670",
		ProcessingCd:         "000000",
		Txamt:                "666",
		Syssn:                "200002",
		PosEntryMode:         "005",
		Cardsequencenumber:   "001",
		NetId:                "226",
		PosCondCd:            "00",
		ICCSystemRelatedData: biccdata,
		TrackData2:           "4336680006896670D22022011193265100000",
		Tid:                  "11111111",
		MchntId:              "188000344333",
		CurrencyCd:           "344",
	}
	pData.Domain63Tags = make(map[string][]byte)
	
	// FA M
	_ = ph.RegisterD63Tag(ctx, "FA", pData, &planet_8583.TagFA{
		Len: "0003", Tag: "FA", FinalAuthIndicator: "P",
	})
	
	if _, err := ph.PackStru(ctx, pData); err != nil {
		t.Fatal(err)
	}
	if err := ph.PackMac(ctx, "BBEFB74400000000"); err != nil {
		t.Fatal(err)
	}
	ph.Pack(ctx)
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
	
	PrintData(ctx, SendData(ctx, []string{fs})...)
}

func TestAuthTradeConfirm(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := planet_8583.NewProtoHandler()
	iccdata := "9F2608BB4FD0027DF6D4EC9F2701809F100706010A03A000009F37045AED67039F36020A31950580C00000009A032510279C01009F02060000000001005F2A02034482021C009F1A0203449F3303E028C89F34031E03009F3501228407A00000000310109F0902008C9F1E0843415357383332309F0306000000000000"
	biccdata, _ := hex.DecodeString(iccdata)
	pData := &planet_8583.ProtoStruct{
		MsgType:                  "0120",
		CardNo:                   "4336680006896670",
		ProcessingCd:             "020000",
		Txamt:                    "666",
		Syssn:                    "200010",
		TimeLocalTransaction:     time.Now().Format("150405"),
		DateLocalTransaction:     time.Now().Format("0102"),
		CardDatetime:             "9912",
		PosEntryMode:             "005",
		NetId:                    "226",
		PosCondCd:                "06",
		RetrievalReferenceNumber: "606166537543",
		AuthorizationIDResponse:  "735746",
		Tid:                      "11111111",
		MchntId:                  "188000344333",
		CurrencyCd:               "344",
		TrackData2:               "4336680006896670D22022011193265100000",
		ICCSystemRelatedData:     biccdata,
	}
	pData.Domain63Tags = make(map[string][]byte)
	
	// TC M
	_ = ph.RegisterD63Tag(ctx, "TC", pData, &planet_8583.TagTC{
		Len: "0003", Tag: "TC", TerminalEntryCapabilities: "5",
	})
	
	if _, err := ph.PackStru(ctx, pData); err != nil {
		t.Fatal(err)
	}
	if err := ph.PackMac(ctx, "BBEFB74400000000"); err != nil {
		t.Fatal(err)
	}
	ph.Pack(ctx)
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
	
	PrintData(ctx, SendData(ctx, []string{fs})...)
}

func GetBatchUploadData(txamt, syssn, rrn, tid string) string {
	//
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := planet_8583.NewProtoHandler()
	pData := &planet_8583.ProtoStruct{
		MsgType: "0320",
		// 2
		CardNo: "4336680006896670",
		// 3
		ProcessingCd: "000000",
		// 4
		Txamt: txamt,
		// 11
		Syssn: syssn,
		// 12
		TimeLocalTransaction: time.Now().Format("150405"),
		// 13
		DateLocalTransaction: time.Now().Format("0102"),
		// 14
		CardDatetime: "3012",
		// 22
		PosEntryMode: "021",
		// 24
		NetId: "226",
		// 25
		PosCondCd: "00",
		// 37
		RetrievalReferenceNumber: rrn,
		// 39
		ResponseCode: "00",
		// 41
		Tid: tid,
		// 42
		MchntId: "188000344333",
		// 49
		CurrencyCd: "344",
		// 55
	}
	if _, err := ph.PackStru(ctx, pData); err != nil {
	}
	if err := ph.PackMac(ctx, "BBEFB74400000000"); err != nil {
		panic("error")
	}
	ph.Pack(ctx)
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
	return fs
}

func GetBatchUploadData1(txamt, syssn, rrn, tid, cardNo, txTime, txDate, expireDate, posEntryMode, authCode, iccdata string) string {
	// 需要
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := planet_8583.NewProtoHandler()
	pData := &planet_8583.ProtoStruct{
		MsgType: "0320",
		// 2
		CardNo: cardNo,
		// 3
		ProcessingCd: "000000",
		// 4
		Txamt: txamt,
		// 11
		Syssn: syssn,
		// 12
		TimeLocalTransaction: txTime,
		// 13
		DateLocalTransaction: txDate,
		// 14
		CardDatetime: expireDate,
		// 22
		PosEntryMode: posEntryMode,
		// 24
		NetId: "226",
		// 25
		PosCondCd: "00",
		// 37
		RetrievalReferenceNumber: rrn,
		// 38
		AuthorizationIDResponse: "",
		// 39
		ResponseCode: "00",
		// 41
		Tid: tid,
		// 42
		MchntId: "188000344333",
		// 49
		CurrencyCd: "344",
	}
	
	if iccdata != "" {
		pData.ICCSystemRelatedData, _ = hex.DecodeString(iccdata)
	}
	
	if authCode != "" {
		pData.AuthorizationIDResponse = authCode
	}
	
	if _, err := ph.PackStru(ctx, pData); err != nil {
	}
	if err := ph.PackMac(ctx, "BBEFB74400000000"); err != nil {
		panic("error")
	}
	ph.Pack(ctx)
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
	return fs
}

func GetSettleUploadData(syssn, tid string, batchTotals []byte) string {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := planet_8583.NewProtoHandler()
	pData := &planet_8583.ProtoStruct{
		MsgType:      "0500",
		ProcessingCd: "960000",
		Syssn:        syssn,
		NetId:        "226",
		Tid:          tid,
		MchntId:      "188000344333",
		BatchNumber:  []byte(syssn),
		BatchTotals:  batchTotals,
	}
	pData.Domain63Tags = make(map[string][]byte)
	if _, err := ph.PackStru(ctx, pData); err != nil {
		panic(err)
	}
	if err := ph.PackMac(ctx, "BBEFB74400000000"); err != nil {
		panic(err)
	}
	ph.Pack(ctx)
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
	return fs
}

func GenSyssn(ctx context.Context) string {
	timestamp := time.Now().Unix()
	tsStr := strconv.FormatInt(timestamp, 10)
	last8Str := tsStr
	if len(tsStr) > 8 {
		last8Str = tsStr[len(tsStr)-8:] // 长度>8，取后8位
	}
	return fmt.Sprintf("%08s", last8Str)
}

func Settle(ctx context.Context, tid string, txcnt, txamt, refundcnt, refundamt, timeoutCnt, adjustcnt int) {
	syssn := GenSyssn(ctx)
	datas := []string{
		GetSettleData(syssn, tid, NewD63BatchTotals(txcnt, txamt, refundcnt, refundamt, timeoutCnt, adjustcnt)),
	}
	PrintData(ctx, SendData(ctx, datas)...)
}

func Clear(ctx context.Context, tid string) {
	syssn := GenSyssn(ctx)
	datas := []string{
		GetSettleData(syssn, tid, NewD63BatchTotals(0, 0, 0, 0, 0, 0)),
		GetSettleUploadData(syssn, tid, NewD63BatchTotals(0, 0, 0, 0, 0, 0)),
	}
	PrintData(ctx, SendData(ctx, datas)...)
}

func TestSettle(t *testing.T) {
	ctx := context.Background()
	datas := []string{
		//GetTradeData("60", "200038", "20000012"),
		//GetTradeData("60", "200039", "20000012"),
		//GetTradeData("60", "200040", "20000012"),
		//GetSettleData("400044", "87654321", NewD63BatchTotals(11, 3436000, 0, 0, 1, 0)),
		GetSettleData("100012", "99998888", NewD63BatchTotals(1, 300500, 0, 0, 0, 1)),
		//GetSettleData("100007", "99998888", NewD63BatchTotals(0, 0, 0, 0, 0, 0)),
		GetSettleUploadData("100007", "99998888", NewD63BatchTotals(0, 0, 0, 0, 0, 0)),
		GetBatchUploadData1(
			"155000", "001284", "607567767946",
			"87654321", "4514617557672096",
			"150429", "0316",
			"2902",
			"072",
			"752224",
			"9F26084D6B62E6AB3B5AD09F2701809F100706011203A000009F3704A4B4BADE9F360202D5950500000000009A032603169C01009F02060000001550005F2A020344820220209F1A0203449F3303E0B8C89F3501228407A00000000310109F0902008C9F6E04207000009B0200009F34031E03009F1E0843415338303136329F03060000000000005F340101",
		),
		GetSettleUploadData("400038", "87654321", NewD63BatchTotals(0, 0, 0, 0, 0, 0)),
	}
	PrintData(ctx, SendData(ctx, datas)...)
}
