package strUtil

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"io"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// CompareVersion 版本比较，
// src > dst  返回 1
// src = dst  返回 0
// src < dst  返回 -1
func CompareVersion(src, dst string) int {
	if src == dst {
		return 0
	}

	sa := strings.Split(src, ".")
	da := strings.Split(dst, ".")
	sar := FillZero(sa, 3)
	dar := FillZero(da, 3)

	for i, v := range sar {
		sn, _ := strconv.Atoi(v)
		dn, _ := strconv.Atoi(dar[i])
		if sn > dn {
			return 1
		} else if sn < dn {
			return -1
		}
	}
	return 0
}

func FillZero(sl []string, i int) []string {
	if len(sl) == i {
		return sl
	} else if len(sl) > i {
		return sl[:i]
	} else {
		t := make([]string, i, i)
		for id, v := range sl {
			t[id] = v
		}
		return t
	}
}

func ToAscii(str string) string {
	return strconv.QuoteToASCII(str)
}

/*****************************************
 *
 * bytes
 *
 *****************************************/

// Bytes2String 直接转换底层指针，两者指向的相同的内存，改一个另外一个也会变。
// 效率是string([]byte{})的百倍以上，且转换量越大效率优势越明显。
func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// String2Bytes 直接转换底层指针，两者指向的相同的内存，改一个另外一个也会变。
// 效率是string([]byte{})的百倍以上，且转换量越大效率优势越明显。
// 转换之后若没做其他操作直接改变里面的字符，则程序会崩溃。
// 如 b:=String2bytes("xxx"); b[1]='d'; 程序将panic。
func String2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func s2b(s string) (b []byte) {
	/* #nosec G103 */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	/* #nosec G103 */
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

/*****************************************
 *
 * random
 *
 *****************************************/

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// GetRandomString 获取指定长度的随机字符串
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%&*"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// RandStringBytesMaskImpr 生成指定长度的随机字符串
func RandStringBytesMaskImpr(n int) []byte {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

/*
*
l <= 9
*/
func RandomString(l int) string {
	i := math.Pow10(l) - 1
	fmtStr := fmt.Sprintf("%%0%vv", l)
	r := fmt.Sprintf(fmtStr, rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(int32(i)))
	return r
}

/*****************************************
 *
 * yaml
 *
 *****************************************/

func ParseFromYaml(str string, v interface{}) error {
	if bs, err := yaml.YAMLToJSON(String2Bytes(str)); err == nil {
		Parse(Bytes2String(bs), &v)
		return nil
	} else {
		return err
	}
}

// ToStr 获取变量的字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func ToStr(value interface{}) string {
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = Bytes2String(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}

func Md5(str string) string {
	if str == "" {
		return ""
	}
	w := md5.New()
	io.WriteString(w, str)
	return fmt.Sprintf("%v", hex.EncodeToString(w.Sum(nil)))

}

// UnsafeEqual 比较字符串与字节数组中的内容是否一致
func UnsafeEqual(a string, b []byte) bool {
	bbp := *(*string)(unsafe.Pointer(&b))
	return a == bbp
}
