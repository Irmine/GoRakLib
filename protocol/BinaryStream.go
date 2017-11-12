package protocol

type BinaryStream struct {
	offset int
	buffer []byte
}

func NewStream() *BinaryStream {
	return &BinaryStream{0, make([]byte, 4096)}
}

func (stream *BinaryStream) SetBuffer(buffer []byte) {
	stream.buffer = buffer
}

func (stream *BinaryStream) Feof() bool {
	return stream.offset >= len(stream.buffer)
}

func (stream *BinaryStream) ReadBool() bool {
	return ReadBool(&stream.buffer, &stream.offset)
}

func (stream *BinaryStream) ReadByte() byte {
	return ReadByte(&stream.buffer, &stream.offset)
}

func (stream *BinaryStream) ReadShort() int16 {
	return ReadShort(&stream.buffer, &stream.offset)
}

func (stream *BinaryStream) ReadInt() int32 {
	return ReadInt(&stream.buffer, &stream.offset)
}

func (stream *BinaryStream) ReadLong() int64 {
	return ReadLong(&stream.buffer, &stream.offset)
}

func (stream *BinaryStream) ReadBigEndianTriad() uint32 {
	return ReadBigEndianTriad(&stream.buffer, &stream.offset)
}

func (stream *BinaryStream) ReadLittleEndianTriad() uint32 {
	return ReadLittleEndianTriad(&stream.buffer, &stream.offset)
}

func (stream *BinaryStream) ReadString() string {
	return ReadString(&stream.buffer, &stream.offset)
}

func (stream *BinaryStream) WriteBool(bool bool) {
	WriteBool(&stream.buffer, bool)
}

func (stream *BinaryStream) WriteByte(byte byte) {
	WriteByte(&stream.buffer, byte)
}

func (stream *BinaryStream) WriteShort(short int16) {
	WriteShort(&stream.buffer, short)
}

func (stream *BinaryStream) WriteInt(int int32) {
	WriteInt(&stream.buffer, int)
}

func (stream *BinaryStream) WriteLong(long int64) {
	WriteLong(&stream.buffer, long)
}

func (stream *BinaryStream) WrightBigEndianTriad(uint uint32) {
	WriteBigEndianTriad(&stream.buffer, uint)
}

func (stream *BinaryStream) WriteLittleEndianTriad(uint uint32) {
	WriteLittleEndianTriad(&stream.buffer, uint)
}

func (stream *BinaryStream) WriteString(string string) {
	WriteString(&stream.buffer, string)
}

func (stream *BinaryStream) WriteBytes(bytes []byte) {
	for _, byte := range bytes {
		stream.WriteByte(byte)
	}
}

func (stream *BinaryStream) ResetStream() {
	stream.offset = 0
	stream.buffer = []byte{}
}
