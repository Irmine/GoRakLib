package protocol

import "bytes"

var magic = []byte{0, 255, 255, 0, 254, 254, 254, 254, 253, 253, 253, 253, 18, 52, 86, 120}

type UnconnectedMessage struct {
	*Packet
	magic []byte
}

func NewUnconnectedMessage(packet *Packet) *UnconnectedMessage {
	return &UnconnectedMessage{packet, make([]byte, 16)}
}

func (message *UnconnectedMessage) WriteMagic() {
	message.PutBytes(magic)
}

func (message *UnconnectedMessage) ReadMagic() {
	message.magic = Read(&message.buffer, &message.offset, 16)
}

func (message *UnconnectedMessage) HasValidMagic() bool {
	return bytes.Equal(message.magic, magic)
}