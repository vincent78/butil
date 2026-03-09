package objUtil

import "github.com/bytedance/sonic"

func ObjToStr(o any) string {
	bytes, _ := sonic.Marshal(o)
	return string(bytes)
}

func ObjFromStr(str string, o any) error {
	return sonic.Unmarshal([]byte(str), o)
}
