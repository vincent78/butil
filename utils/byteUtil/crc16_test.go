package byteUtil

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestCrc16(t *testing.T) {
	src := "6822000000013201060021444001010d56312e30302e303100000000000000000000000011d9"
	bs, err := hex.DecodeString(src)
	if err != nil {
		t.Errorf("src -> byte: %v", err)
	}
	bytes1 := bs[2 : len(bs)-2]
	num := CheckSum(bytes1)
	t.Logf("CRC: %X", num)
}

func TestCrc16Info(t *testing.T) {
	mData := []byte{0x01, 0x02, 0x03, 0x04}
	checksum := CheckSum(mData)
	fmt.Printf("check sum:%X \n", checksum)
	int16buf := new(bytes.Buffer)
	binary.Write(int16buf, binary.LittleEndian, checksum)
	fmt.Printf("write buf is: %+X \n", int16buf.Bytes())
	fmt.Printf("output-before:%X \n", mData)
	mData = append(mData, int16buf.Bytes()...)
	fmt.Printf("output-after:%X \n", mData)
}

func TestCrc16Bytes(t *testing.T) {
	mData, _ := hex.DecodeString("000000013201060021444001010d56312e30302e3031000000000000000000000000")
	checksum := CheckSum(mData)
	fmt.Printf("check sum:%X \n", checksum)
	int16buf := CheckSum2Byte(mData, binary.BigEndian)
	fmt.Printf("output-after:%X \n", int16buf)
}
