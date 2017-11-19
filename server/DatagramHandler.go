package server

import (
	"goraklib/protocol"
	"goraklib/protocol/identifiers"
	"fmt"
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

		session.clientId = request.ClientId

		var accept = protocol.NewConnectionAccept()
		accept.ClientAddress = session.GetAddress()
		accept.ClientPort = session.GetPort()
		accept.PingSendTime = request.PingSendTime
		accept.PongSendTime = request.PingSendTime

		accept.Encode()

		var encPacket = protocol.NewEncapsulatedPacket()
		encPacket.Buffer = accept.Buffer
		encPacket.Reliability = protocol.ReliabilityUnreliable

		var datagram = protocol.NewDatagram()
		session.currentSequenceNumber++
		datagram.SequenceNumber = uint32(session.currentSequenceNumber)
		datagram.AddPacket(&encPacket)

		datagram.Encode()

		manager.SendPacket(datagram, session.GetAddress(), session.GetPort())

	case identifiers.NewIncomingConnection:
		var connection = protocol.NewNewIncomingConnection()
		connection.Buffer = packet.Buffer
		connection.Decode()

		session.SetConnected(true)

	case 0xFE:
		manager.AddProcessedEncapsulatedPacket(*packet)

	default:
		fmt.Println("Unknown encapsulated packet:", packet.Buffer[0])
	}
}
