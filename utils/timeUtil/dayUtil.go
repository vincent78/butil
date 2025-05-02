package timeUtil

import "time"

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
