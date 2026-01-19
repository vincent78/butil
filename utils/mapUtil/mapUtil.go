package mapUtil

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/goinggo/mapstructure"
)

/*
*
go get github.com/goinggo/mapstructure

	func MapToStructDemo(){
	        mapInstance := make(map[string]interface{})
	        mapInstance["Name"] = "jqw"
	        mapInstance["Age"] = 18

	        var people People
	        err := mapstructure.Decode(mapInstance, &people)
	        if err != nil {
	                fmt.Println(err)
	        }
	        fmt.Println(people)
	}
*/
func MapToStruct(m map[string]any, t any) {
	err := mapstructure.Decode(m, t)
	if err != nil {
		panic(err)
	}
}

func StructToMap(obj any) map[string]any {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]any)
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}

func ToJsonStr(obj any) string {
	if obj == nil {
		return ""
	}

	bytes, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func FromJsonStr(str string) any {
	var r any
	e := json.Unmarshal([]byte(str), &r)
	if e != nil {
		println(fmt.Sprintf("FromJsonStr error: %v", e.Error()))
		return nil
	}
	return r
}

func ObjByAnchorFromMap(obj any, anchor []string) any {
	if obj == nil {
		return nil
	}
	if reflect.TypeOf(obj).Name() == "string" {
		obj = FromJsonStr(obj.(string))
	}
	mapObj, err := obj.(map[string]any)
	if !err {
		return nil
	}
	if anchor == nil || len(anchor) == 0 {
		return obj
	}
	k := anchor[0]
	if k == "" {
		return obj
	} else {
		t := mapObj[anchor[0]]
		if len(anchor) > 1 {
			return ObjByAnchorFromMap(t, anchor[1:])
		} else {
			return t
		}
	}
}

func Add(source, target map[string]any) map[string]any {
	if target == nil || len(target) == 0 {
		return source
	}

	if source == nil || len(source) == 0 {
		return target
	}

	for key, value := range source {
		target[key] = value
	}
	return target
}

func mergeMaps(maps ...map[string]any) map[string]any {
	result := make(map[string]any)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
