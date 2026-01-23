package timeUtil

import (
	"fmt"
	"strings"
	"time"
)

func Yesterday() time.Time {
	return CalDay(time.Now(), -1)
}

func Tomorrow() time.Time {
	return CalDay(time.Now(), 1)
}

func CalDay(t time.Time, n int) time.Time {
	return t.AddDate(0, 0, n)
}

func SetTimeZero(t time.Time) time.Time {
	ds := TimeFmtStr(t, DateFormat)
	r, _ := ParseWithFmt(DateFormat, ds)
	return r
}

// 按天进行比较      小于：-1 大于：1 等于： 0
func CompareDay(t1, t2 time.Time) int {
	d1, _ := ParseWithFmt(DateFormat, TimeFmtStr(t1, DateFormat))
	d2, _ := ParseWithFmt(DateFormat, TimeFmtStr(t2, DateFormat))
	if d1.Equal(d2) {
		return 0
	} else if d1.After(d2) {
		return 1
	} else {
		return -1
	}
}

// 获取两个时间相差的天数，0表同一天，正数表t1>t2，负数表t1<t2
func GetDiffDays(t1, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)

	return int(t1.Sub(t2).Hours() / 24)
}

// 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

// 获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func main() {
	// Add 时间相加
	now := time.Now()
	// ParseDuration parses a duration string.
	// A duration string is a possibly signed sequence of decimal numbers,
	// each with optional fraction and a unit suffix,
	// such as "300ms", "-1.5h" or "2h45m".
	//  Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	// 10分钟前
	m, _ := time.ParseDuration("-1m")
	m1 := now.Add(m)
	fmt.Println(m1)

	// 8个小时前
	h, _ := time.ParseDuration("-1h")
	h1 := now.Add(8 * h)
	fmt.Println(h1)

	// 一天前
	d, _ := time.ParseDuration("-24h")
	d1 := now.Add(d)
	fmt.Println(d1)

	printSplit(50)

	// 10分钟后
	mm, _ := time.ParseDuration("1m")
	mm1 := now.Add(mm)
	fmt.Println(mm1)

	// 8小时后
	hh, _ := time.ParseDuration("1h")
	hh1 := now.Add(hh)
	fmt.Println(hh1)

	// 一天后
	dd, _ := time.ParseDuration("24h")
	dd1 := now.Add(dd)
	fmt.Println(dd1)

	printSplit(50)

	// Sub 计算两个时间差
	subM := now.Sub(m1)
	fmt.Println(subM.Minutes(), "分钟")

	sumH := now.Sub(h1)
	fmt.Println(sumH.Hours(), "小时")

	sumD := now.Sub(d1)
	fmt.Printf("%v 天\n", sumD.Hours()/24)

}

func printSplit(count int) {
	fmt.Println(strings.Repeat("#", count))
}
