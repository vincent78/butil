package strUtil

import (
	"testing"
)

func TestFmtSlice(t *testing.T) {
	fmtStr := "p:%v"
	str := FmtSlice(fmtStr, " ", "1", "2", "3")
	t.Log(str)
	str = FmtSlice(fmtStr, " ", "1", 1, 2.0, true, []string{"ddd"})
	t.Log(str)
}
