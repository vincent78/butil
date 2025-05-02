package reflectx

import (
	"errors"
	"fmt"
	"reflect"
)

// IsPoint 判断是否为指针
func IsPoint(m interface{}) bool {
	return reflect.TypeOf(m).Kind() == reflect.Ptr
}

func NilOrZero(m interface{}, fs string) bool {
	if m == nil {
		return true
	}
	//t := reflect.TypeOf(m)
	//if IsPoint(m) {}
	return false
}

func GetType(m interface{}) reflect.Type {
	tp := reflect.TypeOf(m)
	if tp.Kind() != reflect.Ptr {
		return tp
	} else {
		return tp.Elem()
	}
}

func GetFieldByName(m interface{}, n string) (reflect.StructField, error) {
	t := GetType(m)
	if t.Kind() != reflect.Struct {
		return reflect.StructField{}, errors.New("input object is not struct")
	}
	if f, exist := t.FieldByName(n); exist {
		return f, nil
	} else {
		return f, errors.New("the field is not find")
	}
}

func GetAllFieldName(m interface{}) ([]string, error) {
	t := GetType(m)
	if t.Kind() != reflect.Struct {
		return []string{}, errors.New("input object is not struct")
	}
	total := t.NumField()
	r := make([]string, 0, total)
	for i := 0; i < total; i++ {
		r = append(r, t.Field(i).Name)
	}
	return r, nil
}

func GetFieldValueByName(m interface{}, n string) reflect.Value {
	tv := reflect.ValueOf(m)
	if IsPoint(m) {
		tv = tv.Elem()
	}
	return tv.FieldByName(n)
}

func GetValueByName(m interface{}, n string) interface{} {
	tv := reflect.ValueOf(m)
	if IsPoint(m) {
		tv = tv.Elem()
	}
	return tv.FieldByName(n).Interface()
}

func SetFieldValue(m interface{}, n string, v reflect.Value) error {
	tv := reflect.ValueOf(m)
	if IsPoint(m) {
		tv = tv.Elem()
	}
	fv := tv.FieldByName(n)

	if fv.CanSet() {
		fv.Set(v)
		return nil
	} else {
		return errors.New("can't set the field")
	}
}

// 将切片转成 []interface{}
func ToInterfaceSlice(src interface{}) []interface{} {
	var ret []interface{}
	if reflect.TypeOf(src).Kind() == reflect.Slice {
		s := reflect.ValueOf(src)
		for i := 0; i < s.Len(); i++ {
			so := s.Index(i)
			ret = append(ret, so.Interface())
		}
	}
	return ret
}

// PrintInterface print the interface by level
func PrintInterface(v interface{}) {
	val := reflect.ValueOf(v).Elem()
	typ := reflect.TypeOf(v)
	//log.Printf("%+v\n", v)
	nums := val.NumField()
	for i := 0; i < nums; i++ {
		if typ.Elem().Field(i).Type.Kind() == reflect.Ptr {
			//log.Printf("%s: ", typ.Elem().Field(i).Name)
			PrintInterface(val.Field(i).Interface())
		}
	}
	fmt.Println("")
}
