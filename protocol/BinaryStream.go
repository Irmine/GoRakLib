package protocol

type BinaryStream struct {
	offset int
	buffer []byte
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

func (stream *BinaryStream) Reset() {
	stream.offset = 0
	stream.buffer = []byte{}
}
