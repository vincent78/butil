package mathUtil

import "testing"

func TestRandom(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Logf("NO.%v - %v", i, Random(3))
	}
}
