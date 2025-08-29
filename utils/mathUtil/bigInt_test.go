package mathUtil

import (
	"fmt"
	"testing"
)

func TestSuoJinSuanFa(t *testing.T) {

	fmt.Println(SuoJinSuanFa("2071.82")) // 输出: 0.05001
	fmt.Println(SuoJinSuanFa("5000"))    // 输出: 0.05001
	fmt.Println(SuoJinSuanFa("80"))      // 输出: 0.08
	fmt.Println(SuoJinSuanFa("1000"))    // 输出: 0.01001
	fmt.Println(SuoJinSuanFa("1010"))    // 输出: 0.01001
	fmt.Println(SuoJinSuanFa("1002"))    // 输出: 0.01001
	fmt.Println(SuoJinSuanFa("8"))       // 输出: 0.01001

}
func TestSuoJinFa2(t *testing.T) {

	fmt.Println(SuoJinSuanFa2("19785000")) // 输出: 1.1978501
	fmt.Println(SuoJinSuanFa2("5000"))     // 输出: 0.501
	fmt.Println(SuoJinSuanFa2("80"))       // 输出: 0.80
	fmt.Println(SuoJinSuanFa2("69863"))    // 输出: 0.69863

}
func TestSuoJinSuanFaReverse(t *testing.T) {
	testCases := []string{
		"0.05001",  // 5000
		"0.08",     // 80
		"0.069863", // 69863
		"0.01001",  // 1000
		"0.123",    // 123（无补1情况）
	}

	for _, s := range testCases {
		original, err := SuoJinSuanFaReverse(s)
		if err != nil {
			fmt.Printf("错误：%v\n", err)
			continue
		}
		fmt.Printf("格式化：%-10s => 原始值：%d\n", s, original)
	}
}
func TestSuoJinSuanFa2Reverse(t *testing.T) {
	testCases := []struct {
		input  string
		output string
	}{
		{"19785000", "1.9785001"},
		{"5000", "0.5001"},
		{"80", "0.8"},
		{"69863", "0.69863"},
		{"100", "1.11"},
		{"1234500", "1.23451"},
		{"200", "0.21"},
	}

	fmt.Println("=== 格式化测试 ===")
	for _, tc := range testCases {
		result := SuoJinSuanFa2(tc.input)
		fmt.Printf("%s → %s (期望: %s)\n", tc.input, result, tc.output)
		//if result != tc.output {
		//	fmt.Println("  错误！不匹配")
		//}
	}

	fmt.Println("\n=== 反向解析测试 ===")
	for _, tc := range testCases {
		reversed, err := SuoJinSuanFa2Reverse(tc.output)
		if err != nil {
			fmt.Printf("解析 %s 错误: %v\n", tc.output, err)
			continue
		}
		fmt.Printf("%s → %s (原始: %s)\n", tc.output, reversed, tc.input)
		//if reversed != tc.input {
		//	fmt.Println("  错误！不匹配")
		//}
	}

	// 额外测试一些边缘情况
	extraTests := []struct {
		formatted string
		expected  string
	}{
		{"1.2345671", "123456700"},
		{"0.1231", "12300"},
		{"1.11", "100"},
		{"0.21", "200"},
	}

	fmt.Println("\n=== 额外反向解析测试 ===")
	for _, et := range extraTests {
		reversed, err := SuoJinSuanFa2Reverse(et.formatted)
		if err != nil {
			fmt.Printf("解析 %s 错误: %v\n", et.formatted, err)
			continue
		}
		fmt.Printf("%s → %s (期望: %s)\n", et.formatted, reversed, et.expected)
		//if reversed != et.expected {
		//	fmt.Println("  错误！不匹配")
		//}
	}
}

func TestConvertFloat2BigInt(t *testing.T) {
	s := "7357499.492611614930156458"
	r := ConvertFloat2BigInt(s, 19)
	t.Log(r)
}
