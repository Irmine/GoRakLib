package server

import (
	"goraklib/protocol"
	"goraklib/protocol/identifiers"
)

func (manager *SessionManager) HandleDatagram(datagram *protocol.Datagram, session *Session) {
	for _, packet := range *datagram.GetPackets() {
		manager.HandleEncapsulated(packet, session)
	}
}

func (manager *SessionManager) HandleEncapsulated(packet *protocol.EncapsulatedPacket, session *Session) {
	switch packet.Buffer[0] {
	case identifiers.ConnectionRequest:
		var request = protocol.NewConnectionRequest()
		request.Buffer = packet.GetBuffer()
		request.Decode()
	}
}
