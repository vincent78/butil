package secureUtil

import (
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {
	strTest := "I love this beautiful world!"
	strEncrypted := "98b4fc4538115c4980a8b859ff3d27e1"
	fmt.Println(Check(strTest, strEncrypted))
}
