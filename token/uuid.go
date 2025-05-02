package token

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	rand2 "math/rand"
)

type UUID [16]byte

const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GetUUID() (uuidHex string) {
	uuid := NewUUID()
	uuidHex = hex.EncodeToString(uuid[:])
	return
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand2.Intn(len(letterBytes))]
	}
	return string(b)
}

func GetMD5Base64(bytes []byte) (base64Value string) {
	md5Ctx := md5.New()
	md5Ctx.Write(bytes)
	md5Value := md5Ctx.Sum(nil)
	base64Value = base64.StdEncoding.EncodeToString(md5Value)
	return
}

func NewUUID() UUID {
	ns := UUID{}
	safeRandom(ns[:])
	u := newFromHash(md5.New(), ns, RandStringBytes(16))
	u[6] = (u[6] & 0x0f) | (byte(2) << 4)
	u[8] = (u[8]&(0xff>>2) | (0x02 << 6))

	return u
}

func newFromHash(h hash.Hash, ns UUID, name string) UUID {
	u := UUID{}
	h.Write(ns[:])
	h.Write([]byte(name))
	copy(u[:], h.Sum(nil))

	return u
}

func safeRandom(dest []byte) {
	if _, err := rand.Read(dest); err != nil {
		panic(err)
	}
}

func (u UUID) String() string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}

/*******************************************
 *
 * 外部引入
 *
 ********************************************/

// New Simple call
func New() string {
	uuid, _ := GenerateUUID()
	return uuid
}

func NewShort() string {
	return RandStringBytes(9)
}

// GenerateRandomBytes is used to generate random bytes of given size.
func GenerateRandomBytes(size int) ([]byte, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return nil, fmt.Errorf("failed to read random bytes: %v", err)
	}
	return buf, nil
}

const uuidLen = 16

// GenerateUUID is used to generate a random UUID
func GenerateUUID() (string, error) {
	buf, err := GenerateRandomBytes(uuidLen)
	if err != nil {
		return "", err
	}
	return FormatUUID(buf)
}

func FormatUUID(buf []byte) (string, error) {
	if buflen := len(buf); buflen != uuidLen {
		return "", fmt.Errorf("wrong length byte slice (%d)", buflen)
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		buf[0:4],
		buf[4:6],
		buf[6:8],
		buf[8:10],
		buf[10:16]), nil
}

func ParseUUID(uuid string) ([]byte, error) {
	if len(uuid) != 2*uuidLen+4 {
		return nil, fmt.Errorf("uuid string is wrong length")
	}

	if uuid[8] != '-' ||
		uuid[13] != '-' ||
		uuid[18] != '-' ||
		uuid[23] != '-' {
		return nil, fmt.Errorf("uuid is improperly formatted")
	}

	hexStr := uuid[0:8] + uuid[9:13] + uuid[14:18] + uuid[19:23] + uuid[24:36]

	ret, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	if len(ret) != uuidLen {
		return nil, fmt.Errorf("decoded hex is the wrong length")
	}

	return ret, nil
}
