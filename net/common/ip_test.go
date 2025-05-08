package common

import (
	"testing"
)

func TestIp(t *testing.T) {
	ip, err := ExternalIP()
	if err != nil {
		t.Error(err)
	}

	t.Logf("%v", ip.String())
}
