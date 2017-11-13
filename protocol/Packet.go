package protocol

type Packet struct {
	packetId int
	*BinaryStream
}

type IPacket interface {
	SetBuffer([]byte)
	GetBuffer() []byte
	GetId() int
	Encode()
	Decode()
}

func NewPacket(id int) *Packet {
	return &Packet{id, NewStream()}
}

func (packet *Packet) GetId() int {
	return packet.packetId
}

func (packet *Packet) DecodeStep() {
	packet.offset = 1
}

func (packet *Packet) EncodeId() {
	packet.buffer = []byte{}
	var newBuffer = append(packet.buffer, byte(packet.packetId))
	packet.buffer = newBuffer
}

func (packet *Packet) ResetBase() {
	packet.ResetStream()
}
