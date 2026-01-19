package reflectx

import "reflect"

// SafeValue
// 可以传指针，也可以传值
// 传指针则返回指针，传值则返回值
// 遇到空指针会返回该类型的零值的指针
func SafeValue[T any](v T) T {
	vt := reflect.TypeOf(v)
	switch vt.Kind() {
	case reflect.Ptr:
		vv := reflect.ValueOf(v)
		if vv.IsNil() {
			return reflect.New(vt.Elem()).Interface().(T)
		}
		return v
	default:
		return v
	}
}

// SafeStruct
// 必须传入指针
// 它会将空指针的字段设置为该类型的零值的指针
// 传入空指针会返回一个所有字段都为零值的结构体
func SafeStruct(v any) any {
	vt := reflect.TypeOf(v)
	vv := reflect.ValueOf(v)
	if vv.IsNil() {
		vv = reflect.ValueOf(reflect.New(vt.Elem()).Interface()).Elem()
	} else {
		vv = vv.Elem()
	}
	for i := 0; i < vv.NumField(); i++ {
		if vv.Field(i).Kind() == reflect.Ptr && vv.Field(i).IsNil() {
			vv.Field(i).Set(reflect.New(vv.Field(i).Type().Elem()))
		}
	}
	return vv.Interface()
}

// IsPtr
// 判断当前是否为指针
func IsPtr[T any](v T) bool {
	vt := reflect.TypeOf(v)
	return vt.Kind() == reflect.Ptr
}
