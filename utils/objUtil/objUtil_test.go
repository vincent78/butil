package objUtil

import (
	"fmt"
	"testing"
)

func TestInterfaceIsNil(t *testing.T) {
	var a interface{} = nil         // tab = nil, data = nil
	var b interface{} = (*int)(nil) // tab 包含 *int 类型信息, data = nil

	t.Logf("type and data all is nil? %v", a == nil)
	t.Logf("data is nil? %v", b == nil)

}

func TestToInt(t *testing.T) {
	i, e := ToInt(1)
	if e == nil && i == 1 {
		t.Log("成功")
	} else {
		t.Errorf("转换失败:%v", e)
	}
}

func TestFunc1(t *testing.T) {
	var obj interface{}
	obj = []interface{}{
		"onte",
		1,
		"#2",
	}
	//obj = 1.2
	fmt.Printf("output:%v ", ObjToStrSlice(obj))
}

// 类型别名不能转换,支持嵌套copy
func TestSimpleCopyProperties(t *testing.T) {

	person1 := &Person1{
		Name:  "小明",
		Age:   20,
		Sex:   0,
		Hobby: []string{"听音乐", "跑步"},
		Test: Test1{
			Code:  "001",
			Color: []string{"red", "greed"},
			Count: 2,
		},
		Test1: Test1{
			Code:  "002",
			Color: []string{"yellow", "black"},
			Count: 5,
		},
		Test2: Test1{
			Code:  "002",
			Color: []string{"yellow", "black"},
			Count: 5,
		},
	}
	person2 := &Person2{}
	err := SimpleCopyProperties(person2, person1)
	if err != nil {
		return
	}
	t.Logf("person1:%#v", person1)
	t.Logf("person2:%#v", person2)
}

type Person1 struct {
	Name  Str
	Age   int
	Sex   int
	Hobby []string
	Test  Test1
	Test1 Test1
	Test2 Test1
}

type Person2 struct {
	Name  string
	Age   int
	Sex   string
	Hobby []string
	Test  Test2
	Test1 Test3
	Test2 Test1
}

type Test1 struct {
	Code  string
	Color []string
	Count int
}

type Test2 struct {
	Code  string
	Color string
	Count int
}
type Str = string
type Test3 Test1

func Test_Map2Obj(t *testing.T) {
	//str := " {\"address\":\"0x1c726d5f71ead5acf30a60912f839a752b3c6d6d\",\"topics\":[\"0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822\",\"0x000000000000000000000000a56195fe69994c083b9a86ccf63ae16402aebb7c\",\"0x000000000000000000000000a56195fe69994c083b9a86ccf63ae16402aebb7c\"],\"data\":\"0x0000000000000000000000000000000000000000000000000003ac3800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000472b5077a70c131\",\"blockNumber\":\"0x166666c\",\"transactionHash\":\"0xb14225aacee0fa2183b5bb5a3a4ad37b6b39ea5ca047c1095e1d6ae1c24f927d\",\"transactionIndex\":\"0x1\",\"blockHash\":\"0xd716e6e330283b557d7b6b75aef1c08957f32311d682cc8f85a16864c2b05bad\",\"blockTimestamp\":\"0x0\",\"logIndex\":\"0x2b\",\"removed\":false}"
	obj := &Person1{
		Name:  "person1",
		Age:   20,
		Sex:   1,
		Hobby: []string{"1", "2", "3"},
	}
	mp := make(map[string]interface{})
	if r, err := Map2Obj(mp, obj); err != nil {
		t.Errorf("transfer error: %v", err)
	} else {
		t.Logf("transfer result: %v", r)
	}
}

func Test_ObjToMap(t *testing.T) {
	obj := &Person1{
		Name:  "person1",
		Age:   20,
		Sex:   1,
		Hobby: []string{"1", "2", "3"},
		Test1: Test1{
			Code:  "001",
			Color: []string{"red", "greed"},
		},
	}
	if mp, err := ObjToMap(obj); err != nil {
		t.Errorf("transfer error: %v", err)
	} else {
		t.Logf("transfer result: %v", mp)
	}

}
