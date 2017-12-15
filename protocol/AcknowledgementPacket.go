package protocol

import (
	"goraklib/protocol/identifiers"
	"sort"
	"goraklib/binary"
)

type AcknowledgementPacket struct {
	*Packet
	Packets []uint32
}

type ACK struct {
	*AcknowledgementPacket
}

type NAK struct {
	*AcknowledgementPacket
}

func NewACK() *ACK {
	return &ACK{&AcknowledgementPacket{NewPacket(
		identifiers.DatagramAck,
	), []uint32{}}}
}

func NewNAK() *NAK {
	return &NAK{&AcknowledgementPacket{NewPacket(
		identifiers.DatagramNack,
	), []uint32{}}}
}

func (packet *AcknowledgementPacket) Encode() {
	packet.EncodeId()

	sort.Slice(packet.Packets, func(i, j int) bool {
		return packet.Packets[i] < packet.Packets[j]
	})

	var packetCount = len(packet.Packets)

	if packetCount == 0 {
		packet.PutShort(0)
		return
	}

	var stream = binary.NewStream()
	stream.ResetStream()

	var pointer = 1
	var firstPacket = packet.Packets[0]
	var lastPacket = packet.Packets[0]

	for pointer < packetCount {
		currentPacket := packet.Packets[pointer]
		difference := currentPacket - lastPacket

		if difference == 1 {
			lastPacket = currentPacket
		}

		pointer++
	}

	if firstPacket == lastPacket {
		stream.PutByte(01)
		stream.PutLittleTriad(firstPacket)
	} else {
		stream.PutByte(00)
		stream.PutLittleTriad(firstPacket)
		stream.PutLittleTriad(lastPacket)
	}

	packet.PutShort(1)
	packet.PutBytes(stream.Buffer)
}

func (packet *AcknowledgementPacket) Decode() {
	packet.DecodeStep()

	packet.Packets = []uint32{}
	var packetCount = packet.GetShort()
	var count = 0

	for i := int16(0); i < packetCount && !packet.Feof() && count < 4096; i++ {
		if packet.GetByte() == 0 {
			var start = packet.GetLittleTriad()
			var end = packet.GetLittleTriad()

			if (end - start) > 512 {
				end = start + 512
			}

			for pack := start; pack < end; pack++ {
				packet.Packets = append(packet.Packets, pack)
				count++
			}

		} else {
			packet.Packets = append(packet.Packets, packet.GetLittleTriad())
			count++
		}
	}
}