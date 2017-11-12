package protocol

import "goraklib/protocol/identifiers"

type AcknowledgementPacket struct {
	*Packet
	packets []uint32
}

type ACK struct {
	*AcknowledgementPacket
}

type NACK struct {
	*AcknowledgementPacket
}

func NewACK() ACK {
	return ACK{&AcknowledgementPacket{&Packet{
		packetId: identifiers.PacketAck,
	}, []uint32{}}}
}

func NewNACK() NACK {
	return NACK{&AcknowledgementPacket{&Packet{
		packetId: identifiers.PacketNack,
	}, []uint32{}}}
}

func (packet AcknowledgementPacket) Encode() {

}

func (packet AcknowledgementPacket) Decode() {
	packet.DecodeStep()
	packet.packets = []uint32{}
	var packetCount = packet.ReadShort()
	var count = 0

	for i := int16(0); i < packetCount && !packet.Feof() && count < 4096; i++ {
		if packet.ReadByte() == 0 {
			var start = packet.ReadLittleEndianTriad()
			var end = packet.ReadLittleEndianTriad()

			if (end - start) > 512 {
				end = start + 512
			}

			for pack := start; pack < end; pack++ {
				packet.packets = append(packet.packets, pack)
				count++
			}

		} else {
			packet.packets = append(packet.packets, packet.ReadLittleEndianTriad())
			count++
		}
	}
}

func (packet AcknowledgementPacket) Reset() {
	packet.ResetBase()
	packet.packets = []uint32{}
}