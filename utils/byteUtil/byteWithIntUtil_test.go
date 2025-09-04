package byteUtil

import (
	"encoding/binary"
	"encoding/hex"
	"testing"
	"unsafe"

	"github.com/imroc/biu"
)

func TestSizeOfType(t *testing.T) {
	a1 := int(1)
	t.Logf("\n int:%d", unsafe.Sizeof(a1))
	a2 := int8(1)
	t.Logf("\n int8:%d", unsafe.Sizeof(a2))
	a3 := int16(1)
	t.Logf("\n int16:%d", unsafe.Sizeof(a3))
	a4 := int32(1)
	t.Logf("\n int32:%d", unsafe.Sizeof(a4))
	a5 := int64(1)
	t.Logf("\n int64:%d", unsafe.Sizeof(a5))

}

func TestBitGet(t *testing.T) {
	b := byte(1)
	t.Logf("byte: %v", biu.ToBinaryString(b))
	t.Logf("no.0:%v", biu.ToBinaryString(BitGet(b, 0)))
	t.Logf("no.1:%v", biu.ToBinaryString(BitGet(b, 1)))
	t.Logf("no.2:%v", biu.ToBinaryString(BitGet(b, 2)))
	t.Logf("no.3:%v", biu.ToBinaryString(BitGet(b, 3)))
	t.Logf("no.4:%v", biu.ToBinaryString(BitGet(b, 4)))
	t.Logf("no.5:%v", biu.ToBinaryString(BitGet(b, 5)))
	t.Logf("no.6:%v", biu.ToBinaryString(BitGet(b, 6)))
	t.Logf("no.7:%v", biu.ToBinaryString(BitGet(b, 7)))
}

func TestIntToBytes(t *testing.T) {
	x := 255
	t.Logf("x: %b", x)
	t.Logf("\nbig: %b   lit: %b",
		Int2BytesByOrder(x, binary.BigEndian),
		Int2BytesByOrder(x, binary.LittleEndian),
	)
}

func TestInt2Bytes(t *testing.T) {
	t.Logf("x: %d  => %v", 20000, Int2BytesByOrder(20000, binary.BigEndian))
	t.Logf("x: %d  => %v", -20000, Int2BytesByOrder(-20000, binary.BigEndian))
	t.Logf("x: %d  => %v", 9000, Int2BytesByOrder(9000, binary.BigEndian))
	t.Logf("x: %d  => %v", -9000, Int2BytesByOrder(-9000, binary.BigEndian))
	t.Logf("x: %d  => %v", 4500, Int2BytesByOrder(4500, binary.BigEndian))
	t.Logf("x: %d  => %v", -4500, Int2BytesByOrder(-4500, binary.BigEndian))
	t.Logf("x: %d  => %v", 13500, Int2BytesByOrder(13500, binary.BigEndian))
	t.Logf("x: %d  => %v", -13500, Int2BytesByOrder(-13500, binary.BigEndian))
	t.Logf("x: %d  => %v", 18000, Int2BytesByOrder(18000, binary.BigEndian))
	t.Logf("x: %d  => %v", -18000, Int2BytesByOrder(-18000, binary.BigEndian))
	t.Logf("x: %d  => %v", -2527, Int2BytesByOrder(-2527, binary.BigEndian))
}

func TestBytes2Int64(t *testing.T) {
	x := []byte{0, 0, 18, 32}
	r, _ := Bytes2IntBigEndian(x)
	t.Logf("x: %d  => %v", x, r)

	x = []byte{0, 0, 18, 92}
	r, _ = Bytes2IntBigEndian(x)
	t.Logf("x: %d  => %v", x, r)

	x = []byte{0, 0, 25, 100}
	r, _ = Bytes2IntBigEndian(x)
	t.Logf("x: %d  => %v", x, r)

	x = []byte{0, 0, 11, 184}
	r, _ = Bytes2IntBigEndian(x)
	t.Logf("x: %d  => %v", x, r)

	x = []byte{0, 0, 7, 208}
	r, _ = Bytes2IntBigEndian(x)
	t.Logf("x: %d  => %v", x, r)

	x = []byte{220, 216}
	r, _ = Bytes2IntBigEndian(x)
	t.Logf("x: %d  => %v", x, r)

	x = []byte{33, 246}
	r, _ = Bytes2IntBigEndian(x)
	t.Logf("x: %d  => %v", x, r)

}

func TestByteJoin(t *testing.T) {
	b1 := []byte("1")
	b2 := []byte("this is the test")
	b3 := Join(b1, b2)

	t.Logf("b1: %v", b1)
	t.Logf("b2: %v", b2)
	t.Logf("result: %v", b3)
}

func TestByte2Hex(t *testing.T) {
	b1 := []byte("this is the test")
	t.Logf("b1 byte: %v", b1)
	t.Logf("b1 hx: %v", hex.EncodeToString(b1))
	t.Logf("b1 hx: %x", b1)
}

func TestIntToBytes2(t *testing.T) {
	v := 10000
	b := Int2Bytes(v)
	t.Logf("the bytes is :%x", b)
}

func TestBytesToInt(t *testing.T) {
	bytes, _ := hex.DecodeString("a0860100")
	v := Bytes2IntByOrder(bytes, binary.LittleEndian)
	t.Logf("the int is %v ", v)
}

func TestBytesToInt2(t *testing.T) {
	num := 100 * 1000
	v := Int2BytesByOrder(num, binary.LittleEndian)
	t.Logf("the int is %x ", v)
}
