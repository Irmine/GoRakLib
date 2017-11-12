package protocol

import (
	"fmt"
)

func Read(buffer *[]byte, offset *int, length int) []byte {
	bytes := make([]byte, 0)
	if *offset == (len( *buffer) - 1) {
		fmt.Printf("An error occurred: %v", "no bytes left to Write")
		panic("Aborting...")
	}
	if length > 1 {
		for i := 0; i < length; i++ {
			bytes = append(bytes, (*buffer)[*offset])
			*offset++
		}
		*offset++
		return bytes
	}
	bytes = append(bytes, (*buffer)[*offset])
	*offset++
	return bytes
}

func Write(buffer *[]byte, v byte){
	*buffer = append(*buffer, v)
}

func WriteBool(buffer *[]byte, bool bool) {
	if bool {
		WriteByte(buffer, 0x01)
		return
	}
	WriteByte(buffer, 0x00)
}

func ReadBool(buffer *[]byte, offset *int) bool {
	out := Read(buffer, offset, 1)
	return out[0] != 0x00
}

func WriteByte(buffer *[]byte, byte byte) {
	Write(buffer, byte)
}

func ReadByte(buffer *[]byte, offset *int) byte {
	out := Read(buffer, offset, 1)
	return byte(out[0])
}

func WriteUnsignedByte(buffer *[]byte, unsigned uint8) {
	WriteByte(buffer, byte(unsigned))
}

func ReadUnsignedByte(buffer *[]byte, offset *int) byte {
	out := Read(buffer, offset, 1)
	return byte(out[0])
}

func WriteShort(buffer *[]byte, signed int16) {
	var i uint
	len2 := 2
	for i = 0; i < uint(len2) * 8; i += 8 {
		Write(buffer, byte(signed >> i))
	}
}

func ReadShort(buffer *[]byte, offset *int) int16 {
	var v int
	var i uint
	var out int
	bytes := Read(buffer, offset, 2)
	len2 := len(bytes)
	v = len2
	for i = 0; i < uint(len2) * 8; i += 8 {
		if i == 0 {
			out = int(bytes[v])
			continue
		}
		out |= int(bytes[v]) << i
		v--
	}
	return int16(out)
}

func WriteLittleEndianShort(buffer *[]byte, signed int16) {
	var i uint
	for i = 2 * 8; i > 0; i -= 8 {
		Write(buffer, byte(signed >> i))
	}
}

func ReadLittleEndianShort(buffer *[]byte, offset *int) int16 {
	var v int
	var i uint
	var out int
	bytes := Read(buffer, offset, 2)
	len2 := len(bytes)
	v = len2
	for i = uint(len2) * 8; i > 0; i -= 8 {
		if i == 0 {
			out = int(bytes[v])
			continue
		}
		out |= int(bytes[v]) << i
		v--
	}
	return int16(out)
}

func WriteInt(buffer *[]byte, int int32) {
	var i uint
	len2 := 4
	for i = 0; i < uint(len2) * 8; i += 8 {
		Write(buffer, byte(int >> i))
	}
}

func ReadInt(buffer *[]byte, offset *int) int32 {
	var v int
	var i uint
	var out int
	bytes := Read(buffer, offset, 4)
	len2 := len(bytes)
	v = len2
	for i = 0; i < uint(len2) * 8; i += 8 {
		if i == 0 {
			out = int(bytes[v])
			continue
		}
		out |= int(bytes[v]) << i
		v--
	}
	return int32(out)
}

func WriteLong(buffer *[]byte, int int64) {
	var i uint
	len2 := 8
	for i = 0; i < uint(len2) * 8; i += 8 {
		Write(buffer, byte(int >> i))
	}
}

func ReadLong(buffer *[]byte, offset *int) int64 {
	var v int
	var i uint
	var out int
	bytes := Read(buffer, offset, 8)
	len2 := len(bytes)
	v = len2
	for i = 0; i < uint(len2) * 8; i += 8 {
		if i == 0 {
			out = int(bytes[v - 1])
			continue
		}
		out |= int(bytes[v - 1]) << i
		v--
	}
	return int64(out)
}

func WriteUnsignedLong(buffer *[]byte, int uint64) {
	var i uint
	len2 := 8
	for i = 0; i < uint(len2) * 8; i += 8 {
		Write(buffer, byte(int >> i))
	}
}

func ReadUnsignedLong(buffer *[]byte, offset *int) uint64 {
	var v int
	var i uint
	var out int
	bytes := Read(buffer, offset, 8)
	len2 := len(bytes)
	v = len2
	for i = 0; i < uint(len2) * 8; i += 8 {
		if i == 0 {
			out = int(bytes[v])
			continue
		}
		out |= int(bytes[v]) << i
		v--
	}
	return uint64(out)
}

func WriteFloat(buffer *[]byte, float float32) {
	var i uint
	len2 := 4
	for i = 0; i < uint(len2) * 8; i += 8 {
		Write(buffer, byte(uint(float) >> i))
	}
}

func ReadFloat(buffer *[]byte, offset *int) float32 {
	var v int
	var i uint
	var out int
	bytes := Read(buffer, offset, 4)
	len2 := len(bytes)
	v = len2
	for i = 0; i < uint(len2) * 8; i += 8 {
		if i == 0 {
			out = int(bytes[v])
			continue
		}
		out |= int(bytes[v]) << i
		v--
	}
	return float32(out)
}

func WriteDouble(buffer *[]byte, double float64) {
	var i uint
	len2 := 4
	for i = 0; i < uint(len2) * 8; i += 8 {
		Write(buffer, byte(uint(double) >> i))
	}
}

func ReadDouble(buffer *[]byte, offset *int) float64 {
	var v int
	var i uint
	var out int
	bytes := Read(buffer, offset, 8)
	len2 := len(bytes)
	v = len2
	for i = 0; i < uint(len2) * 8; i += 8 {
		if i == 0 {
			out = int(bytes[v])
			continue
		}
		out |= int(bytes[v]) << i
		v--
	}
	return float64(out)
}

func WriteString(buffer *[]byte, string string) {
	len2 := len(string)
	WriteVarInt(buffer, int32(len2))
	for i := 0; i < len2; i++ {
		WriteByte(buffer, byte(string[i]))
	}
}

func ReadString(buffer *[]byte, offset *int) string {
	bytes := Read(buffer, offset, int(ReadVarInt(buffer, offset)))
	return string(bytes)
}

func WriteVarInt(buffer *[]byte, int int32) {
	var int2 uint32
	for int2 != 0 {
		out := int & 0x7F
		int2 = uint32(int) >> 7
		if int2 != 0 {
			out |= 0x7F
		}
		WriteByte(buffer, byte(out))
	}
}

func ReadVarInt(buffer *[]byte, offset *int) int32 {
	var out int32
	var next byte
	var bytesRead int32

	for (next & 0x7F) != 0 {
		next = ReadByte(buffer, offset)
		out |= int32(next & 0x7F) << 7 * bytesRead
		bytesRead++
		if bytesRead > 5 {
			fmt.Printf("An error occurred: var int is too big")
			panic("Aborting...")
		}
	}

	return out
}

func WriteVarLong(buffer *[]byte, int int64) {
	var int2 uint64
	for int2 != 0 {
		out := int & 0x7F
		int2 = uint64(int) >> 7
		if int2 != 0 {
			out |= 0x7F
		}
		WriteByte(buffer, byte(out))
	}
}

func ReadVarLong(buffer *[]byte, offset *int) int64 {
	var out int64
	var next byte
	var bytesRead int64

	for (next & 0x7F) != 0 {
		next = ReadByte(buffer, offset)
		out |= int64(next & 0x7F) << 7 * bytesRead
		bytesRead++
		if bytesRead > 10 {
			fmt.Printf("An error occurred: var long is too big")
			panic("Aborting...")
		}
	}

	return out
}

func ReadBigEndianTriad(buffer *[]byte, offset *int) uint32 {
	var out uint32
	var bytes = Read(buffer, offset, 3)
	out = uint32(bytes[0] | (bytes[1] << 8) | (bytes[2] << 16))

	return out
}

func WriteBigEndianTriad(buffer *[]byte, uint uint32) {
	Write(buffer, byte(uint & 0xFF))
	Write(buffer, byte(uint >> 8 & 0xFF))
	Write(buffer, byte(uint >> 16))
}

func ReadLittleEndianTriad(buffer *[]byte, offset *int) uint32 {
	var out uint32
	var bytes = Read(buffer, offset, 3)
	out = uint32(bytes[2] | (bytes[1] << 8) | (bytes[0] << 16))

	return out
}

func WriteLittleEndianTriad(buffer *[]byte, uint uint32) {
	Write(buffer, byte(uint >> 16))
	Write(buffer, byte(uint >> 8 & 0xFF))
	Write(buffer, byte(uint & 0xFF))
}