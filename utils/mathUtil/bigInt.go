package mathUtil

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

func DivideWithPrecision(value *big.Int, precision int) string {
	divisor := big.NewInt(1_000_000)

	// 计算整数部分
	intPart := new(big.Int).Div(value, divisor)

	// 计算小数部分 (value % divisor)
	remainder := new(big.Int).Mod(value, divisor)

	// 格式化为字符串: "整数.小数"
	return fmt.Sprintf("%s.%0*d", intPart.String(), precision, remainder)
}

func ConvertBigIntToFloat(amount *big.Int, decimals int) float64 {
	f := new(big.Float).SetInt(amount)
	divisor := new(big.Float).SetFloat64(math.Pow10(decimals))
	result, _ := new(big.Float).Quo(f, divisor).Float64()
	return result
}

func ConvertFloatToBigInt(amount float64, decimals int) *big.Int {
	multiplier := math.Pow10(decimals)
	bf := new(big.Float).SetFloat64(amount * multiplier)
	result := new(big.Int)
	bf.Int(result)
	return result
}

// SuoJinSuanFa 将整数转换为特定缩进格式的字符串（如 5000 → "0.05001"）
func SuoJinSuanFa(numStr string) string {
	//numStr := strconv.Itoa(num)

	// 检查是否以 "00" 结尾
	if strings.HasSuffix(numStr, "00") {
		return "0.0" + numStr[:len(numStr)-2] + "1" // 去掉00并补1
	}
	return "0.0" + numStr // 直接拼接
}

// SuoJinSuanFaReverse 反向还原 SuoJinSuanFa 生成的字符串为原始整数
func SuoJinSuanFaReverse(formatted string) (int, error) {
	// 1. 移除前缀 "0.0"
	if !strings.HasPrefix(formatted, "0.0") {
		return 0, fmt.Errorf("无效格式：必须以 '0.0' 开头")
	}
	numStr := formatted[len("0.0"):]

	// 2. 检查是否以 "1" 结尾
	if strings.HasSuffix(numStr, "1") {
		// 去掉末尾的 "1" 并补 "00"
		numStr = numStr[:len(numStr)-1] + "00"
	}

	// 3. 转换为整数（自动去除前导零）
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, fmt.Errorf("转换失败：%v", err)
	}

	return num, nil
}

// 1开头用1. 其他开头用0.，尾数00用1占位。例如19785000=1.9785001，5000=0.5001，80=0.8，69863=0.69863
func SuoJinSuanFa2(numStr string) string {
	// 处理前缀
	var prefix string
	if strings.HasPrefix(numStr, "1") {
		prefix = "1."
	} else {
		prefix = "0."
	}

	// 处理后缀
	var body string
	if strings.HasSuffix(numStr, "00") {
		body = numStr[:len(numStr)-2] + "1"
	} else {
		body = numStr
	}

	// 组合结果
	return prefix + body
}

// SuoJinSuanFa2Reverse 反向还原 SuoJinSuanFa2 生成的字符串为原始整数
func SuoJinSuanFa2Reverse(formatted string) (string, error) {
	// 检查格式是否正确
	if !strings.Contains(formatted, ".") {
		return "0", fmt.Errorf("invalid format: missing dot")
	}

	parts := strings.Split(formatted, ".")
	if len(parts) != 2 {
		return "0", fmt.Errorf("invalid format: expected exactly one dot")
	}

	prefix := parts[0]
	body := parts[1]

	// 处理前缀
	var original string
	if prefix == "1" {
		original = "1" + body
	} else if prefix == "0" {
		original = body
	} else {
		return "0", fmt.Errorf("invalid prefix: %s", prefix)
	}

	// 处理后缀
	if strings.HasSuffix(original, "1") {
		// 检查是否真的是占位符1
		if len(original) > 1 && !strings.HasSuffix(original, "11") {
			original = original[:len(original)-1] + "00"
		}
	}
	//
	//// 转换为整数
	//num, err := strconv.Atoi(original)
	//if err != nil {
	//	return 0, fmt.Errorf("invalid number format: %v", err)
	//}

	return original, nil
}

// SuoJinSuanFa3 所有订单都用这个金额
func SuoJinSuanFa3(numStr string) string {
	return "0.001" // 直接拼接
}
