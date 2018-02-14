package binary

type BinaryStream struct {
	Offset int
	Buffer []byte
}

func NewStream() *BinaryStream {
	return &BinaryStream{0, []byte{}}
}

func (stream *BinaryStream) SetBuffer(Buffer []byte) {
	stream.Buffer = Buffer
}

func (stream *BinaryStream) GetBuffer() []byte {
	return stream.Buffer
}

func (stream *BinaryStream) Feof() bool {
	return stream.Offset >= len(stream.Buffer)-1
}

func (stream *BinaryStream) Get(length int) []byte {
	if length < 0 {
		length = len(stream.Buffer) - stream.Offset
	}
	return Read(&stream.Buffer, &stream.Offset, length)
}

func (stream *BinaryStream) PutBool(v bool) {
	WriteBool(&stream.Buffer, v)
}

func (stream *BinaryStream) GetBool() bool {
	return ReadBool(&stream.Buffer, &stream.Offset)
}

func (stream *BinaryStream) PutByte(v byte) {
	WriteByte(&stream.Buffer, v)
}

func (stream *BinaryStream) GetByte() byte {
	return ReadByte(&stream.Buffer, &stream.Offset)
}

func (stream *BinaryStream) PutUnsignedByte(v byte) {
	WriteUnsignedByte(&stream.Buffer, v)
}

func (stream *BinaryStream) GetUnsignedByte() byte {
	return ReadUnsignedByte(&stream.Buffer, &stream.Offset)
}

func (stream *BinaryStream) PutShort(v int16) {
	WriteShort(&stream.Buffer, v)
}

func (stream *BinaryStream) GetShort() int16 {
	return ReadShort(&stream.Buffer, &stream.Offset)
}

func (stream *BinaryStream) PutUnsignedShort(v uint16) {
	WriteUnsignedShort(&stream.Buffer, v)
}

func (stream *BinaryStream) GetUnsignedShort() uint16 {
	return ReadUnsignedShort(&stream.Buffer, &stream.Offset)
}

func (stream *BinaryStream) PutInt(v int32) {
	WriteInt(&stream.Buffer, v)
}

func (stream *BinaryStream) GetInt() int32 {
	return ReadInt(&stream.Buffer, &stream.Offset)
}

func (stream *BinaryStream) PutLong(v int64) {
	WriteLong(&stream.Buffer, v)
}

func (stream *BinaryStream) GetLong() int64 {
	return ReadLong(&stream.Buffer, &stream.Offset)
}

func (stream *BinaryStream) PutUnsignedLong(v uint64) {
	WriteUnsignedLong(&stream.Buffer, v)
}

func (stream *BinaryStream) GetUnsignedLong() uint64 {
	return ReadUnsignedLong(&stream.Buffer, &stream.Offset)
}

func (stream *BinaryStream) PutFloat(v float32) {
	WriteFloat(&stream.Buffer, v)
}

func (stream *BinaryStream) GetFloat() float32 {
	return ReadFloat(&stream.Buffer, &stream.Offset)
}

func (stream *BinaryStream) PutDouble(v float64) {
	WriteDouble(&stream.Buffer, v)
}

func (stream *BinaryStream) GetDouble() float64 {
	return ReadDouble(&stream.Buffer, &stream.Offset)
}

func (stream *BinaryStream) PutString(v string) {
	WriteUnsignedShort(&stream.Buffer, uint16(len(v)))
	stream.Buffer = append(stream.Buffer, []byte(v)...)
}

func (stream *BinaryStream) GetString() string {
	return string(Read(&stream.Buffer, &stream.Offset, int(stream.GetUnsignedShort())))
}

func (stream *BinaryStream) PutLittleShort(v int16) {
	WriteLittleShort(&stream.Buffer, v)
}

func (stream *BinaryStream) GetLittleShort() int16 {
	return ReadLittleShort(&stream.Buffer, &stream.Offset)
}

func (stream *BinaryStream) PutLittleUnsignedShort(v uint16) {
	WriteLittleUnsignedShort(&stream.Buffer, v)
}

func (stream *BinaryStream) GetLittleUnsignedShort() uint16 {
	return ReadLittleUnsignedShort(&stream.Buffer, &stream.Offset)
}

func (stream *BinaryStream) PutLittleInt(v int32) {
	WriteLittleInt(&stream.Buffer, v)
}

func (stream *BinaryStream) GetLittleInt() int32 {
	return ReadLittleInt(&stream.Buffer, &stream.Offset)
}

func (stream *BinaryStream) PutLittleLong(v int64) {
	WriteLittleLong(&stream.Buffer, v)
}

func (stream *BinaryStream) GetLittleLong() int64 {
	return ReadLittleLong(&stream.Buffer, &stream.Offset)
}

func (stream *BinaryStream) PutLittleUnsignedLong(v uint64) {
	WriteLittleUnsignedLong(&stream.Buffer, v)
}

func (stream *BinaryStream) GetLittleUnsignedLong() uint64 {
	return ReadLittleUnsignedLong(&stream.Buffer, &stream.Offset)
}

func (stream *BinaryStream) PutLittleFloat(v float32) {
	WriteLittleFloat(&stream.Buffer, v)
}

func (stream *BinaryStream) GetLittleFloat() float32 {
	return ReadLittleFloat(&stream.Buffer, &stream.Offset)
}

func (stream *BinaryStream) PutLittleDouble(v float64) {
	WriteDouble(&stream.Buffer, v)
}

func (stream *BinaryStream) GetLittleDouble() float64 {
	return ReadDouble(&stream.Buffer, &stream.Offset)
}

func (stream *BinaryStream) PutTriad(v uint32) {
	WriteBigTriad(&stream.Buffer, v)
}

func (stream *BinaryStream) GetTriad() uint32 {
	return ReadBigTriad(&stream.Buffer, &stream.Offset)
}

func (stream *BinaryStream) PutLittleTriad(v uint32) {
	WriteLittleTriad(&stream.Buffer, v)
}

func (stream *BinaryStream) GetLittleTriad() uint32 {
	return ReadLittleTriad(&stream.Buffer, &stream.Offset)
}

func (stream *BinaryStream) PutBytes(bytes []byte) {
	stream.Buffer = append(stream.Buffer, bytes...)
	stream.Offset += len(bytes)
}

func (stream *BinaryStream) ResetStream() {
	stream.Offset = 0
	stream.Buffer = []byte{}
}
