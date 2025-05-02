package timeUtil

import "testing"

/*
*

98B70E11100314  2020-03-16 17:14:47
*/
func TestCp56time2aToTime(t *testing.T) {
	ti, e := Cp56time2aToTime("98B70E11100314")
	if e != nil {
		t.Error(e)
	} else {
		t.Logf("time : %v", ti)
	}
}

func TestTimeToCP56time2a(t *testing.T) {
	ti, e := Parse("2020-03-16 17:14:47")
	if e != nil {
		t.Error(e)
	} else {
		t.Logf("hex: %v", TimeToCp56time2a(ti))
	}
}
