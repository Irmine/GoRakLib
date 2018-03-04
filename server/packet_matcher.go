package server

import "github.com/irmine/goraklib/protocol"

// GetPacketFor selects the appropriate packet by a buffer.
// It uses hasSession to check for appropriate messages.
func getPacketFor(buffer []byte, hasSession bool) protocol.IPacket {
	header := buffer[0]
	var packet protocol.IPacket
	if hasSession {
		switch {
		case header & protocol.BitFlagIsAck != 0:
			packet = protocol.NewACK()
		case header & protocol.BitFlagIsNak != 0:
			packet = protocol.NewNAK()
		case header & protocol.BitFlagValid != 0:
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
	packet.Decode()
	return packet
}
