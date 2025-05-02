package byteUtil

import "testing"

func TestBitSet1(t *testing.T) {
	var a uint8 = 30
	t.Logf("BIN: %v", ShowBin(a))
	BitSet1(&a, 6)
	t.Logf("第6位设置成1: %v", ShowBin(a))
	BitSet0(&a, 2)
	t.Logf("第2位设置成0: %v", ShowBin(a))
}
