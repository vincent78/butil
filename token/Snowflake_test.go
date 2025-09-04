package token

import (
	"testing"
	"time"

	utime "github.com/vincent78/butil/utils/timeUtil"
)

func TestToken(t *testing.T) {
	for i := 0; i < 10; i++ {
		println("token:%v", NextStrToken())
	}
}

func TestStartTime(t *testing.T) {
	println(utime.TimeStr(utime.TimeFromStampMillis(1577836800000)))
}

func TestGetTimestamp(t *testing.T) {
	tm, _ := time.Parse("01/02/2006 03:04:05", "01/01/2020 00:00:00")
	println(utime.TimeMillis(tm))
}
