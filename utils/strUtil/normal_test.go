package strUtil

import (
	"regexp"
	"testing"
)

func TestTrimSpace(t *testing.T) {
	str := " \t\nHello, World!\n\t\x00\x01\x02" // 包含控制字符和普通空白字符
	re := regexp.MustCompile(`[\p{C}\s]+`)      // \p{C} 匹配所有控制字符，\s 匹配所有空白字符
	trimmedStr := re.ReplaceAllString(str, "")  // 替换为空字符串
	t.Logf("result1: %s", trimmedStr)

	t.Logf("result2: %s", TrimSpace(str))
}
