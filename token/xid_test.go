package token

import (
	"testing"

	"github.com/rs/xid"
)

func Test_gen_NewID(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Logf("%d : %v", i, xid.New().String())
	}
	s, _ := X.NewID()
	t.Logf("%v", s)
}
