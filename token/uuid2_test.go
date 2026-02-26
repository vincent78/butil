package token

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	fmt.Println("snowflake:", GenSonyflake())
}

func TestGenSonyflake(t *testing.T) {
	for i := 0; i < 20; i++ {
		t.Logf("id: %s", GenSonyflake())
	}
}
