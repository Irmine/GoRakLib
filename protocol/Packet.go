package protocol

import (
	"strings"
	"strconv"
	"goraklib/binary"
)

type Packet struct {
	packetId int
	*binary.BinaryStream
}

type IPacket interface {
	SetBuffer([]byte)
	GetBuffer() []byte
	GetId() int
	Encode()
	Decode()
	HasMagic() bool
}

func NewPacket(id int) *Packet {
	return &Packet{id, binary.NewStream()}
}

func (packet *Packet) GetId() int {
	return packet.packetId
}

func (packet *Packet) HasMagic() bool {
	var magicString = string(magic)
	return strings.Contains(string(packet.Buffer), magicString)
}

func (packet *Packet) DecodeStep() {
	packet.Offset = 1
}

func (packet *Packet) EncodeId() {
	packet.Buffer = []byte{}
	var newBuffer = append(packet.Buffer, byte(packet.packetId))
	packet.Buffer = newBuffer
}

func (packet *Packet) ResetBase() {
	packet.ResetStream()
}

func (packet *Packet) GetAddress() (address string, port uint16, ipVersion byte) {
	ipVersion = packet.GetByte()
	switch ipVersion {
	default:
	case 4:
		var parts = []byte{(-packet.GetByte() - 1) & 0xff, (-packet.GetByte() - 1) & 0xff, (-packet.GetByte() - 1) & 0xff, (-packet.GetByte() - 1) & 0xff}
		var stringArr = []string{}
		for _, part := range parts {
			stringArr = append(stringArr, strconv.Itoa(int(part)))
		}

		address = strings.Join(stringArr, ".")
		port = packet.GetUnsignedShort()
	case 6:
		packet.GetLittleShort()
		port = packet.GetUnsignedShort()
		packet.GetInt()
		address = string(packet.Get(16))
		packet.GetInt()
	}
	return
}

func (packet *Packet) PutAddress(address string, port uint16, ipVersion byte) {
	packet.PutByte(ipVersion)
	switch ipVersion {
	default:
	case 4:
		var stringArr = strings.Split(address, ".")
		for _, str := range stringArr {
			var digit, _ = strconv.Atoi(str)
			packet.PutByte(byte(digit))
		}
		packet.PutUnsignedShort(port)
	case 6:
		packet.PutLittleShort(23)
		packet.PutUnsignedShort(port)
		packet.PutInt(0)
		packet.PutBytes([]byte(address))
		packet.PutInt(0)
	}
}