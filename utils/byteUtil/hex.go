package byteUtil

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
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

//func HexStr2BigInt(str string) (*big.Int, error) {
//	if byteValue, err := hex.DecodeString(str); err == nil {
//		return new(big.Int).SetBytes(byteValue), nil
//	} else {
//		return nil, err
//	}
//}

func HexStr2Int64(str string) (int64, error) {
	numberStr := strings.Replace(str, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)
	return strconv.ParseInt(numberStr, 16, 64)
}

func Int642HexStr(i int64) string {
	return fmt.Sprintf("0x%v", strconv.FormatInt(i, 16))
}
