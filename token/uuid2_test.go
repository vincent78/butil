package token

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	fmt.Println("snowflake:", GenSonyflake())
}

func TestGenSonyflake(t *testing.T) {
	t.Logf("id: %s", GenSonyflake())
}
