package timeUtil

import (
	"testing"
	"time"
)

func TestCompareDay(t *testing.T) {
	t1 := time.Now()
	t2 := CalDay(t1, 1)

	r := CompareDay(t1, t2)
	t.Logf("t1: %v\t\t t2: %v\t result: %v", t1, t2, r)

	r1 := CompareDay(t2, t1)
	t.Logf("t2: %v\t\t t1: %v \t result: %v", t2, t1, r1)
}

func TestGetDiffDays(t *testing.T) {
	t1, _ := time.Parse(TimeFormat, "2025-06-25 23:12:11")
	t2, _ := time.Parse(TimeFormat, "2025-06-24 23:59:59")

	n := GetDiffDays(t1, t2)
	t.Logf("t1: %v\t\t t2: %v\t result: %v", t1, t2, n)

	t3, _ := time.Parse(TimeFormat, "2025-06-25 00:12:11")
	t4, _ := time.Parse(TimeFormat, "2025-06-24 23:59:59")

	n = GetDiffDays(t3, t4)
	t.Logf("t3: %v\t\t t4: %v\t result: %v", t3, t4, n)
}
