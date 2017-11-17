package protocol

import (
	"errors"
)

const (
	ReliabilityUnreliable = 0
	ReliabilityUnreliableSequenced = 1
	ReliabilityReliable = 2
	ReliabilityReliableOrdered = 3
	ReliabilityReliableSequenced = 4
	ReliabilityUnreliableWithAck = 5
	ReliabilityReliableWithAck = 6
	ReliabilityReliableOrderedWithAck = 7
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

	IdentifierAck int
	NeedsAck bool
}

func NewEncapsulatedPacket(stream *Datagram) (EncapsulatedPacket, error) {
	var packet = EncapsulatedPacket{NewPacket(0), 0, false, 0, 0, 0, 0, 0, 0, 0, 0, false}

	var flags = stream.GetByte()
	packet.Reliability = (flags & 224) >> 5
	packet.HasSplit = (flags & 16) != 0

	if stream.Feof() {
		return EncapsulatedPacket{}, errors.New("no bytes left to read")
	}
	packet.Length = uint(stream.GetUnsignedShort() / 8)

	if packet.Length == 0 {
		return EncapsulatedPacket{}, errors.New("null encapsulated packet")
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

	return packet, nil
}

func (packet *EncapsulatedPacket) Decode() {

}

func (packet *EncapsulatedPacket) Encode() {

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