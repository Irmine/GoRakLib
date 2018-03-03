package protocol

import (
	"bytes"
	"github.com/irmine/binutils"
)

var magic = []byte{0x00, 0xff, 0xff, 0x00, 0xfe, 0xfe, 0xfe, 0xfe, 0xfd, 0xfd, 0xfd, 0xfd, 0x12, 0x34, 0x56, 0x78}

type UnconnectedMessage struct {
	*Packet
	magic []byte
}

func NewUnconnectedMessage(packet *Packet) *UnconnectedMessage {
	return &UnconnectedMessage{packet, make([]byte, 16)}
}

func (message *UnconnectedMessage) PutMagic() {
	message.PutBytes(magic)
}

func (message *UnconnectedMessage) ReadMagic() {
	message.magic = binutils.Read(&message.Buffer, &message.Offset, 16)
}

func (message *UnconnectedMessage) HasValidMagic() bool {
	return bytes.Equal(message.magic, magic)
}
