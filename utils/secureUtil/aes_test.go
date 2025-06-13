package secureUtil

import (
	"fmt"
	"testing"
)

func TestAes(t *testing.T) {
	orig := "hello world"
	key := "123456781234567812345678"
	fmt.Println("原文：", orig)

	encryptCode := AesEncrypt(orig, key)
	fmt.Println("密文：", encryptCode)

	decryptCode := AesDecrypt(encryptCode, key)
	fmt.Println("解密结果：", decryptCode)
}

func TestDecrypt(t *testing.T) {
	orig := "iEwkOunnw0OEUoHUFB7M02D3NoC_UYcsUpsYUMkgQX27uZWxA3t1cnXP6Xzj6_w0eF8uOHCA-QoMerpt7Bj1-at2H5eymZcX8KFGVpWjUVY"
	key := "B3I*SM$y.yQyG4AY"

	decryptCode := AesDecrypt(orig, key)
	fmt.Println("解密结果：", decryptCode)
}
