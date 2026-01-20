package planet_8583

import (
	"encoding/hex"
	"fmt"
	"testing"
)

//func TestNewD63BatchTotals(t *testing.T) {
//	d63 := NewD63BatchTotals(1, 200, 0, 0, 0, 0)
//	fmt.Println(d63.Len, d63.BatchTotals, len(d63.BatchTotals))
//}

func TestYYK(t *testing.T) {
	batchNumber := "400001"
	length := len(batchNumber)
	b, _ := hex.DecodeString(fmt.Sprintf("%04d", length))
	fmt.Println(b, length)
	//ret := append(b, []byte(batchNumber)...)
	//fmt.Println(ret)
	//fmt.Println([]byte(string(ret)))
	
	ret := string(b) + batchNumber
	fmt.Println(ret)
	fmt.Println([]byte(ret))
	fmt.Println([]byte(string([]byte{0, 6}) + "400001"))
}
