package byteUtil

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"hash/crc32"
)

func Unbox(s []byte) (string, error) {
	b, err := base64.RawURLEncoding.DecodeString(string(s))
	if err != nil {
		return "", err
	}
	if len(b) < 4 {
		return "", errors.New("invalid name")
	}
	v, err := base64.RawURLEncoding.DecodeString(string(b[4:]))
	if err != nil {
		return "", err
	}
	if crc32.ChecksumIEEE(v) != binary.BigEndian.Uint32(b[:4]) {
		return "", errors.New("invalid name")
	}
	return string(v), nil
}

func Box(s string) ([]byte, error) {
	es := []byte(base64.RawURLEncoding.EncodeToString([]byte(s)))
	crc := crc32.ChecksumIEEE(es)
	crcBytes := make([]byte, 4)
	binary.BigEndian.AppendUint32(crcBytes, crc)
	bs := make([]byte, 0, 4+len(es))
	bs = append(bs, crcBytes...)
	bs = append(bs, es...)

	rs := []byte(base64.RawURLEncoding.EncodeToString(bs))
	return rs, nil
}
