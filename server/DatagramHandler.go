package server

import (
	"goraklib/protocol"
	"goraklib/protocol/identifiers"
	"fmt"
)

func (manager *SessionManager) ProcessDatagram(datagram *protocol.Datagram, session *Session) {
	var ack = protocol.NewACK()
	ack.Packets = []uint32{datagram.SequenceNumber}
	manager.SendPacket(ack, session)

	session.receiveWindow.SubmitDatagram(datagram)
}

func (manager *SessionManager) HandleDatagram(datagram *protocol.Datagram, session *Session) {
	if datagram == nil {
		return
	}

	if datagram.SequenceNumber <= session.receiveSequence && datagram.SequenceNumber != 0 {
		return
	}

	session.receiveSequence = datagram.SequenceNumber
	for _, packet := range *datagram.GetPackets() {
		manager.HandleEncapsulated(packet, session)
	}
}

func (manager *SessionManager) HandleAck(ack *protocol.ACK, session *Session) {
	session.recoveryQueue.FlagForDeletion(ack.Packets)
}

func (manager *SessionManager) HandleNak(nak *protocol.NAK, session *Session) {
	var datagrams = session.recoveryQueue.Recover(nak.Packets)
	for _, datagram := range datagrams {
		manager.server.udp.WriteBuffer(datagram.GetBuffer(), session.GetAddress(), session.GetPort())
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

		session.SendConnectedPacket(accept, protocol.ReliabilityReliableOrdered, PriorityImmediate)

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

		session.SendConnectedPacket(pong, protocol.ReliabilityUnreliable, PriorityLow)

	case identifiers.ConnectedPong:
		var pong = protocol.NewConnectedPong()
		pong.Buffer = packet.Buffer
		pong.Decode()

		ping := uint64(manager.server.GetRunTime() - pong.PingSendTime)

		session.SetPing(ping)

	case identifiers.DisconnectNotification:
		session.Close(true)

	default:
		if !session.IsConnected() {
			fmt.Println("Unknown encapsulated packet:", packet.Buffer[0])
			return
		}
		session.AddProcessedEncapsulatedPacket(*packet)
	}
}

func (manager *SessionManager) HandleSplitEncapsulated(packet *protocol.EncapsulatedPacket, session *Session) {
	var id = int(packet.SplitId)

	if session.splits[id] == nil {
		session.splits[id] = make([]*protocol.EncapsulatedPacket, packet.SplitCount)
		session.splitCounts[id] = 0
	}

	if pk := session.splits[id][packet.SplitIndex]; pk == nil {
		session.splitCounts[id]++
	}

	session.splits[id][packet.SplitIndex] = packet

	if session.splitCounts[id] == packet.SplitCount {
		var newPacket = protocol.NewEncapsulatedPacket()
		for _, pk := range session.splits[id] {
			newPacket.PutBytes(pk.Buffer)
		}

		manager.HandleEncapsulated(newPacket, session)
		delete(session.splits, id)
	}
}

func (manager *SessionManager) SendPing(session *Session) {
	var ping = protocol.NewConnectedPing()
	ping.PingSendTime = manager.server.GetRunTime()

	session.SendConnectedPacket(ping, protocol.ReliabilityUnreliable, PriorityImmediate)
}
