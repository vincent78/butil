package strUtil

import (
	"encoding/hex"
	"math/big"
	"strconv"
)

func Hex2Str(b []byte) string {
	return hex.EncodeToString(b)
}

func Hex2Byte(str string) []byte {
	slen := len(str)
	bHex := make([]byte, len(str)/2)
	ii := 0
	for i := 0; i < len(str); i = i + 2 {
		if slen != 1 {
			ss := string(str[i]) + string(str[i+1])
			bt, _ := strconv.ParseInt(ss, 16, 32)
			bHex[ii] = byte(bt)
			ii = ii + 1
			slen = slen - 2
		}
	}
	return bHex
}

func HexStr2BigInt(str string) (*big.Int, error) {
	if byteValue, err := hex.DecodeString(str); err == nil {
		return new(big.Int).SetBytes(byteValue), nil
	} else {
		return nil, err
	}
}

//func Int2HexStr(n *big.Int) string {
//
//}
