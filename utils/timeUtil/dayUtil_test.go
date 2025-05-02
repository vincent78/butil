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
