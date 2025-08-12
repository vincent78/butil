package strUtil

import (
	"testing"
)

func TestHexStr2BigInt(t *testing.T) {
	str := "00000000000000000000000000000000000000000000000000000000055d4a80"
	n, e := HexStr2BigInt(str)
	if e != nil {
		t.Error(e)
	} else {
		t.Logf("number: %v", n)
	}
}
