package timeUtil

import "time"

func UTCbyTime(t time.Time) (time.Time, error) {
	// time.LoadLocation("") //等同于"UTC"
	// time.LoadLocation("Local")//服务器设置的时区
	// time.LoadLocation("America/Los_Angeles")

	local, err := time.LoadLocation("Local")
	if err != nil {
		return time.Now(), err
	} else {
		t1 := t.In(local)
		return t1, nil
	}
}

func ChangeLoc(t time.Time, target string) time.Time {
	// "Asia/Shanghai"
	//fmt.Println(t)
	loc, _ := time.LoadLocation(target)
	t1 := t.In(loc)
	//fmt.Println(t1)
	return t1
}
