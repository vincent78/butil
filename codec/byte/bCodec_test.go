package codec

import (
	"encoding/json"
	"fmt"
	"testing"
)

type PdInfo struct {
	RobotState   uint8 `json:"robotState"` //小车状态
	OrderState   uint8 `json:"orderState"` //
	RobotX       int32 `json:"x"`          //小车当前坐标X
	RobotY       int32 `json:"y"`          //小车当前从标Y
	RobotHeading int16 `json:"heading"`    //小车当前角度
}

var testData1 = `
{
    "robotState": 192,
    "orderState": 11,
    "x": 2000,
    "y": 10000,
    "heading": 10
}
`

var configStr = `
{
    "key": "pdinfo",
    "total": 12,
    "items":
    [
        {
            "name": "RobotState",
            "index": 0,
            "size": 1,
            "EncodeFunc": "Int2BytesFunc",
            "DecodeFunc": "Bytes2IntFunc",
            "Desc": "测试字段1"
        },
        {
            "name": "OrderState",
            "index": 1,
            "size": 1,
            "EncodeFunc": "Int2BytesFunc",
            "DecodeFunc": "Bytes2IntFunc",
            "Desc": "测试字段2"
        },
        {
            "name": "RobotX",
            "index": 2,
            "size": 4,
            "EncodeFunc": "Int2BytesFunc",
            "DecodeFunc": "Bytes2IntFunc",
            "Desc": "测试字段3"
        },
        {
            "name": "RobotY",
            "index": 6,
            "size": 4,
            "EncodeFunc": "Int2BytesFunc",
            "DecodeFunc": "Bytes2IntFunc",
            "Desc": "测试字段4"
        },
        {
            "name": "RobotHeading",
            "index": 10,
            "size": 2,
            "EncodeFunc": "Int2BytesFunc",
            "DecodeFunc": "Bytes2IntFunc",
            "Desc": "测试字段5"
        }
    ]
}
`

func TestByteEncode(t *testing.T) {
	p := &PdInfo{}
	e := json.Unmarshal([]byte(testData1), p)
	if e != nil {
		fmt.Printf("PdInfo parse error : %v\n", e)
	}

	config := ByteCodecConfig{}
	e = json.Unmarshal([]byte(configStr), &config)
	if e != nil {
		fmt.Printf("ByteCodecConfig parse error : %v\n", e)
	}
	bs, e := ByteEncode(p, config)
	fmt.Printf("the result is : %v", bs)
}

func TestByteDecode(t *testing.T) {
	p := &PdInfo{}
	e := json.Unmarshal([]byte(testData1), p)
	if e != nil {
		fmt.Printf("PdInfo parse error : %v\n", e)
	}
	bs := []byte{192, 11, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	config := ByteCodecConfig{}
	e = json.Unmarshal([]byte(configStr), &config)
	if e != nil {
		fmt.Printf("ByteCodecConfig parse error : %v\n", e)
	}
	e = ByteDecode(bs, p, config)
	if e != nil {
		fmt.Printf("ByteDecode error : %v\n", e)
	}
	fmt.Printf("the result is : %+v", p)
}
