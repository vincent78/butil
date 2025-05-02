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
	bytes, err := hex.DecodeString(src)
	if err != nil {
		t.Errorf("src -> byte: %v", err)
	}
	bytes1 := bytes[2 : len(bytes)-2]
	num := CheckSum(bytes1)
	t.Logf("CRC: %X", num)
}

func TestCrc16Info(t *testing.T) {
	m_data := []byte{0x01, 0x02, 0x03, 0x04}
	checksum := CheckSum(m_data)
	fmt.Printf("check sum:%X \n", checksum)
	int16buf := new(bytes.Buffer)
	binary.Write(int16buf, binary.LittleEndian, checksum)
	fmt.Printf("write buf is: %+X \n", int16buf.Bytes())
	fmt.Printf("output-before:%X \n", m_data)
	m_data = append(m_data, int16buf.Bytes()...)
	fmt.Printf("output-after:%X \n", m_data)
}

func TestCrc16Bytes(t *testing.T) {
	m_data, _ := hex.DecodeString("000000013201060021444001010d56312e30302e3031000000000000000000000000")
	checksum := CheckSum(m_data)
	fmt.Printf("check sum:%X \n", checksum)
	int16buf := CheckSum2Byte(m_data, binary.BigEndian)
	fmt.Printf("output-after:%X \n", int16buf)
}
