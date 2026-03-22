package mathUtil

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
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

func ConvertFloatToBigInt(amount float64, decimals int32) *big.Int {
	return decimal.NewFromFloatWithExponent(amount, decimals).BigInt()
}

func ConvertStringToBigInt(amount string, decimals int) *big.Int {
	if decimals == 0 {
		r, _ := new(big.Int).SetString(amount, 10)
		return r
	}
	multiplier := math.Pow10(decimals)
	bf, _ := big.NewFloat(0).SetString(amount)
	f := big.NewFloat(0).Mul(big.NewFloat(multiplier), bf)
	i := new(big.Int)
	f.Int(i)
	return i
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

func ConvertFloat2BigInt(num string, decimals int) *big.Int {
	s := num
	l := len(s)
	for i := 0; i < decimals; i++ {
		p := strings.Index(s, ".")
		if p < 0 {
			s = s + "0"
		} else if p >= l-2 {
			s = s[:p] + s[p+1:p+2]
		} else {
			s = s[:p] + s[p+1:p+2] + "." + s[p+2:]
		}
	}
	dot := strings.Index(s, ".")
	if dot >= 0 {
		s = s[:dot]
	}
	r, _ := new(big.Int).SetString(s, 10)
	return r
}

/***********************************************************
 *
 * 高精度操作
 *
 ***********************************************************/

func BigDiv(num1, num2 string) (string, error) {
	var d1, d2 decimal.Decimal
	var err error
	if d1, err = decimal.NewFromString(num1); err != nil {
		return "", err
	}
	if d2, err = decimal.NewFromString(num2); err != nil {
		return "", err
	}
	return d1.Div(d2).String(), nil
}

func BigDivExp(num string, exp int32) (string, error) {
	//d2 := big.NewInt(int64(math.Pow10(int(exp))))
	d2 := decimal.New(1, exp)
	return BigDiv(num, d2.String())
}

func BigMul(num1, num2 string) (string, error) {
	var d1, d2 decimal.Decimal
	var err error
	if d1, err = decimal.NewFromString(num1); err != nil {
		return "", err
	}
	if d2, err = decimal.NewFromString(num2); err != nil {
		return "", err
	}
	return d1.Mul(d2).String(), nil
}

func BigMulExp(num string, exp int32) (string, error) {
	d2 := decimal.New(1, exp)
	return BigMul(num, d2.String())
}

func BigFormat(num string, n int32) (string, error) {
	if d1, err := decimal.NewFromString(num); err != nil {
		return "", err
	} else {
		rounded := d1.Round(n)
		return rounded.String(), nil
	}
}

func BigInt2HexStr(num int64, bitSizes ...int) string {
	n := new(big.Int).SetInt64(num)
	l := 256
	if len(bitSizes) > 0 {
		l = bitSizes[0]
	}

	// %x: 小写十六进制
	// %X: 大写十六进制
	// %#x: 带 0x 前缀的小写
	// %064x: 补齐到 64 位长度（常用于以太坊 Hash 或 Address）

	//fmt.Printf("%x\n", n)    // ff
	//fmt.Printf("%X\n", n)    // FF
	//fmt.Printf("%#x\n", n)   // 0xff
	//fmt.Printf("%064x\n", n) // 00000000000000000000000000000000000000000000000000000000000000ff
	return fmt.Sprintf("%0"+strconv.Itoa(l)+"x", n)
}

// bitsizes是指hexStr最大的长度,不传默认为256位（这是以太坊中常用的)
func HexStr2BigInt(hexStr string, bitSizes ...int) *big.Int {
	bitSize := 256
	if len(bitSizes) > 0 {
		bitSize = bitSizes[0]
	}

	str := strings.Replace(hexStr, "0x", "", 1)
	str = strings.Replace(str, "0X", "", 1)

	str = strings.ToLower(str)

	// 1. 将十六进制字符串解析为 big.Int
	n := new(big.Int)
	n.SetString(str, 16)

	if CheckOverflow(n, bitSizes...) == 0 {
		return n
	}

	// 2. 定义 2^256 (用于补码计算)
	// 256 位十六进制的最大值边界
	maxVal := new(big.Int).Lsh(big.NewInt(1), uint(bitSize))
	// 3. 计算补码对应的负数值
	// 如果 n > 2^255，说明符号位为 1，是一个负数
	// 这里我们直接计算 n - 2^256 即可得到正确的负值
	return new(big.Int).Sub(n, maxVal)
}

// CheckOverflow 检查 val 在指定的 bitSize 下是否溢出（有符号整数）
func CheckOverflow(val *big.Int, bitSize ...int) int {
	// 1. 设置默认位数为 256
	n := 256
	if len(bitSize) > 0 {
		n = bitSize[0]
	}

	// 2. 计算最大正数边界 (2^(n-1) - 1)
	// limit = 1 << (n-1)
	limit := new(big.Int).Lsh(big.NewInt(1), uint(n-1))
	maxPos := new(big.Int).Sub(limit, big.NewInt(1))

	// 3. 计算最小负数边界 (-2^(n-1))
	minNeg := new(big.Int).Neg(limit)

	// 4. 执行范围检查： minNeg <= val <= maxPos
	if val.Cmp(maxPos) > 0 {
		//return fmt.Errorf("上溢出: 数值大于 %d 位下的最大正数", n)
		return 1
	}
	if val.Cmp(minNeg) < 0 {
		//return fmt.Errorf("下溢出: 数值小于 %d 位下的最小负数", n)
		return -1
	}

	return 0
}
