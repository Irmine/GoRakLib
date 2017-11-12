package protocol

type Packet struct {
	*BinaryStream
	packetId int
}

func (packet *Packet) DecodeStep() {
	packet.offset = 1
}

func (packet *Packet) EncodeId() {
	packet.buffer = append(packet.buffer, byte(packet.packetId))
}
