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

func Unpack(bcd string) map[string]string {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	// 2025-11-19
	// 200015
	//bcd := "0210303801000e8082000000000000000001002000152136071119022635333233353836373631323235383832353630303230303030303030034400059f36020a31"
	// 200016
	//bcd := "0210303801000e8082000000000000000002002000162137451119022635333233353836373632303135383832353930303230303030303030034400059f36020a31"
	// 200017
	//bcd := "0210303801000e8082000000000000000003002000172139081119022635333233353836373632373135383832363130303230303030303030034400059f36020a31"
	// 200018
	//bcd := "0210303801000e8082000000000000000030102000182200141119022635333233353836373733363035383832383930303230303030303130034400059f36020a31"
	
	b, _ := hex.DecodeString(bcd)
	uph := planet_8583.NewProtoHandler()
	ups := &planet_8583.ProtoStruct{}
	if err := uph.Unpack(ctx, b, ups); err != nil {
		panic(err)
	}
	
	return GetNonZeroFields(ups)
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

func GetTradeData(txamt, syssn, tid string) string {
	ctx := context.WithValue(context.Background(), "trace_id", util.GenXid())
	ph := planet_8583.NewProtoHandler()
	iccdata := "9F2608BB4FD0027DF6D4EC9F2701809F100706010A03A000009F37045AED67039F36020A31950580C00000009A032510279C01009F02060000000001005F2A02034482021C009F1A0203449F3303E028C89F34031E03009F3501228407A00000000310109F0902008C9F1E0843415357383332309F0306000000000000"
	biccdata, _ := hex.DecodeString(iccdata)
	pData := &planet_8583.ProtoStruct{
		MsgType:      "0200",
		CardNo:       "4336680006896670",
		ProcessingCd: "000000",
		
		// 第一笔
		//Txamt: "200",
		//Syssn: "200001",
		
		// 第二笔
		//Txamt: "1200",
		//Syssn: "200002",
		
		// 第三笔
		//Txamt: "2400",
		//Syssn: "200003",
		
		// 第四笔
		//Txamt: "1900",
		//Syssn: "200004",
		//Tid:   "20000000",
		
		// 第六笔
		//Txamt: "1910",
		//Syssn: "200005",
		//Tid:   "20000001",
		
		//Txamt: "90",
		//Syssn: "200006",
		//Tid:   "20000001",
		
		// 第7笔
		//Txamt: "9910",
		//Syssn: "200005",
		//Tid:   "20000002",
		
		// 2025-11-16
		// 第8笔
		//Txamt: "110",
		//Syssn: "200007",
		//Tid:   "20000000",
		// 第9笔
		//Txamt: "110",
		//Syssn: "200008",
		//Tid:   "20000000",
		// 第10笔
		//Txamt: "100",
		//Syssn: "200009",
		//Tid:   "20000001",
		// 第11笔
		//Txamt: "100",
		//Syssn: "200009",
		//Tid:   "20000002",
		// 12
		//Txamt: "100",
		//Syssn: "200012",
		//Tid:   "20000003",
		// 13
		//Txamt: "100",
		//Syssn: "200013",
		//Tid:   "20000004",
		// 14
		//Txamt: "100",
		//Syssn: "200014",
		//Tid:   "20000005",
		
		// 2025-11-19
		//Txamt: "100",
		//Syssn: "200015",
		//Tid:   "20000000",
		
		//Txamt: "200",
		//Syssn: "200016",
		//Tid:   "20000000",
		
		//Txamt: "300",
		//Syssn: "200017",
		//Tid:   "20000000",
		
		Txamt: txamt,
		Syssn: syssn,
		Tid:   tid,
		
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
		ProcessingCd: "960000",
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
	sock, err := getSocket()
	if err != nil {
		fmt.Printf("获取连接失败: %v\n", err)
		os.Exit(1)
	}
	defer sock.Close() // 确保连接最终关闭
	
	datas := [][3]string{
		// 20251119
		//{"3010", "200018", "20000010"},
		//{"2010", "200019", "20000010"},
		//{"2010", "200020", "20000010"},
		//{"1010", "200021", "20000010"},
		
		// 20251120
		{"10", "200022", "20000011"},
		{"10", "200023", "20000011"},
		{"10", "200024", "20000011"},
	}
	
	var allData []map[string]string
	for _, data := range datas {
		// 4. 发送报文（对应Python的sock.send）
		_, err = sock.Write(addReqHeader(GetTradeData(data[0], data[1], data[2])))
		if err != nil {
			fmt.Printf("发送数据失败: %v\n", err)
			os.Exit(1)
		}
		
		// 5. 接收响应（对应Python的sock.recv().hex()）
		// 注意：Go的recv需要指定缓冲区大小，根据实际场景调整（示例用4096）
		buf := make([]byte, 9016)
		n, err := sock.Read(buf)
		if err != nil {
			panic(err)
		}
		// 只取实际接收到的字节数，转十六进制字符串
		recvHex := hex.EncodeToString(buf[:n])
		fmt.Printf("收到响应: %s\n", recvHex[14:])
		tmp := Unpack(recvHex[14:])
		fmt.Printf("收到响应: %s\n", tmp)
		tmp["txamt"] = data[0]
		tmp["syssn"] = data[1]
		tmp["tid"] = data[2]
		allData = append(allData, tmp)
	}
	
	for _, data := range allData {
		for k, v := range data {
			fmt.Printf("%s: %s\n", k, v)
		}
	}
}

func TestTrade2(t *testing.T) {
	bcds := []string{
		"0210303801000e8082000000000000000001002000152136071119022635333233353836373631323235383832353630303230303030303030034400059f36020a31",
		"0210303801000e8082000000000000000002002000162137451119022635333233353836373632303135383832353930303230303030303030034400059f36020a31",
		"0210303801000e8082000000000000000003002000172139081119022635333233353836373632373135383832363130303230303030303030034400059f36020a31",
		"0210303801000e8082000000000000000030102000182200141119022635333233353836373733363035383832383930303230303030303130034400059f36020a31",
	}
	
	datas := []map[string]string{}
	for _, bcd := range bcds {
		datas = append(datas, Unpack(bcd))
	}
	
	for _, data := range datas {
		for k, v := range data {
			fmt.Printf("%s: %s\n", k, v)
		}
		fmt.Println("")
	}
}

func TestSettle(t *testing.T) {
	settleData := []string{
		//GetSettleData("400014", "20000010", NewD63BatchTotals(4, 8040, 0, 0, 0, 0)),
		//GetSettleData("400017", "20000000", NewD63BatchTotals(4, 700, 0, 0, 0, 0)),
		//GetSettleData("400018", "20000011", NewD63BatchTotals(3, 30, 0, 0, 0, 0)),
		GetSettleData("400019", "20000000", NewD63BatchTotals(4, 700, 0, 0, 0, 0)),
	}
	
	sock, err := getSocket()
	if err != nil {
		fmt.Printf("获取连接失败: %v\n", err)
		os.Exit(1)
	}
	defer sock.Close()
	
	var allData []map[string]string
	for _, data := range settleData {
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
		fmt.Printf("收到响应: %s\n", recvHex[14:])
		tmp := Unpack(recvHex[14:])
		allData = append(allData, tmp)
	}
	
	for _, data := range allData {
		for k, v := range data {
			fmt.Printf("%s: %s\n", k, v)
		}
		fmt.Println("")
		fmt.Println("")
	}
}
