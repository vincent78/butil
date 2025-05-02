package secureUtil

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func Sha256(secret, data string) string {
	//secret := "mysecret"
	//data := "data"
	//fmt.Printf("Secret: %s Data: %s\n", secret, data)
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secret))
	// Write Data to it
	h.Write([]byte(data))
	// Get result and encode as hexadecimal string
	return hex.EncodeToString(h.Sum(nil))
}
