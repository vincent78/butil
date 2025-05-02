package strUtil

import (
	"fmt"
	"testing"
)

func TestBcd(t *testing.T) {
	//bcd := Hex2Byte("32010600214440")
	bcd := Hex2Byte("ffff1012")
	fmt.Printf("bcd: %x\n", bcd)
	number := Bcd2Number(bcd)
	fmt.Printf("bcd2number: %v\n", number)
	b := Number2bcd(number)
	fmt.Printf("number2bcd: %x\n", b)
}
