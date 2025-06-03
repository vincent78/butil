package timeUtil

import (
	"testing"
	"time"
)

func TestASyncEmitter(t *testing.T) {
	tn := time.Now()
	tsnano := TimeNano(tn)
	println("当前时间（秒）:%v", tn.Unix())
	println("当前时间（毫秒）:%v", tsnano/1e6)
	println("当前时间（纳秒）:%v", tsnano)
	println("当前时间(字符串):%v", TimeStr(tn))
	println("转换成time:%v", TimeStr(TimeFromStamp(tsnano/1e9, 0)))
	println("转换成time:%v", TimeStr(TimeFromStampMillis(TimeMillis(tn))))

}

func TestTimeFormat(t *testing.T) {
	t.Logf("now:%v", NowStr())
	t.Logf("now:%v", NowFmtStr(CompactFormat))
}

func TestNowFmtStr(t *testing.T) {
	t.Log(NowFmtStr(DateFormat))
}

func TestTimeStamp(t *testing.T) {
	tt := 1747611165000
	utc := 8 * 60 * 60 * 1000
	println("转换成time:%v", TimeStr(TimeFromStamp(int64((tt-utc)/1000), 0)))
}
