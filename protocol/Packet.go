package protocol

import (
	"strings"
	"strconv"
)

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
		for _, string := range stringArr {
			var str, _ = strconv.Atoi(string)
			packet.PutByte(byte(str))
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
