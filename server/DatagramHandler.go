package server

import (
	"goraklib/protocol"
	"goraklib/protocol/identifiers"
	"fmt"
)

func (manager *SessionManager) HandleDatagram(datagram *protocol.Datagram, session *Session) {

	var ack = protocol.NewACK()
	ack.Packets = []uint32{datagram.SequenceNumber}
	ack.Encode()
	manager.SendPacket(ack, session)

	for _, packet := range *datagram.GetPackets() {
		manager.HandleEncapsulated(packet, session)
	}
}

func (manager *SessionManager) HandleEncapsulated(packet *protocol.EncapsulatedPacket, session *Session) {
	if packet.HasSplit {
		manager.HandleSplitEncapsulated(packet, session)
		return
	}

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
		var pongTime = uint64(manager.server.GetRunTime())
		accept.PongSendTime = pongTime

		manager.sendRawPacket(accept, session)

		session.SetPing(pongTime - request.PingSendTime)

	case identifiers.NewIncomingConnection:
		var connection = protocol.NewNewIncomingConnection()
		connection.Buffer = packet.Buffer
		connection.Decode()

		session.SetConnected(true)

	case identifiers.ConnectedPing:
		var ping = protocol.NewConnectedPing()
		ping.Buffer = packet.Buffer
		ping.Decode()

		var pong = protocol.NewConnectedPong()
		pong.PingSendTime = ping.PingSendTime
		var pongTime = manager.server.GetRunTime()
		pong.PongSendTime = pongTime

		manager.sendRawPacket(pong, session)

		manager.SendPing(session)

	case identifiers.ConnectedPong:
		var pong = protocol.NewConnectedPong()
		pong.Buffer = packet.Buffer
		pong.Decode()

		ping := uint64(manager.server.GetRunTime() - pong.PingSendTime)

		session.SetPing(ping)

	case 0xFE:
		session.AddProcessedEncapsulatedPacket(*packet)

	default:
		fmt.Println("Unknown encapsulated packet:", packet.Buffer[0])
	}
}

func (manager *SessionManager) HandleSplitEncapsulated(packet *protocol.EncapsulatedPacket, session *Session) {
	var id = int(packet.SplitId)

	if session.splits[id] == nil {
		session.splits[id] = make(chan *protocol.EncapsulatedPacket, packet.SplitCount)
	}

	session.splits[id] <- packet

	if len(session.splits[id]) == int(packet.SplitCount) {
		var newPacket = protocol.NewEncapsulatedPacket()

		for len(session.splits[id]) != 0 {
			pk := <-session.splits[id]
			newPacket.PutBytes(pk.GetBuffer())
		}

		manager.HandleEncapsulated(&newPacket, session)

		delete(session.splits, id)
	}
}

func (manager *SessionManager) SendPing(session *Session) {
	var ping = protocol.NewConnectedPing()
	ping.PingSendTime = manager.server.GetRunTime()

	manager.sendRawPacket(ping, session)
}

func (manager *SessionManager) sendRawPacket(packet protocol.IPacket, session *Session) {
	var encPacket = protocol.NewEncapsulatedPacket()
	packet.Encode()
	encPacket.Buffer = packet.GetBuffer()
	encPacket.Reliability = protocol.ReliabilityUnreliable

	var datagram = protocol.NewDatagram()
	session.currentSequenceNumber++
	datagram.SequenceNumber = uint32(session.currentSequenceNumber)
	datagram.AddPacket(&encPacket)

	datagram.Encode()

	manager.SendPacket(datagram, session)
}