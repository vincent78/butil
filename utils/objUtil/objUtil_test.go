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
