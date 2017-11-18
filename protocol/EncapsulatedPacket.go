package protocol

import (
	"errors"
)

const (
	ReliabilityUnreliable = iota
	ReliabilityUnreliableSequenced
	ReliabilityReliable
	ReliabilityReliableOrdered
	ReliabilityReliableSequenced
	ReliabilityUnreliableWithAck
	ReliabilityReliableWithAck
	ReliabilityReliableOrderedWithAck
)

type EncapsulatedPacket struct {
	*Packet

	Reliability byte
	HasSplit bool
	Length uint
	MessageIndex uint32
	OrderIndex uint32
	OrderChannel byte
	SplitId int16
	SplitCount uint
	SplitIndex uint
}

func NewEncapsulatedPacket() EncapsulatedPacket {
	var packet = EncapsulatedPacket{NewPacket(0), 0, false, 0, 0, 0, 0, 0, 0, 0}
	return packet
}

func (packet *EncapsulatedPacket) GetFromBinary(stream *Datagram) (EncapsulatedPacket, error) {
	var flags = stream.GetByte()
	packet.Reliability = (flags & 224) >> 5
	packet.HasSplit = (flags & 16) != 0

	if stream.Feof() {
		return *packet, errors.New("no bytes left to read")
	}
	packet.Length = uint(stream.GetUnsignedShort() / 8)

	if packet.Length == 0 {
		return *packet, errors.New("null encapsulated packet")
	}

	if packet.IsReliable() {
		packet.MessageIndex = stream.GetLittleTriad()
	}

	if packet.IsSequenced() {
		packet.OrderIndex = stream.GetLittleTriad()
		packet.OrderChannel = stream.GetByte()
	}

	if packet.HasSplit {
		packet.SplitCount = uint(stream.GetInt())
		packet.SplitId = stream.GetShort()
		packet.SplitIndex = uint(stream.GetInt())
	}

	packet.SetBuffer(stream.Get(int(packet.Length)))

	return *packet, nil
}

func (packet *EncapsulatedPacket) Encode() {
	var buffer = packet.GetBuffer()

	var splitValue = 0
	if packet.HasSplit {
		splitValue = 16
	}
	packet.SetBuffer([]byte{})

	packet.PutByte(byte((packet.Reliability << 5) | byte(splitValue)))

	packet.PutShort(int16(len(buffer) << 3))

	if packet.IsReliable() {
		packet.PutLittleTriad(packet.MessageIndex)
	}
	if packet.IsSequenced() {
		packet.PutLittleTriad(packet.OrderIndex)
		packet.PutByte(packet.OrderChannel)
	}

	if packet.HasSplit {
		packet.PutInt(int32(packet.SplitCount))
		packet.PutShort(packet.SplitId)
		packet.PutInt(int32(packet.SplitIndex))
	}

	packet.PutBytes(buffer)
}

func (packet *EncapsulatedPacket) IsReliable() bool {
	switch packet.Reliability {
	case ReliabilityReliable:
		return true
	case ReliabilityReliableOrdered:
		return true
	case ReliabilityReliableSequenced:
		return true
	case ReliabilityReliableWithAck:
		return true
	case ReliabilityReliableOrderedWithAck:
		return true
	}
	return false
}

func (packet *EncapsulatedPacket) IsSequenced() bool {
	switch packet.Reliability {
	case ReliabilityUnreliableSequenced:
		return true
	case ReliabilityReliableOrdered:
		return true
	case ReliabilityReliableSequenced:
		return true
	case ReliabilityReliableOrderedWithAck:
		return true
	}
	return false
}