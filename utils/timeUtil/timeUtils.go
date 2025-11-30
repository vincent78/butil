package timeUtil

import (
	"time"
)

type Time time.Time

const (
	DateFormat        = "2006-01-02"
	DateCompactFormat = "20060102"
	TimeFormat        = "2006-01-02 15:04:05"
	LongFormat        = "2006-01-02 15:04:05.999"
	CompactFormat     = "20060102150405" // 紧凑型格式
	AllFormat         = "2006-01-02 15:04:05.999999999 -0700 MST"
)

/**********************************************

	转换

 **********************************************/

// NowStr 得到当前时间的通用格式字符串
func NowStr() string {
	return TimeStr(time.Now())
}

// NowStr2 得到当前时间的字符串（秒）
func NowStr2() string {
	return TimeFmtStr(time.Now(), CompactFormat)
}

func NowFmtStr(fmtStr string) string {
	return TimeFmtStr(time.Now(), fmtStr)
}

func TimeStr(t time.Time) string {
	return TimeFmtStr(t, TimeFormat)
}

func TimeFmtStr(t time.Time, fmt string) string {
	return t.Format(fmt)
}

func NowNano() int64 {
	return TimeSecond(time.Now())
}
func TimeNano(t time.Time) int64 {
	return t.UnixNano()
}

// NowMillis 得到当前时间戳 毫秒
func NowMillis() int64 {
	return TimeMillis(time.Now())
}
func NowUtcMillis() int64 {
	return UTCTimeMillis(time.Now().UTC())
}
func TimeMillis(t time.Time) int64 {
	return t.UnixNano() / 1e6
}
func UTCTimeMillis(t time.Time) int64 {
	// shanghai, _ := time.LoadLocation("Asia/Shanghai")
	if t.Location() != time.UTC {
		t.In(time.UTC)
	}
	return TimeMillis(t)
}

// NowSecond 得到当前时间戳 秒
func NowSecond() int64 {
	return TimeSecond(time.Now())
}
func TimeSecond(t time.Time) int64 {
	return t.Unix()
}

func UTCTimeSecond(t time.Time) int64 {
	// shanghai, _ := time.LoadLocation("Asia/Shanghai")
	if t.Location() != time.UTC {
		t.In(time.UTC)
	}
	return TimeMillis(t)
}

// TimeFromStamp 时间戳 参数  s 为秒   ns 为纳秒
// 注：最后的时间为s+ns的结果
// millis := nanos / 1000000 // 毫秒
func TimeFromStamp(s, ns int64) time.Time {
	return time.Unix(s, ns)
}

func TimeFromStampMillis(ms int64) time.Time {
	return time.Unix(0, ms*1e6)
}

func Parse(str string) (time.Time, error) {
	return ParseWithFmt(TimeFormat, str)
}

func ParseWithFmt(fmt, str string) (time.Time, error) {
	return time.ParseInLocation(fmt, str, time.Local)
}

/**********************************************

	json

 **********************************************/

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(TimeFormat)
}
