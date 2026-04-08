package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	
	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/planet_8583/pkg/config"
	"github.com/lanwenhong/planet_8583/pkg/models"
	"github.com/lanwenhong/planet_8583/pkg/utils"
	"github.com/lanwenhong/planet_8583/planet_8583"
)

func getSocket(ctx context.Context) net.Conn {
	host := config.Conf.PlanetAddr
	port := config.Conf.PlanetPort
	caChainPath := config.Conf.PlanetCertPath
	caCert, err := os.ReadFile(caChainPath)
	utils.MustNil(err)
	caCertPool := x509.NewCertPool()
	utils.MustTrue(caCertPool.AppendCertsFromPEM(caCert), fmt.Errorf("添加CA证书到证书池失败"))
	tlsConfig := &tls.Config{
		RootCAs:            caCertPool, // 服务端证书验证的CA池
		InsecureSkipVerify: true,       // 对应原verify_mode=CERT_NONE（禁用服务端证书验证）
		ServerName:         host,       // SNI（Server Name Indication）
	}
	st := time.Now()
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), tlsConfig)
	utils.MustNil(err)
	logger.Infof(ctx, "成功连接到 s（SSL双向认证通过） cost=%.6f", time.Since(st).Seconds())
	return conn
}

func addReqHeader(req string) []byte {
	raw := "6000100000" + strings.Join(strings.Fields(req), "")
	length := len(raw) / 2 // 每个十六进制字符占4位，2个字符=1字节
	lenHex := fmt.Sprintf("%04x", length)
	finalHex := lenHex + raw
	data, _ := hex.DecodeString(finalHex)
	return data
}

func Unpack(ctx context.Context, bcd string) *planet_8583.ProtoStruct {
	b, _ := hex.DecodeString(bcd)
	uph := planet_8583.NewProtoHandler()
	ups := &planet_8583.ProtoStruct{}
	utils.MustNil(uph.Unpack(ctx, b, ups))
	return ups
}

func SendData(ctx context.Context, sock net.Conn, data string) *planet_8583.ProtoStruct {
	_, err := sock.Write(addReqHeader(data))
	utils.MustNil(err)
	buf := make([]byte, 9016)
	n, err := sock.Read(buf)
	utils.MustNil(err)
	recvHex := hex.EncodeToString(buf[:n])
	return Unpack(ctx, recvHex[14:])
}

func GetSettleData(ctx context.Context, captureData *models.CaptureData) string {
	return ProtoToString(ctx, &planet_8583.ProtoStruct{
		MsgType:      "0500",
		ProcessingCd: "920000",
		Syssn:        captureData.BatchNumber,
		NetId:        "226",
		Tid:          captureData.Tid,
		MchntId:      captureData.MchntId,
		BatchNumber:  []byte(captureData.BatchNumber),
		BatchTotals:  captureData.GetBatchTotals(),
	})
}

func GetBatchUploadData(ctx context.Context, record *models.CaptureRecord) string {
	pData := &planet_8583.ProtoStruct{
		MsgType:                  "0320",
		CardNo:                   record.CardNo,
		ProcessingCd:             "000000",
		Txamt:                    strconv.FormatInt(record.TxAmt, 10),
		Syssn:                    record.Clisn,
		TimeLocalTransaction:     record.TxTime,
		DateLocalTransaction:     record.TxDate,
		CardDatetime:             record.ExpireDate,
		PosEntryMode:             record.PosEntryMode,
		NetId:                    "226",
		PosCondCd:                "00",
		RetrievalReferenceNumber: record.RRN,
		AuthorizationIDResponse:  "",
		ResponseCode:             "00",
		Tid:                      record.Tid,
		MchntId:                  record.MchntId,
		CurrencyCd:               record.Txcurrcd,
	}
	if record.Type == "refund" {
		pData.ProcessingCd = "200000"
	}
	if record.ICCData != "" {
		pData.ICCSystemRelatedData, _ = hex.DecodeString(record.ICCData)
	}
	if record.AuthCode != "" {
		pData.AuthorizationIDResponse = record.AuthCode
	}
	return ProtoToString(ctx, pData)
}

func GetSettleUploadData(ctx context.Context, captureData *models.CaptureData) string {
	return ProtoToString(ctx, &planet_8583.ProtoStruct{
		MsgType:      "0500",
		ProcessingCd: "960000",
		Syssn:        captureData.BatchSyssn,
		NetId:        "226",
		Tid:          captureData.Tid,
		MchntId:      captureData.MchntId,
		BatchNumber:  []byte(captureData.BatchNumber),
		BatchTotals:  captureData.GetBatchTotals(),
	})
}

func ProtoToString(ctx context.Context, pData *planet_8583.ProtoStruct) string {
	ph := planet_8583.NewProtoHandler()
	pData.Domain63Tags = make(map[string][]byte)
	_, err := ph.PackStru(ctx, pData)
	utils.MustNil(err)
	ph.Pack(ctx)
	return planet_8583.FormatByte(ctx, ph.Tbuf)
}

func Settle(ctx context.Context, captureData *models.CaptureData) {
	sock := getSocket(ctx)
	defer sock.Close()
	
	logger.Infof(ctx, "[%s]发送结算数据", captureData)
	// 发送结算数据
	closeResp := SendData(ctx, sock, GetSettleData(ctx, captureData))
	if closeResp != nil && closeResp.ResponseCode == "00" {
		return
	}
	
	// 批量清算
	logger.Infof(ctx, "[%s]upload", captureData)
	for _, record := range captureData.Records {
		uploadData := GetBatchUploadData(ctx, record)
		logger.Infof(ctx, "[%s]开始upload", record)
		for i := 0; ; i++ {
			if uploadResp := SendData(ctx, sock, uploadData); uploadResp != nil && uploadResp.ResponseCode == "00" {
				logger.Infof(ctx, "[%s]upload success", record)
				break
			}
			utils.MustTrue(i < 3, errors.New("清算失败"))
		}
	}
	
	// batch upload
	logger.Infof(ctx, "[%s]batch upload", captureData)
	if resp := SendData(ctx, sock, GetSettleUploadData(ctx, captureData)); resp == nil || resp.ResponseCode != "00" {
		msg := "清算失败"
		if resp != nil {
			msg += ", responseCode=" + resp.ResponseCode
		}
		panic(errors.New(msg))
	}
}
