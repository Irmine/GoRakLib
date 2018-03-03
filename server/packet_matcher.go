package server

import "github.com/irmine/goraklib/protocol"

// GetPacketFor selects the appropriate packet by a buffer.
// It uses hasSession to check for appropriate messages.
func getPacketFor(buffer []byte, hasSession bool) protocol.IPacket {
	header := buffer[0]
	var packet protocol.IPacket
	if hasSession {
		if header & protocol.BitFlagIsAck != 0 {
			packet = protocol.NewACK()
		} else if header & protocol.BitFlagIsNak != 0 {
			packet = protocol.NewNAK()
		} else if header & protocol.BitFlagValid != 0 {
			packet = protocol.NewDatagram()
		}
	} else {
		switch header {
		case protocol.IdUnconnectedPing:
			packet = protocol.NewUnconnectedPing()
		case protocol.IdOpenConnectionRequest1:
			packet = protocol.NewOpenConnectionRequest1()
		case protocol.IdOpenConnectionRequest2:
			packet = protocol.NewOpenConnectionRequest2()
		}
	}
	if packet == nil {
		packet = NewRawPacket()
	}
	packet.SetBuffer(buffer)
	return packet
}
