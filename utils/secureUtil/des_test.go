package secureUtil

import (
	"fmt"
	"log"
	"testing"
)

func TestDes(t *testing.T) {
	key := []byte("2fa6c1e9")
	str := "I love this beautiful world!"
	strEncrypted, err := DesEncrypt(str, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Encrypted:", strEncrypted)
	strDecrypted, err := DesDecrypt(strEncrypted, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Decrypted:", strDecrypted)

}
