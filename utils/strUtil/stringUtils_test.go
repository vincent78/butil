package strUtil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatString(t *testing.T) {
	originStr := "this is the test : %v"
	valueStr := "value123"
	targetStr := fmt.Sprintf(originStr, valueStr)
	println(targetStr)
}

func TestCompareVersion(t *testing.T) {
	src := "1.0.2"
	dst := "1"
	r := CompareVersion(src, dst)
	assert.Equal(t, r, 1)

	src = "1.0.0"
	dst = "1"
	r = CompareVersion(src, dst)
	assert.Equal(t, r, 0)

	src = "1.0"
	dst = "1.0.0"
	r = CompareVersion(src, dst)
	assert.Equal(t, r, 0)

	src = "1.0.0"
	dst = "1.1.0"
	r = CompareVersion(src, dst)
	assert.Equal(t, r, 0)
}

func TestHash(t *testing.T) {
	str := "this the test"
	hash, _ := Hash(str)
	assert.Equal(t, hash, uint64(2921343956878442447))
}

func TestMd5(t *testing.T) {
	str := "this is the test"
	assert.Equal(t, Md5(str), "495faac23c0b2be85b83c140ff909ac1")
}

func TestRandomString(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Logf("%v", RandomString(9))
	}

	for i := 0; i < 10; i++ {
		t.Logf("%v", GetRandomString(10))
	}

}
