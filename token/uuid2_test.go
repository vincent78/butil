package token

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	fmt.Println("snowflake:", GenSonyflake())
}
