package secureUtil

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"log"
	"os"
	"strings"
)

// 加密
func RsaEncrypt(origData []byte, publicKey []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecrypt(ciphertext []byte, privateKey []byte) ([]byte, error) {
	//解密
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

// ecdsaCmd represents the doc command
// keyName 需要输入全路径
func KeyPairs(keyName string) error {
	//elliptic.P256(),elliptic.P384(),elliptic.P521()

	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		log.Fatal(err)
		return err
	}
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	privateBs := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	privateFile, err := os.Create(keyName + ".private.pem")
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = privateFile.Write(privateBs)
	if err != nil {
		log.Fatal(err)
		return err
	}
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(privateKey.Public())
	publicBs := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
	publicKeyFile, err := os.Create(keyName + ".public.pem")
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = publicKeyFile.Write(publicBs)
	if err != nil {
		log.Fatal(err)
		return err
	} else {
		return nil
	}
}

func chunkSplit(s string, chunkSize int, endLine string) string {
	var chunks []string

	for len(s) > 0 {
		if len(s) < chunkSize {
			chunkSize = len(s)
		}
		chunks = append(chunks, s[:chunkSize])
		s = s[chunkSize:]
	}

	return strings.Join(chunks, endLine)
}

func LoadRsaPrivateKey(private64 string) (*rsa.PrivateKey, error) {
	pemSplit := chunkSplit(private64, 64, "\n")
	privatePem := "-----BEGIN PRIVATE KEY-----\n" + pemSplit + "\n-----END PRIVATE KEY-----\n"
	privateKeyBlock, _ := pem.Decode([]byte(privatePem))
	privateKeyParsed, err := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKeyParsed.(*rsa.PrivateKey), nil
}

func LoadRsaPublicKey(public64 string) (*rsa.PublicKey, error) {
	pemSplit := chunkSplit(public64, 64, "\n")
	publicPem := "-----BEGIN PUBLIC KEY-----\n" + pemSplit + "\n-----END PUBLIC KEY-----\n"
	publicKeyBlock, _ := pem.Decode([]byte(publicPem))
	publicKeyParsed, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return publicKeyParsed.(*rsa.PublicKey), nil
}

func RsaSign(privateKey *rsa.PrivateKey, data []byte) (string, error) {
	// 使用私钥对消息进行签名
	hash := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString([]byte(signature)), nil
}

func RsaVerify(signature string, publicKey *rsa.PublicKey, data []byte) error {
	hash := sha256.Sum256(data)
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}
	// 使用公钥对签名进行验证
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], decoded)
	if err != nil {
		return err
	}
	return nil
}

/**
 * @brief  获取RSA公钥长度
 * @param[in]       PubKey				    RSA公钥
 * @return   成功返回 RSA公钥长度，失败返回error	错误信息
 */
func GetPubKeyLen(PubKey []byte) (int, error) {
	if PubKey == nil {
		return 0, errors.New("input arguments error")
	}

	block, _ := pem.Decode(PubKey)
	if block == nil {
		return 0, errors.New("public rsaKey error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return 0, err
	}
	pub := pubInterface.(*rsa.PublicKey)

	return pub.N.BitLen(), nil
}

/**
 * @brief  获取RSA私钥长度
 * @param[in]       PriKey				    RSA私钥
 * @return   成功返回 RSA私钥长度，失败返回error	错误信息
 */
func GetPriKeyLen(PriKey []byte) (int, error) {
	if PriKey == nil {
		return 0, errors.New("input arguments error")
	}

	block, _ := pem.Decode(PriKey)
	if block == nil {
		return 0, errors.New("private rsaKey error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return 0, err
	}

	return priv.N.BitLen(), nil
}

func GenRsaKey(bits int, path string) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "私钥",
		Bytes: derStream,
	}
	file, err := os.Create(path + ".private.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "公钥",
		Bytes: derPkix,
	}
	file, err = os.Create(path + ".public.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}
