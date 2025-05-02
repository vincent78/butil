package dup

import (
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

// Copy deepcopier(v0.0.0-20200430083143-45decc6639b6) 库封装
//
// dst需要创建好的对象, 支持sql.NullString, 支持WithContext
//
// src中empty, dst中的ptr不会复制. str="" -> str=nil
//
// 基础类型指针问题[UserID *uint64 -> UserID uint64], 结构在创建新的实例时,如果是基础类型(uint64)本身就是零值
// uint64 -> *uint64
// FIXME: 不能解决value -> ptr的复制, 弃用deepcopier库， 使用copier最新版本解决
//func Copy(src, dst interface{}) error {
//	return deepcopier.Copy(src).To(dst)
//}

// Copy copier封装
//
// dst需要创建好的对象
//
// value -> ptr
//
// src 中 empty, dst中的ptr会赋值为empty. str="" -> *str="" (如果需要nil不能使用这个方法)
// FIXME : 不能解决 empty -> ptr 复制后 ptr为empty, copier新版本v0.1.0 使用 option 解决
//
func Copy(src, dst interface{}) {
	err := CopyIgnoreEmpty(src, dst)
	_ = err
	return
}

func CopyWithOption(src, dst interface{}, opt copier.Option) error {

	return copier.CopyWithOption(dst, src, opt)
}

// CopyWithFieldStr 拷贝指定的属性
func CopyWithFieldStr(src, dst interface{}, fs ...string) error {
	if src == nil {
		return fmt.Errorf("src is nil")
	}
	if dst == nil {
		return fmt.Errorf("dst is nil")
	}

	if len(fs) == 0 {
		return fmt.Errorf("the field is empty")
	}

	stt := reflect.TypeOf(src)
	if stt.Kind() != reflect.Ptr {
		return fmt.Errorf("src must be a pointer")
	}

	dtt := reflect.TypeOf(dst)
	if dtt.Kind() != reflect.Ptr {
		return fmt.Errorf("dst must be a pointer")
	}

	sv := reflect.ValueOf(src).Elem()
	dv := reflect.ValueOf(dst).Elem()
	for _,f := range fs {
		sf := sv.FieldByName(f)
		// src没有指定的属性
		if sf == reflect.Zero(reflect.TypeOf(stt)) {
			continue
		}

		df := dv.FieldByName(f)
		// dst没有指定的属性
		if sf == reflect.Zero(reflect.TypeOf(dtt)) {
			continue
		}

		if df.CanSet() {
			df.Set(sf)
		}
	}
	return nil
}

func CopyIgnoreEmpty(src, dst interface{}) error {

	return CopyWithOption(src, dst, copier.Option{IgnoreEmpty: true})
}

func CopyMap(src, dst interface{}) {
	err := copier.Copy(dst, src)
	_ = err
	return
}

func Decode(src interface{}, dst interface{}) error {
	return mapstructure.Decode(src, dst)
}

func DecodeSpec(data interface{}, out interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:   nil,
		//DecodeHook: ToTimeHookFunc(),
		Result:     out,
	})
	if err != nil {
		return err
	}
	if err := decoder.Decode(data); err != nil {
		return err
	}
	return nil
}

func DecodeMap(data map[string]interface{}, out interface{}) error {
	config := &mapstructure.DecoderConfig{
		//DecodeHook:       mapstructure.StringToTimeHookFunc(timex.LayoutLong),
		//DecodeHook:       ToTimeHookFunc(),
		WeaklyTypedInput: true,
		Result:           out,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}

	err = decoder.Decode(data)
	if err != nil {
		return err
	}
	return nil
}
