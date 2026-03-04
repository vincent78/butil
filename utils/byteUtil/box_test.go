package byteUtil

import (
	"testing"
)

func TestBox(t *testing.T) {
	//TODO:  未验证成功

	testStr := "hello world"
	bs, err := Box(testStr)
	if err != nil {
		t.Fatal(err)
		return
	}
	s, e := Unbox(bs)
	if e != nil {
		t.Fatal(e)
	}
	t.Log(s)
}
