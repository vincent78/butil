package strUtil

import (
	"fmt"
	"testing"

	"github.com/vincent78/butil/utils/byteUtil"
)

func TestBcd(t *testing.T) {
	//bcd := Hex2Byte("32010600214440")
	bcd := byteUtil.Hex2Byte("ffff1012")
	fmt.Printf("bcd: %x\n", bcd)
	number := Bcd2Number(bcd)
	fmt.Printf("bcd2number: %v\n", number)
	b := Number2bcd(number)
	fmt.Printf("number2bcd: %x\n", b)
}
