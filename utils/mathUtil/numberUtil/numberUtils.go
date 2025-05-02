package numberUtil

import (
	sj "encoding/json"
	"gitee.com/vincent78/gcutil/logger"
	j "github.com/json-iterator/go"
	"strconv"
)

func Int64ToInt(n1 int64) int {
	t := strconv.FormatInt(n1, 10)
	if r, e := strconv.Atoi(t); e != nil {
		logger.Error("Int64ToInt error: %v", e.Error())
		return 0
	} else {
		return r
	}
}

func ToInt(v interface{}) int {
	switch v.(type) {
	case sj.Number:
		u, err := v.(sj.Number).Int64()
		if err != nil {
			u = 0
		}
		return int(u)
	case j.Number:
		u, err := v.(sj.Number).Int64()
		if err != nil {
			u = 0
		}
		return int(u)
	case string:
		u, err := strconv.ParseInt(v.(string), 10, 64)
		if err != nil {
			u = 0
		}
		return int(u)
	case float64:
		return int(v.(float64))
	}
	return 0
}

func ToUInt(v interface{}) uint {
	switch v.(type) {
	case sj.Number:
		u, err := v.(sj.Number).Int64()
		if err != nil {
			u = 0
		}
		return uint(u)
	case j.Number:
		u, err := v.(sj.Number).Int64()
		if err != nil {
			u = 0
		}
		return uint(u)
	case string:
		u, err := strconv.ParseInt(v.(string), 10, 64)
		if err != nil {
			u = 0
		}
		return uint(u)
	case int:
		return uint(v.(int))
	case int64:
		return uint(v.(int64))
	case float64:
		return uint(v.(float64))
	}
	return 0
}

func ToInt64(v interface{}) int64 {
	switch v.(type) {
	case sj.Number:
		u, err := v.(sj.Number).Int64()
		if err != nil {
			u = 0
		}
		return u
	case j.Number:
		u, err := v.(sj.Number).Int64()
		if err != nil {
			u = 0
		}
		return u
	case string:
		u, err := strconv.ParseInt(v.(string), 10, 64)
		if err != nil {
			u = 0
		}
		return u
	case float64:
		return int64(v.(float64))
	}
	return 0
}

func ToUInt64(v interface{}) uint64 {
	switch v.(type) {
	case sj.Number:
		u, err := v.(sj.Number).Int64()
		if err != nil {
			u = 0
		}
		return uint64(u)
	case j.Number:
		u, err := v.(sj.Number).Int64()
		if err != nil {
			u = 0
		}
		return uint64(u)
	case string:
		u, err := strconv.ParseUint(v.(string), 10, 64)
		if err != nil {
			u = 0
		}
		return u
	case float64:
		return uint64(v.(float64))
	}
	return 0
}

func ToString(v interface{}) string {
	switch v.(type) {
	case sj.Number:
		return v.(sj.Number).String()
	case j.Number:
		return v.(j.Number).String()
	case uint64:
		return strconv.FormatUint(v.(uint64), 10)
	case uint8:
		i := v.(uint8)
		return strconv.Itoa(int(i))
	case int64:
		return strconv.FormatInt(v.(int64), 10)
	case *uint64:
		vu := v.(*uint64)
		return strconv.FormatUint(*vu, 10)
	case int:
		return strconv.Itoa(v.(int))
	case uint32:
		return strconv.FormatUint(uint64(v.(uint32)), 10)
	case string:
		return v.(string)
	}
	return ""
}
