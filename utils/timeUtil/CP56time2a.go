package timeUtil

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/vincent78/butil/utils/byteUtil"
)

/*


解析下面这个CP56time2a表示的时间

14 2 d3 d 23 46 50
首先按照上面解析规则的纵坐标，标好位置

第n个字节	解析数组
1			50
2			46
3			23
4			d
5			d3
6			2
7			14

第1-2个字节表示ms，首先将这俩合并，结果是0x4650，然后将其转换成十进制，则为18000，单位是ms，转换成s，即18s
第3个字节的第0位到第五位表示分钟。
先将第3个字节 0x23 转换成 二进制 0010 0011 ，从0开始，截取到第五位是 10 0011，将其转换成十进制是35，即 35min。
第4个字节的第0位到第4位表示小时。
先将第4个字节 0x0d 转换成 二进制 0000 1101 ，从0开始，截取到第四位是 01101，将其转换成十进制是13，即 13h。
综上，按照前四个字节解析出来的时分秒是13:36:18。
年月日的解析也是按照上面的路子来，就不在赘述了，最后说下最终的解析结果是
2020-02-19 13:36:18



注意，上面是大端的排列，实际用的是小端排列，所以下面是按小端排的进行解析的

*/

func Cp56time2aToTime(hexStr string) (time.Time, error) {
	bs, err := hex.DecodeString(hexStr)
	if err != nil {
		return time.Now(), err
	}
	if len(bs) != 7 {
		return time.Now(), errors.New("the length is not equal 7")
	}
	ms := byteUtil.Bytes2IntByOrder([]byte{bs[0], bs[1], 0, 0}, binary.LittleEndian)

	min := bs[2]
	byteUtil.BitSet0(&min, 7)
	byteUtil.BitSet0(&min, 6)

	hour := bs[3]
	byteUtil.BitSet0(&hour, 7)
	byteUtil.BitSet0(&hour, 6)
	byteUtil.BitSet0(&hour, 5)

	day := bs[4]
	byteUtil.BitSet0(&day, 7)
	byteUtil.BitSet0(&day, 6)
	byteUtil.BitSet0(&day, 5)

	month := bs[5]
	byteUtil.BitSet0(&month, 7)
	byteUtil.BitSet0(&month, 6)
	byteUtil.BitSet0(&month, 5)
	byteUtil.BitSet0(&month, 4)

	year := bs[6]
	byteUtil.BitSet0(&year, 7)

	t, err := Parse(fmt.Sprintf("20%02d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, min, ms/1000))
	return t, err
}

func TimeToCp56time2a(t time.Time) string {
	bs := make([]byte, 7)
	ms := t.Second() * 1000
	msBytes := byteUtil.Int2BytesByOrder(ms, binary.LittleEndian)
	copy(bs[0:2], msBytes[0:2])
	bs[2] = byteUtil.Int2Byte(t.Minute())
	bs[3] = byteUtil.Int2Byte(t.Hour())
	bs[4] = byteUtil.Int2Byte(t.Day())
	bs[5] = byteUtil.Int2Byte(int(t.Month()))
	bs[6] = byteUtil.Int2Byte(t.Year() - 2000)
	return hex.EncodeToString(bs)
}
