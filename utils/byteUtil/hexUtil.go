package byteUtil

import (
	"encoding/hex"
	"strings"
)

var (
	EmptyString = &hexError{"empty hex string"}
)

type hexError struct {
	msg string
}

func (h *hexError) Error() string {
	return h.msg
}

// BytesToHexString encodes bytes as a hex string.
func BytesToHexString(bytes []byte) string {
	encode := make([]byte, len(bytes)*2)
	hex.Encode(encode, bytes)
	return "0x" + string(encode)
}

// HexStringToBytes hex string as bytes
func HexStringToBytes(input string) ([]byte, error) {
	if len(input) == 0 {
		return nil, EmptyString
	}

	return hex.DecodeString(strings.Replace(input, "0x", "", -1))
}

// 工具函数: Hex字符串转字节数组
func HexToBytes(hexStr string) string {
	// 简单实现，实际应该处理各种格式
	return hex.EncodeToString([]byte(hexStr))
}

// ToHexArray creates a array of hex-string based on []byte
func ToHexArray(b [][]byte) []string {
	r := make([]string, len(b))
	for i := range b {
		r[i] = BytesToHexString(b[i])
	}
	return r
}

// Has0xPrefix validates str begins with '0x' or '0X'.
func Has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// isHex validates whether each byte is valid hexadecimal string.
func IsHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}
