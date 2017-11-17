package protocol

const (
	BitFlagPacketPair = 0x10
	BitFlagContinuousSend = 0x08
	BitFlagNeedsBAndAs = 0x04
)

type Datagram struct {
	*Packet

	PacketPair bool
	ContinuousSend bool
	NeedsBAndAs bool

	SequenceNumber uint32

	packets *[]*EncapsulatedPacket
}

func NewDatagram() *Datagram {
	return &Datagram{NewPacket(0), false, false, false, 0, &[]*EncapsulatedPacket{}}
}

func (datagram *Datagram) GetPackets() *[]*EncapsulatedPacket {
	return datagram.packets
}

func (datagram *Datagram) Encode() {

}

func (datagram *Datagram) Decode() {
	var flags = datagram.GetByte()
	datagram.PacketPair = (flags & BitFlagPacketPair) != 0
	datagram.ContinuousSend = (flags & BitFlagContinuousSend) != 0
	datagram.NeedsBAndAs = (flags & BitFlagNeedsBAndAs) != 0

	datagram.SequenceNumber = datagram.GetLittleTriad()

	for !datagram.Feof() {
		packet, err := NewEncapsulatedPacket(datagram)
		if err == nil {
			var packets = append(*datagram.packets, &packet)
			datagram.packets = &packets
		}
	}
}