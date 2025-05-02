package gcnet

import (
	"testing"
)

func TestIp(t *testing.T) {
	ip, err := ExternalIP()
	if err != nil {
		t.Error(err)
	}

	t.Logf(ip.String())
}
