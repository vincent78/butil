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

func TestHexStr2BigInt1(t *testing.T) {
	tests := []struct {
		name   string
		hexStr string
	}{
		{
			name:   "test1",
			hexStr: "00000000000000000000000000000000000000000000000000000000055d4a80",
		},
		{
			name:   "test2",
			hexStr: "000000000000000000000000000000000000000000000000002ffd6130a7c5b7",
		},
		{
			name:   "test3",
			hexStr: "000000000000000000000000000000000000000000000000021710ae10e95a20",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := HexStr2BigInt(tt.hexStr)
			t.Logf("%v = %v ", tt.hexStr, got)
		})
	}
}

func TestInt642HexStr(t *testing.T) {
	tests := []struct {
		name string
		i    int64
	}{
		{
			name: "test1",
			i:    13507917775357367, // 0x2ffd6130a7c5b7
		},
		{
			name: "test1",
			i:    150607452334283296, // 0x21710ae10e95a20
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Int642HexStr(tt.i)
			t.Logf("%v -> %v", tt.i, got)
		})
	}
}
