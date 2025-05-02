package jsonUtil

import (
	"fmt"
	"github.com/tidwall/gjson"
)

// 使用路径的方式获取到对应的值
// 层级比较多时取值方便
// gjson.Parse(json).Get("name").Get("last")
// gjson.Get(json, "name.last").Exists()
// gjson.GetBytes(json, path)
// gjson.GetMany(json, "name.first", "name.last", "age")

/**
{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}

"name.last"          >> "Anderson"
"age"                >> 37
"children"           >> ["Sara","Alex","Jack"]
"children.#"         >> 3
"children.1"         >> "Alex"
"child*.2"           >> "Jack"
"c?ildren.0"         >> "Sara"
"fav\.movie"         >> "Deer Hunter"
"friends.#.first"    >> ["Dale","Roger","Jane"]
"friends.1.last"     >> "Craig"


friends.#(last=="Murphy").first    >> "Dale"
friends.#(last=="Murphy")#.first   >> ["Dale","Jane"]
friends.#(age>45)#.last            >> ["Craig","Murphy"]
friends.#(first%"D*").last         >> "Murphy"
friends.#(first!%"D*").last        >> "Craig"
friends.#(nets.#(=="fb"))#.first   >> ["Dale","Roger"]
*/

func GetStringValue(src, pos string) string {
	if gjson.Valid(src) {
		return gjson.Get(src, pos).String()
	} else {
		return ""
	}
}

func GetIntValue(src, pos string) int64 {
	if gjson.Valid(src) {
		return gjson.Get(src, pos).Int()
	} else {
		return -1
	}
}

func GetValue(src, pos string) (gjson.Result, error) {
	if gjson.Valid(src) {
		return gjson.Get(src, pos), nil
	} else {
		return gjson.Parse("{}"), fmt.Errorf("is not json format str: %v", src)
	}
}

func ExistValue(src, path string) bool {
	if gjson.Valid(src) {
		return gjson.Get(src, path).Exists()
	} else {
		return false
	}
}
