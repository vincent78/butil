package strUtil

import "strings"

func Empty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func Constains(str string, sub string) bool {
	return strings.Contains(str, sub)
}

// 去掉字符串头尾的不可见字符（空格，换行，制表）
// \p{C}来匹配所有的控制字符（包括不可见字符）
func TrimSpace(str string) string {
	return strings.TrimSpace(str)
}
