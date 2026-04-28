package objUtil

import "github.com/tidwall/sjson"

func Set2JsonStr(json, key string, value any) (string, error) {
	return sjson.Set(json, key, value)
}
