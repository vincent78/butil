package randUtil

import "testing"

func TestGetRandomNumber(t *testing.T) {
	for i := 0; i < 20; i++ {
		t.Logf("GetRandomNumber %d", GetRandomNumber())
	}
}

func TestGetRandomNumber100(t *testing.T) {
	for i := 0; i < 20; i++ {
		t.Logf("GetRandomNumber %02d", GetRandomNumber100())
		t.Logf("GetRandomNumber %03d", GetRandomNumber100())
	}
}
