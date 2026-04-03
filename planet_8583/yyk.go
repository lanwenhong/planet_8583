package planet_8583

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	
	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/util"
)

func isZero(v reflect.Value) bool {
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

func NewD63BatchTotals(txCnt, txAmt, refundCnt, refundAmt, timeoutCnt, adjustCnt int) []byte {
	batchTotals := fmt.Sprintf("%03d%012d%03d%012d%03d%03d", txCnt, txAmt, refundCnt, refundAmt, timeoutCnt, adjustCnt)
	return []byte(batchTotals + strings.Repeat("0", 90-len(batchTotals)))
}

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

func Unpack(ctx context.Context, bcd string) *ProtoStruct {
	b, _ := hex.DecodeString(bcd)
	uph := NewProtoHandler()
	ups := &ProtoStruct{}
	if err := uph.Unpack(ctx, b, ups); err != nil {
		panic(err)
	}
	return ups
}

func UnpackWithMap(ctx context.Context, bcd string) map[string]string {
	return GetNonZeroFields(Unpack(ctx, bcd))
}

func PrintData(ctx context.Context, datas ...*ProtoStruct) {
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
	crtStr := `-----BEGIN CERTIFICATE-----
MIIDxTCCAq2gAwIBAgIBADANBgkqhkiG9w0BAQsFADCBgzELMAkGA1UEBhMCVVMx
EDAOBgNVBAgTB0FyaXpvbmExEzARBgNVBAcTClNjb3R0c2RhbGUxGjAYBgNVBAoT
EUdvRGFkZHkuY29tLCBJbmMuMTEwLwYDVQQDEyhHbyBEYWRkeSBSb290IENlcnRp
ZmljYXRlIEF1dGhvcml0eSAtIEcyMB4XDTA5MDkwMTAwMDAwMFoXDTM3MTIzMTIz
NTk1OVowgYMxCzAJBgNVBAYTAlVTMRAwDgYDVQQIEwdBcml6b25hMRMwEQYDVQQH
EwpTY290dHNkYWxlMRowGAYDVQQKExFHb0RhZGR5LmNvbSwgSW5jLjExMC8GA1UE
AxMoR28gRGFkZHkgUm9vdCBDZXJ0aWZpY2F0ZSBBdXRob3JpdHkgLSBHMjCCASIw
DQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAL9xYgjx+lk09xvJGKP3gElY6SKD
E6bFIEMBO4Tx5oVJnyfq9oQbTqC023CYxzIBsQU+B07u9PpPL1kwIuerGVZr4oAH
/PMWdYA5UXvl+TW2dE6pjYIT5LY/qQOD+qK+ihVqf94Lw7YZFAXK6sOoBJQ7Rnwy
DfMAZiLIjWltNowRGLfTshxgtDj6AozO091GB94KPutdfMh8+7ArU6SSYmlRJQVh
GkSBjCypQ5Yj36w6gZoOKcUcqeldHraenjAKOc7xiID7S13MMuyFYkMlNAJWJwGR
tDtwKj9useiciAF9n9T521NtYJ2/LOdYq7hfRvzOxBsDPAnrSTFcaUaz4EcCAwEA
AaNCMEAwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMCAQYwHQYDVR0OBBYE
FDqahQcQZyi27/a9BUFuIMGU2g/eMA0GCSqGSIb3DQEBCwUAA4IBAQCZ21151fmX
WWcDYfF+OwYxdS2hII5PZYe096acvNjpL9DbWu7PdIxztDhC2gV7+AJ1uP2lsdeu
9tfeE8tTEH6KRtGX+rcuKxGrkLAngPnon1rpN5+r5N9ss4UXnT3ZJE95kTXWXwTr
gIOrmgIttRD02JDHBHNA7XIloKmf7J6raBKZV8aPEjoJpL1E/QYVN8Gb5DKj7Tjo
2GTzLH4U/ALqn83/B2gX2yKQOC16jdFU8WnjXzPKej17CuPKf1855eJ1usV2GDPO
LPAvTK33sefOT6jEm0pUBsV/fdUID+Ic/n4XuKxe9tQWskMJDE32p2u0mYRlynqI
4uJEvlz36hz1
-----END CERTIFICATE-----`
	
	caCert := []byte(crtStr)
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

func SendData(ctx context.Context, datas []string) []*ProtoStruct {
	sock, err := getSocket()
	if err != nil {
		fmt.Printf("获取连接失败: %v\n", err)
		os.Exit(1)
	}
	defer sock.Close()
	
	var ret []*ProtoStruct
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
	ph := NewProtoHandler()
	iccdata := "9F2608BB4FD0027DF6D4EC9F2701809F100706010A03A000009F37045AED67039F36020A31950580C00000009A032510279C01009F02060000000001005F2A02034482021C009F1A0203449F3303E028C89F34031E03009F3501228407A00000000310109F0902008C9F1E0843415357383332309F0306000000000000"
	biccdata, _ := hex.DecodeString(iccdata)
	pData := &ProtoStruct{
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
	_ = ph.RegisterD63Tag(ctx, "12", pData, &Tag12{
		Len: "0003", Tag: "12", IndiCator: "X",
	})
	// IA O (Tag IA – Host Key Index)
	_ = ph.RegisterD63Tag(ctx, "IA", pData, &TagIA{
		Len: "0004", Tag: "IA", HostKeyIndex: "220",
	})
	// IB O (Tag IB – MAC Check Digits)
	_ = ph.RegisterD63Tag(ctx, "IB", pData, &TagIB{
		Len: "0006", Tag: "IB", MacCheckDigits: "F9EA",
	})
	// IC O (Tag IB – MAC Check Digits)
	_ = ph.RegisterD63Tag(ctx, "IC", pData, &TagIC{
		Len: "0003", Tag: "IC", InteracTerminalClass: "03",
	})
	// ID
	_ = ph.RegisterD63Tag(ctx, "ID", pData, &TagID{
		Len: "0003", Tag: "ID", InteracCustomerPresent: "1",
	})
	// IE
	_ = ph.RegisterD63Tag(ctx, "IE", pData, &TagIE{
		Len: "0003", Tag: "IE", InteracCardPresent: "0",
	})
	// IF
	_ = ph.RegisterD63Tag(ctx, "IF", pData, &TagIF{
		Len: "0003", Tag: "IF", InteracCardCaptureCapability: "0",
	})
	// IG
	_ = ph.RegisterD63Tag(ctx, "IG", pData, &TagIG{
		Len: "0003", Tag: "IG", BalanceinResponse: "0",
	})
	// IH
	_ = ph.RegisterD63Tag(ctx, "IH", pData, &TagIH{
		Len: "0003", Tag: "IH", InteracSecurity: "0",
	})
	// IL
	_ = ph.RegisterD63Tag(ctx, "IL", pData, &TagIL{
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
	fs := FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
	return fs
}

func GetSettleData(syssn, tid string, batchTotals []byte) string {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := NewProtoHandler()
	pData := &ProtoStruct{
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
	fs := FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
	return fs
}

func GetBatchUploadData(txamt, syssn, rrn, tid string) string {
	//
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := NewProtoHandler()
	pData := &ProtoStruct{
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
	fs := FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
	return fs
}

func GetBatchUploadData1(txamt, syssn, rrn, tid, cardNo, txTime, txDate, expireDate, posEntryMode, authCode, iccdata string) string {
	// 需要
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := NewProtoHandler()
	pData := &ProtoStruct{
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
	fs := FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
	return fs
}

func GetSettleUploadData(syssn, tid string, batchTotals []byte) string {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := NewProtoHandler()
	pData := &ProtoStruct{
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
	fs := FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
	return fs
}

func GenSyssn(ctx context.Context) string {
	timestamp := time.Now().Unix()
	tsStr := strconv.FormatInt(timestamp, 10)
	last8Str := tsStr
	if len(tsStr) > 6 {
		last8Str = tsStr[len(tsStr)-6:] // 长度>8，取后8位
	}
	return fmt.Sprintf("%06s", last8Str)
}

func Settle(ctx context.Context, tid string, txcnt, txamt, refundcnt, refundamt, timeoutCnt, adjustcnt int) {
	syssn := GenSyssn(ctx)
	fmt.Printf("tid=%s, syssn=%s, txcnt=%d, txamt=%d, refundcnt=%d, refundamt=%d, timeoutCnt=%d, adjustcnt=%d\n", tid, syssn, txcnt, txamt, refundcnt, refundamt, timeoutCnt, adjustcnt)
	datas := []string{
		GetSettleData(syssn, tid, NewD63BatchTotals(txcnt, txamt, refundcnt, refundamt, timeoutCnt, adjustcnt)),
	}
	fmt.Println("执行结果")
	PrintData(ctx, SendData(ctx, datas)...)
}

func Clear(ctx context.Context, tid string) {
	syssn := GenSyssn(ctx)
	fmt.Printf("tid=%s, syssn=%s\n", tid, syssn)
	datas := []string{
		GetSettleData(syssn, tid, NewD63BatchTotals(0, 0, 0, 0, 0, 0)),
		GetSettleUploadData(syssn, tid, NewD63BatchTotals(0, 0, 0, 0, 0, 0)),
	}
	fmt.Println("执行结果")
	PrintData(ctx, SendData(ctx, datas)...)
}
