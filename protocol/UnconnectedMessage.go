package protocol

import (
	"bytes"
	"goraklib/binary"
)

var magic = []byte{0x00, 0xff, 0xff, 0x00, 0xfe, 0xfe, 0xfe, 0xfe, 0xfd, 0xfd, 0xfd, 0xfd, 0x12, 0x34, 0x56, 0x78}

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
	message.magic = binary.Read(&message.Buffer, &message.Offset, 16)
}

func (message *UnconnectedMessage) HasValidMagic() bool {
	return bytes.Equal(message.magic, magic)
}