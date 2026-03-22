package byteUtil

import (
	"testing"
)

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
