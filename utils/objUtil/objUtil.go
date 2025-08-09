package objUtil

import (
	"errors"
	"fmt"
	"github.com/vincent78/butil/utils/mapUtil"
	"github.com/vincent78/butil/utils/strUtil"
	"reflect"
	"strconv"
	"strings"
)

func ToString(obj interface{}) string {
	if obj == nil {
		return ""
	}
	str := ""
	switch obj.(type) {
	case string:
		str = obj.(string)
		break
	case int:
		str = strconv.Itoa(obj.(int))
		break
	default:
		str = strUtil.ToJsonStr(obj)
		break
	}
	return str
}

func ToInt(obj interface{}) (int, error) {
	s := ToString(obj)
	return strconv.Atoi(s)
}

func ObjToMap(obj interface{}) map[string]interface{} {
	var r = make(map[string]interface{})
	switch obj.(type) {
	case string:
		err := strUtil.Parse(obj.(string), &r)
		if err != nil {
			return nil
		}
	case []interface{}:
		t := obj.([]interface{})
		for i := 0; i < len(t); i++ {
			r[fmt.Sprintf("key%+v", i)] = t[i]
		}
	case map[interface{}]interface{}:
		t := obj.(map[interface{}]interface{})
		for key, value := range t {
			r[ToString(key)] = value
		}
	case map[string]interface{}:
		r = obj.(map[string]interface{})
	default:
		r["context"] = obj
		break
	}
	return r
}

func ObjsToInterface[T any](objs []T) []interface{} {
	if objs != nil {
		rt := make([]interface{}, 0, len(objs))
		for _, o := range objs {
			rt = append(rt, o)
		}
		return rt
	} else {
		return nil
	}
}

func ObjToStrMap(obj interface{}) map[string]string {
	if obj == nil {
		return nil
	}
	m := ObjToMap(obj)
	r := map[string]string{}
	for key, val := range m {
		r[key] = strUtil.ToJsonStr(val)
	}
	return r
}

func ObjToStrSlice(obj interface{}) []string {
	var r []string
	//fmt.Printf("obj type:%v",reflect.TypeOf(obj))
	switch obj.(type) {
	case []interface{}:
		r = ArrayToStrSlice(obj.([]interface{}))
	default:
		r = append(r, fmt.Sprintf("%v", obj))
		break
	}
	return r
}

func ArrayToStrSlice(obj []interface{}) []string {
	var r []string
	if obj == nil {
		return r
	}
	for i := 0; i < len(obj); i++ {
		r = append(r, ToString(obj[i]))
	}
	return r
}

func ObjByAnchor(obj interface{}, anchor string) interface{} {
	if obj == nil {
		return nil
	}
	return mapUtil.ObjByAnchorFromMap(obj, strings.Split(anchor, ":"))
}

// HasObj 判断 target 中是否包含 obj
// !objUtil.HasObj(ext[1], [...]string{"jpg", "png"})
func HasObj(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}

// SimpleCopyProperties 对象copy
// target 目标对象指针
// src 源对象(或源对象指针)
func SimpleCopyProperties(target, src interface{}) (err error) {
	// 防止意外panic
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	dstType, dstValue := reflect.TypeOf(target), reflect.ValueOf(target)
	srcType, srcValue := reflect.TypeOf(src), reflect.ValueOf(src)

	// dst必须结构体指针类型
	if dstType.Kind() != reflect.Ptr || dstType.Elem().Kind() != reflect.Struct {
		return errors.New("dst type should be a struct pointer")
	}

	// src必须为结构体或者结构体指针，.Elem()类似于*ptr的操作返回指针指向的地址反射类型
	if srcType.Kind() == reflect.Ptr {
		srcType, srcValue = srcType.Elem(), srcValue.Elem()
	}
	if srcType.Kind() != reflect.Struct {
		return errors.New("src type should be a struct or a struct pointer")
	}

	// 取具体内容
	dstType, dstValue = dstType.Elem(), dstValue.Elem()

	// 属性个数
	propertyNums := dstType.NumField()

	for i := 0; i < propertyNums; i++ {
		// 属性
		property := dstType.Field(i)
		// 待填充属性值
		propertyValue := srcValue.FieldByName(property.Name)

		// 无效，说明src没有这个属性 || 属性同名但类型不同
		if !propertyValue.IsValid() || property.Type != propertyValue.Type() {
			continue
		}

		if dstValue.Field(i).CanSet() {
			dstValue.Field(i).Set(propertyValue)
		}
	}
	return nil
}

// 此时 T 不能为指针
func ParseList[T any](list []interface{}) ([]T, error) {
	r := make([]T, 0, len(list))
	for _, o := range list {
		s := strUtil.ToJsonStr(o)
		n := new(T)
		if err := strUtil.Parse(s, n); err != nil {
			return nil, err
		}
		r = append(r, *n)
	}
	return r, nil
}

// 判断 interface{} 是否为空
// 输入必须是 chan, func, interface, map, pointer, or slice
func IsNil(i interface{}) bool {
	defer func() {
		recover()
	}()
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr || vi.Kind() == reflect.Func {
		return vi.IsNil()
	}
	return false
}
