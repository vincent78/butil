//math 的增强功能
package mathUtil

//int 类型取绝对值
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

//int 类型取绝对值
func AbsInt16(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}

// 判断当前数字是否在范围中,前包后不包
// eg  between(1,1,3) true
// eg   between(3,1,3) false
func Between(x int, start int, end int) bool {
	if x >= start && x < end {
		return true
	}
	return false
}

// 判断当前数字是否在范围中,前包后不包
// eg  between(1,1,3) true
// eg   between(3,1,3) false
func BetweenInt16(x int16, start int16, end int16) bool {
	if x >= start && x < end {
		return true
	}
	return false
}
