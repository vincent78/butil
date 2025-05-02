package secureUtil

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"testing"
)

// 私钥生成
// openssl genrsa -out rsa_private_key.pem 1024
var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQDcGsUIIAINHfRTdMmgGwLrjzfMNSrtgIf4EGsNaYwmC1GjF/bM
h0Mcm10oLhNrKNYCTTQVGGIxuc5heKd1gOzb7bdTnCDPPZ7oV7p1B9Pud+6zPaco
qDz2M24vHFWYY2FbIIJh8fHhKcfXNXOLovdVBE7Zy682X1+R1lRK8D+vmQIDAQAB
AoGAeWAZvz1HZExca5k/hpbeqV+0+VtobMgwMs96+U53BpO/VRzl8Cu3CpNyb7HY
64L9YQ+J5QgpPhqkgIO0dMu/0RIXsmhvr2gcxmKObcqT3JQ6S4rjHTln49I2sYTz
7JEH4TcplKjSjHyq5MhHfA+CV2/AB2BO6G8limu7SheXuvECQQDwOpZrZDeTOOBk
z1vercawd+J9ll/FZYttnrWYTI1sSF1sNfZ7dUXPyYPQFZ0LQ1bhZGmWBZ6a6wd9
R+PKlmJvAkEA6o32c/WEXxW2zeh18sOO4wqUiBYq3L3hFObhcsUAY8jfykQefW8q
yPuuL02jLIajFWd0itjvIrzWnVmoUuXydwJAXGLrvllIVkIlah+lATprkypH3Gyc
YFnxCTNkOzIVoXMjGp6WMFylgIfLPZdSUiaPnxby1FNM7987fh7Lp/m12QJAK9iL
2JNtwkSR3p305oOuAz0oFORn8MnB+KFMRaMT9pNHWk0vke0lB1sc7ZTKyvkEJW0o
eQgic9DvIYzwDUcU8wJAIkKROzuzLi9AvLnLUrSdI6998lmeYO9x7pwZPukz3era
zncjRK3pbVkv0KrKfczuJiRlZ7dUzVO0b6QJr8TRAA==
-----END RSA PRIVATE KEY-----
`)

// 公钥: 根据私钥生成
// openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem
var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDcGsUIIAINHfRTdMmgGwLrjzfM
NSrtgIf4EGsNaYwmC1GjF/bMh0Mcm10oLhNrKNYCTTQVGGIxuc5heKd1gOzb7bdT
nCDPPZ7oV7p1B9Pud+6zPacoqDz2M24vHFWYY2FbIIJh8fHhKcfXNXOLovdVBE7Z
y682X1+R1lRK8D+vmQIDAQAB
-----END PUBLIC KEY-----
`)

func TestRsa(t *testing.T) {
	data, _ := RsaEncrypt([]byte("hello world"), publicKey)
	fmt.Println(base64.StdEncoding.EncodeToString(data))
	origData, _ := RsaDecrypt(data, privateKey)
	fmt.Println(string(origData))
}

func TestCreateKeys(t *testing.T) {
	KeyPairs("/tmp/test1")
}

func TestRsaSign(t *testing.T) {
	message := []byte("Hello, world!")
	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Error(err)
	}

	hash := sha256.Sum256(message)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hash[:])
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Message: ", string(message))
	fmt.Println("Signature: ", base64.StdEncoding.EncodeToString(signature))

	pubKey := privKey.PublicKey
	err = rsa.VerifyPKCS1v15(&pubKey, crypto.SHA256, hash[:], signature)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Signature verified")
}

func TestRsaLoadAndSign(t *testing.T) {
	private64 := "MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQCY/i3wvgeqo1gSusHf1AuPYU+nRMNLkZeIPpA6YjvE5iQ26DQXdrQ0WamsvAKz3GEQMwRqmtesEQhLygJj4wJdO4PDhGNqlObsWQx4N1cyrbvGouTbO2hUhnesP3PURMRBEpJd1A/koKmh606i4oKcXKuVBwjoBMdPPCMZ4QchcJ30KFLtiXT9OBpGim7SgFBiKnWxO/4CGd1KkixRj4ID0lyYhRnnCUvFptN822V4g5yr7vSUsH3M7IxVO4SgFzFZjM9pvuo0G2vf51NJ+fK3y9rK24/Vt671sEf59s28OXuyRXErjgfIoDdM2triJ2Fq8jZQV9kVvR6gn7rmZfUnAgMBAAECggEAQ248J1BKJr5Jsi+YBaP62F4Gcm3POb5YsFcK0IC9YSIiMgUT+Id8E1q1ewl+k3F9YltqBeZrSk5TfrvxY78JKrhxcbom6zHnuaHh6hZSG2cRTRI8lhfP+vktQ8DPt237pcaetjYiLx1UxqXkicwVzv7VLSDlnwWEJvsVaXGR5/2BT8q+/2VEK4qCe8DESNpWNlDfonXAK0FDtDzWkjwLeWzJtzWQLw0ps8gSTQsUYRA2GUBtcp3MWOy+GOAIzhbTawOYTi3EjvAsRB7YuLYLOnueid0vYVRu6IHETcOBJIpGbBxV0IpbNvYNJ53A1bgyELvKIM9xUYs/3m5HIc6mmQKBgQDKvaYx1nUTjikkkn88IC+TgnMGSBSDSKcxZd0IOUPC1ohtnB0x/IcH+mEBot8GEkn7CjnyvtbaBq3I+RxGpfrLzdzt8BH9tLxyrGA872iXfB8owRMpoOs0hDMRTT2gZPsXpNdQwDP4UvLqmsQPw0QO5id7gLnc5Rm6OSMcIaXx/QKBgQDBLvl3we7ZzN+PydXBY9AKnvAl9BeDFPgynsRXn0dYNuKDWR3PXF/IOLGraa7LHZ3L6WJY3fLRr5CMV+k8RjWo6aZMHRqFzsQGRW3ta7XczTO1yq6/ks6xHje/yQqeGdbJLD07StCwslA7JDukA5u0WuPkaozjRKLrN9ShiHDq8wKBgQCOxpoo1NekSuQsjkKuTBhVMHPiw5Y2kk60GgFbzkArETwIvP1Oe4F4m9n+9f1L4EtbUGtYyQ6zgiqWsuA33KHPLw3cPsncupBPzZcEsrEcpVuoLrhZA6tAU61HDPdOYm71yq+bfY/b3EaX8yAJ3cCrIWhCsHez2V+R5rUUFZow3QKBgQC81Sr7OfE8qrt49OTh5awNRbEemEuHUS8PZAwuTj5R50xg8fJmqDfkIi7hjCtU1f1Rvi7pCQL6nm9gD+qnhUWcd8+bJPOxChyouKMsaZXaYCcEszs/fcRWc2AxMtYTFtTRzlGILKhzn8k3FkLKHtDLafDLbK+M06Gg5PEOeK1PqwKBgQCn0r+NDE09ImX9PVymwomScpPRWC/SxVgmzx3mGDG0AaRqGjfa1hskqxwx8eRfT3exwvQ9dvYaPyfATyQZ0uEkg7bJ47Jfr7YHkcKMWM1+u8iVxN4mn6Kj3aSfy71iumPzm8J/9BZHX1cTLzXDi1OiP+mQs+UXqwXGfvx/UZHSgw=="
	public64 := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmP4t8L4HqqNYErrB39QLj2FPp0TDS5GXiD6QOmI7xOYkNug0F3a0NFmprLwCs9xhEDMEaprXrBEIS8oCY+MCXTuDw4RjapTm7FkMeDdXMq27xqLk2ztoVIZ3rD9z1ETEQRKSXdQP5KCpoetOouKCnFyrlQcI6ATHTzwjGeEHIXCd9ChS7Yl0/TgaRopu0oBQYip1sTv+AhndSpIsUY+CA9JcmIUZ5wlLxabTfNtleIOcq+70lLB9zOyMVTuEoBcxWYzPab7qNBtr3+dTSfnyt8vaytuP1beu9bBH+fbNvDl7skVxK44HyKA3TNra4idhavI2UFfZFb0eoJ+65mX1JwIDAQAB"

	privateKey, err := LoadRsaPrivateKey(private64)
	if err != nil {
		t.Error(err)
	}

	publicKey, err := LoadRsaPublicKey(public64)
	if err != nil {
		t.Error(err)
	}

	message := []byte("1020000189025321" + "-" + "1681178471210")
	encoded, err := RsaSign(privateKey, message)
	if err != nil {
		t.Error(err)
	}
	// 使用公钥对签名进行验证
	err = RsaVerify(encoded, publicKey, message)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Signature:" + encoded)
	fmt.Println("Signature verified")
}

func TestIntMinBasic(t *testing.T) {
	t.Errorf("IntMin(2, -2) = %d; want -2", 123)
}

func TestGetLen(t *testing.T) {
	//获取rsa 公钥长度
	PubKeyLen, _ := GetPubKeyLen(publicKey)
	fmt.Println("public key len is ", PubKeyLen)

	//获取rsa 私钥长度
	PriKeyLen, _ := GetPriKeyLen(privateKey)
	fmt.Println("private key len is ", PriKeyLen)
}

func TestGenKeysByLen(t *testing.T) {
	err := GenRsaKey(1024, "/tmp/xxx")
	if err != nil {
		return
	}
}
