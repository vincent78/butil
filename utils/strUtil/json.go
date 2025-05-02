package strUtil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/vincent78/gcutil/utils/mapUtil"
	"gitee.com/vincent78/gcutil/utils/reflectx"
	"strings"
)

// ToJsonStr 将对象转为json字符串
func ToJsonStr(obj interface{}) string {

	if str, ok := obj.(string); ok {
		return str
	}

	b, err := json.Marshal(obj)
	if err != nil {
		fmt.Errorf("strUtil.ToJsonStr error:%v", err.Error())
		return ""
	}
	s := fmt.Sprintf("%+v", Bytes2String(b))
	r := strings.Replace(s, `\u003c`, "<", -1)
	r = strings.Replace(r, `\u003e`, ">", -1)
	return r
}

func ToBytes(obj interface{}) []byte {
	if str, ok := obj.(string); ok {
		return []byte(str)
	}

	b, err := json.Marshal(obj)
	if err != nil {
		fmt.Errorf("strUtil.ToJsonStr error:%v", err.Error())
		return nil
	} else {
		return b
	}
}

func Parse(str string, v interface{}) error {
	return json.Unmarshal([]byte(str), &v)
}

func ParseByBytes(bytes []byte, v interface{}) error {
	return json.Unmarshal(bytes, &v)
}

func ObjByAnchor(str, anchor string) interface{} {
	if str == "" {
		return str
	}
	if len(anchor) == 0 {
		return str
	}
	return mapUtil.ObjByAnchorFromMap(mapUtil.FromJsonStr(str), strings.Split(anchor, ":"))
}

// IsArray 判断当前字符串是数组还是对象
func IsArray(data []byte) (bool, error) {
	d := json.NewDecoder(bytes.NewBuffer(data))
	t, err := d.Token()
	if err != nil {
		return false, err
	}
	return t == json.Delim('['), nil
}

// 长字符串通过decode的方法进行解析(包括数组）
func Decoder2Json(str string, v interface{}) error {
	if !reflectx.IsPoint(v) {
		return errors.New("v must be point")
	}
	return json.NewDecoder(strings.NewReader(str)).Decode(v)
}
