package codec

type ByteCodecConfig struct {
	Key        string                `json:"key"`        // 转换标识
	Total      uint                  `json:"total"`      //总长度 字节的长度
	StructName string                `json:"structName"` // 结构体的全路径
	DspVersion string                `json:"dspVersion"` // 下位机的版本
	Items      []ByteCodecItemConfig `json:"items"`      //每项的设置
}

type ByteCodecItemConfig struct {
	Name       string `json:"name"`       // 字段名   TODO:如果有层级，以逗号分隔
	Index      uint   `json:"index"`      // 初始位置(从0开始）
	Size       uint   `json:"size"`       // 长度 字节数 就是几个字节
	EncodeFunc string `json:"encodeFunc"` // 编码的方法名
	DecodeFunc string `json:"decodeFunc"` // 解码的方法名
	Desc       string `json:"desc"`       //解释
}
