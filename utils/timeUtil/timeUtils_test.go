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
	t.Logf(NowFmtStr(DateFormat))
}
