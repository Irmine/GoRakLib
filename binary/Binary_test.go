package binary

import (
	"testing"
	"math"
)

var correct =
	"0100ABCDFF637FFF8000FFFF00787FFFFFFFF1DA38F47FFFFFFFFFFFFFFF8000000000000000FFFFFFFFFFFF" +
		"FFFF000000003ADE68B140490FDB7F7FFFFF4005BF0A8B1457697FEFFFFFFFFFFFFFFF7FC7CFFFFFD204FFFF" +
		"FF7F85FFFFFFFFFFFFFFFFFFFF7FEB81CE64CDE633FBFFFFFFFFFFFFFFFF7FC9120000000000DB0F4940FFFF" +
		"7F7F6957148B0ABF0540FFFFFFFFFFFFEF7F0000000A676F2072616B206C696201E24000000140E201010000"

func TestBinaryWrite(t *testing.T) {
	buf := new([]byte)

	doWriteTest(buf)

	if bytesToHex(buf) != correct {
		t.Error("Incorrect buffer writing")
	}
}

func TestBinaryRead(t *testing.T) {
	buf := hexToBytes(correct)

	off := 0
	if ReadBool					(buf, &off) != true { t.Error("Incorrect buffer reading: ReadBool")}
	if ReadBool					(buf, &off) != false { t.Error("Incorrect buffer reading: ReadBool")}
	if ReadByte					(buf, &off) != 0xab { t.Error("Incorrect buffer reading: ReadByte")}
	if ReadByte					(buf, &off) != 0xcd { t.Error("Incorrect buffer reading: ReadByte")}
	if ReadUnsignedByte			(buf, &off) != 255 { t.Error("Incorrect buffer reading: ReadUnsignedByte")}
	if ReadUnsignedByte			(buf, &off) != 99 { t.Error("Incorrect buffer reading: ReadUnsignedByte")}
	if ReadShort				(buf, &off) != math.MaxInt16 { t.Error("Incorrect buffer reading: ReadShort")}
	if ReadShort				(buf, &off) != -32768 { t.Error("Incorrect buffer reading: ReadShort")}
	if ReadUnsignedShort		(buf, &off) != math.MaxUint16 { t.Error("Incorrect buffer reading: ReadUnsignedShort")}
	if ReadUnsignedShort		(buf, &off) != 120 { t.Error("Incorrect buffer reading: ReadUnsignedShort")}
	if ReadInt					(buf, &off) != math.MaxInt32 { t.Error("Incorrect buffer reading: ReadInt")}
	if ReadInt					(buf, &off) != -237356812 { t.Error("Incorrect buffer reading: ReadInt")}
	if ReadLong					(buf, &off) != math.MaxInt64 { t.Error("Incorrect buffer reading: ReadLong")}
	if ReadLong					(buf, &off) != -9223372036854775808 { t.Error("Incorrect buffer reading: ReadLong")}
	if ReadUnsignedLong			(buf, &off) != math.MaxUint64 { t.Error("Incorrect buffer reading: ReadUnsignedLong")}
	if ReadUnsignedLong			(buf, &off) != 987654321 { t.Error("Incorrect buffer reading: ReadUnsignedLong")}
	if ReadFloat				(buf, &off) != math.Pi { t.Error("Incorrect buffer reading: ReadFloat")}
	if ReadFloat				(buf, &off) != math.MaxFloat32 { t.Error("Incorrect buffer reading: ReadFloat")}
	if ReadDouble				(buf, &off) != math.E { t.Error("Incorrect buffer reading: ReadDouble")}
	if ReadDouble				(buf, &off) != math.MaxFloat64 { t.Error("Incorrect buffer reading: ReadDouble")}
	if ReadLittleShort			(buf, &off) != math.MaxInt16 { t.Error("Incorrect buffer reading: ReadLittleShort")}
	if ReadLittleShort			(buf, &off) != -12345 { t.Error("Incorrect buffer reading: ReadLittleShort")}
	if ReadLittleUnsignedShort	(buf, &off) != math.MaxUint16 { t.Error("Incorrect buffer reading: ReadLittleUnsignedShort")}
	if ReadLittleUnsignedShort	(buf, &off) != 1234 { t.Error("Incorrect buffer reading: ReadLittleUnsignedShort")}
	if ReadLittleInt			(buf, &off) != math.MaxInt32 { t.Error("Incorrect buffer reading: ReadLittleInt")}
	if ReadLittleInt			(buf, &off) != -123 { t.Error("Incorrect buffer reading: ReadLittleInt")}
	if ReadLittleLong			(buf, &off) != math.MaxInt64 { t.Error("Incorrect buffer reading: ReadLittleLong")}
	if ReadLittleLong			(buf, &off) != -345678976543456789 { t.Error("Incorrect buffer reading: ReadLittleLong")}
	if ReadLittleUnsignedLong	(buf, &off) != math.MaxUint64 { t.Error("Incorrect buffer reading: ReadLittleUnsignedLong")}
	if ReadLittleUnsignedLong	(buf, &off) != 1231231 { t.Error("Incorrect buffer reading: ReadLittleUnsignedLong")}
	if ReadLittleFloat			(buf, &off) != math.Pi { t.Error("Incorrect buffer reading: ReadLittleFloat")}
	if ReadLittleFloat			(buf, &off) != math.MaxFloat32 { t.Error("Incorrect buffer reading: ReadLittleFloat")}
	if ReadLittleDouble			(buf, &off) != math.E { t.Error("Incorrect buffer reading: ReadLittleDouble")}
	if ReadLittleDouble			(buf, &off) != math.MaxFloat64 { t.Error("Incorrect buffer reading: ReadLittleDouble")}
	if ReadString				(buf, &off) != "" { t.Error("Incorrect buffer reading: ReadString")}
	if ReadString				(buf, &off) != "go rak lib" { t.Error("Incorrect buffer reading: ReadString")}
	if ReadBigEndianTriad		(buf, &off) != 123456 { t.Error("Incorrect buffer reading: ReadBigEndianTriad")}
	if ReadBigEndianTriad		(buf, &off) != 0x1 { t.Error("Incorrect buffer reading: ReadBigEndianTriad")}
	if ReadLittleEndianTriad	(buf, &off) != 123456 { t.Error("Incorrect buffer reading: ReadLittleEndianTriad")}
	if ReadLittleEndianTriad	(buf, &off) != 0x1 { t.Error("Incorrect buffer reading: ReadLittleEndianTriad")}

}

func BenchmarkWrite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := new([]byte)
		doWriteTest(buf)
	}
}

func doWriteTest(buf *[]byte) {
	WriteBool					(buf, true)
	WriteBool					(buf, false)
	WriteByte					(buf, 0xab)
	WriteByte					(buf, 0xcd)
	WriteUnsignedByte			(buf, 255)
	WriteUnsignedByte			(buf, 99)
	WriteShort					(buf, math.MaxInt16)
	WriteShort					(buf, -32768)
	WriteUnsignedShort			(buf, math.MaxUint16)
	WriteUnsignedShort			(buf, 120)
	WriteInt					(buf, math.MaxInt32)
	WriteInt					(buf, -237356812)
	WriteLong					(buf, math.MaxInt64)
	WriteLong					(buf, -9223372036854775808)
	WriteUnsignedLong			(buf, math.MaxUint64)
	WriteUnsignedLong			(buf, 987654321)
	WriteFloat					(buf, math.Pi)
	WriteFloat					(buf, math.MaxFloat32)
	WriteDouble					(buf, math.E)
	WriteDouble					(buf, math.MaxFloat64)
	WriteLittleShort			(buf, math.MaxInt16)
	WriteLittleShort			(buf, -12345)
	WriteLittleUnsignedShort	(buf, math.MaxUint16)
	WriteLittleUnsignedShort	(buf, 1234)
	WriteLittleInt				(buf, math.MaxInt32)
	WriteLittleInt				(buf, -123)
	WriteLittleLong				(buf, math.MaxInt64)
	WriteLittleLong				(buf, -345678976543456789)
	WriteLittleUnsignedLong		(buf, math.MaxUint64)
	WriteLittleUnsignedLong		(buf, 1231231)
	WriteLittleFloat			(buf, math.Pi)
	WriteLittleFloat			(buf, math.MaxFloat32)
	WriteLittleDouble			(buf, math.E)
	WriteLittleDouble			(buf, math.MaxFloat64)
	WriteString					(buf, "")
	WriteString					(buf, "go rak lib")
	WriteBigEndianTriad			(buf, 123456)
	WriteBigEndianTriad			(buf, 0x1)
	WriteLittleEndianTriad		(buf, 123456)
	WriteLittleEndianTriad		(buf, 0x1)
}

var hexCode = []string{"0","1","2","3","4","5","6","7","8","9","A","B","C","D","E","F"}
var hexMap = map[uint8] byte {
	'0': 0x0,'1': 0x1,'2': 0x2,'3': 0x3,
	'4': 0x4,'5': 0x5,'6': 0x6,'7': 0x7,
	'8': 0x8,'9': 0x9,'A': 0xA,'B': 0xB,
	'C': 0xC,'D': 0xD,'E': 0xE,'F': 0xF,
}

func bytesToHex(b *[]byte) string {
	ans := ""
	for i := 0; i < len(*b); i++ {
		ans += hexCode[((*b)[i]>>4)&0x0f]
		ans += hexCode[((*b)[i])&0x0f]
	}
	return ans
}

func hexToBytes(s string) *[]byte {
	ans := *new([]byte)
	for i := 0; i < len(s); i+=2 {
		ans = append(ans, hexMap[s[i]] << 4 + hexMap[s[i+1]])
	}
	return &ans
}