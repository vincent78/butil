package codec

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/vincent78/butil/utils/byteUtil"
	"github.com/vincent78/butil/utils/reflectx"
)

/************************************************************************
 *
 * Encode相关
 *
 ************************************************************************/

func Int2BytesFunc(v any, config ByteCodecItemConfig) ([]byte, error) {
	if t, ok := v.(int64); ok {
		return byteUtil.Int2BytesBigEndian(t, byte(config.Size))
	} else if t, ok := v.(uint64); ok {
		return byteUtil.Int2BytesBigEndian(int64(t), byte(config.Size))
	} else if t, ok := v.(int32); ok {
		return byteUtil.Int2BytesBigEndian(int64(t), byte(config.Size))
	} else if t, ok := v.(uint32); ok {
		return byteUtil.Int2BytesBigEndian(int64(t), byte(config.Size))
	} else if t, ok := v.(int16); ok {
		return byteUtil.Int2BytesBigEndian(int64(t), byte(config.Size))
	} else if t, ok := v.(uint16); ok {
		return byteUtil.Int2BytesBigEndian(int64(t), byte(config.Size))
	} else if t, ok := v.(int8); ok {
		return byteUtil.Int2BytesBigEndian(int64(t), byte(config.Size))
	} else if t, ok := v.(uint8); ok {
		return byteUtil.Int2BytesBigEndian(int64(t), byte(config.Size))
	} else if t, ok := v.(string); ok {
		n, e := strconv.Atoi(t)
		if e != nil {
			return nil, e
		}
		return byteUtil.Int2BytesBigEndian(int64(n), byte(config.Size))
	} else if t, ok := v.(bool); ok {
		if t {
			return byteUtil.Int2BytesBigEndian(1, byte(config.Size))
		} else {
			return byteUtil.Int2BytesBigEndian(0, byte(config.Size))
		}
	} else {
		return nil, errors.New("the value type is not catch")
	}
}

func GetEncodeFunc(str string) func(v any, config ByteCodecItemConfig) ([]byte, error) {
	switch str {
	case "Int2BytesFunc":
		{
			return Int2BytesFunc
		}
	default:
		return nil
	}
}

func ByteEncode(m any, config ByteCodecConfig) ([]byte, error) {
	result := make([]byte, config.Total)
	for _, ic := range config.Items {
		v := reflectx.GetValueByName(m, ic.Name)
		f := GetEncodeFunc(ic.EncodeFunc)
		if f == nil {
			return nil, errors.New("the encodeFunc not catch")
		}
		bs, e := f(v, ic)
		if e != nil {
			return nil, e
		}
		for i := 0; i < int(ic.Size); i++ {
			result[i+int(ic.Index)] = bs[i]
		}
		//fmt.Printf("field: %v  %v \n", ic.Config, result)
	}
	return result, nil
}

/************************************************************************
 *
 * Decode相关
 *
 ************************************************************************/

func Bytes2IntFunc(bs []byte, config ByteCodecItemConfig) (int64, error) {
	if len(bs) != int(config.Size) {
		return 0, errors.New("the bytes size is not right")
	}

	b, e := byteUtil.Bytes2IntBigEndian(bs)
	if e != nil {
		return 0, e
	} else {
		return int64(b), e
	}

}

func GetDecodeFunc(str string) func(bs []byte, config ByteCodecItemConfig) (int64, error) {
	switch str {
	case "Bytes2IntFunc":
		{
			return Bytes2IntFunc
		}
	default:
		return nil
	}
}

func ByteDecode(bs []byte, obj any, config ByteCodecConfig) error {
	if int(config.Total) != len(bs) {
		return fmt.Errorf("the bytes length[%v] is not equal config size[%v]", len(bs), config.Total)
	}
	if !reflectx.IsPoint(obj) {
		return errors.New("the obj is not point")
	}

	for _, ic := range config.Items {
		vs := bs[ic.Index : ic.Index+ic.Size]
		v, e := byteUtil.Bytes2IntBigEndian(vs)
		if e != nil {
			return e
		}
		f, e := reflectx.GetFieldByName(obj, ic.Name)
		if e != nil {
			return fmt.Errorf("the field [%v] error: %v", ic.Name, e)
		}

		if f.Type.Kind() == reflect.String {
			reflectx.SetFieldValue(obj, ic.Name, reflect.ValueOf(strconv.Itoa(v)))
		} else {
			if ic.Size == 1 {
				if f.Type.Kind() == reflect.Uint8 {
					reflectx.SetFieldValue(obj, ic.Name, reflect.ValueOf(uint8(v)))
				} else if f.Type.Kind() == reflect.Int8 {
					reflectx.SetFieldValue(obj, ic.Name, reflect.ValueOf(int8(v)))
				} else {
					return fmt.Errorf("unknow field[%v] type[%v] in size 1", ic.Name, f.Type.Kind())
				}
			} else if ic.Size == 2 {
				if f.Type.Kind() == reflect.Uint16 {
					reflectx.SetFieldValue(obj, ic.Name, reflect.ValueOf(uint16(v)))
				} else if f.Type.Kind() == reflect.Int16 {
					reflectx.SetFieldValue(obj, ic.Name, reflect.ValueOf(int16(v)))
				} else {
					return fmt.Errorf("unknow field[%v] type[%v] in size 2", ic.Name, f.Type.Kind())
				}
			} else if ic.Size > 2 && ic.Size <= 4 {
				if f.Type.Kind() == reflect.Uint32 {
					reflectx.SetFieldValue(obj, ic.Name, reflect.ValueOf(uint32(v)))
				} else if f.Type.Kind() == reflect.Int32 {
					reflectx.SetFieldValue(obj, ic.Name, reflect.ValueOf(int32(v)))
				} else {
					return fmt.Errorf("unknow field[%v] type[%v] in size 4", ic.Name, f.Type.Kind())
				}
			} else if ic.Size > 4 && ic.Size <= 8 {
				if f.Type.Kind() == reflect.Uint64 {
					reflectx.SetFieldValue(obj, ic.Name, reflect.ValueOf(uint64(v)))
				} else if f.Type.Kind() == reflect.Uint64 {
					reflectx.SetFieldValue(obj, ic.Name, reflect.ValueOf(int64(v)))
				} else {
					return fmt.Errorf("unknow field[%v] type[%v] in size 8", ic.Name, f.Type.Kind())
				}
			}
		}
	}
	return nil
}
