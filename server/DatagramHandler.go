package server

import (
	"goraklib/protocol"
	"goraklib/protocol/identifiers"
	"fmt"
)

const (
	MinecraftHeader = 0xFE
)

func (manager *SessionManager) HandleDatagram(datagram *protocol.Datagram, session *Session) {
	var ack = protocol.NewACK()
	ack.Packets = []uint32{datagram.SequenceNumber}
	manager.SendPacket(ack, session)

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

		manager.SendPing(session)

	case identifiers.ConnectedPong:
		var pong = protocol.NewConnectedPong()
		pong.Buffer = packet.Buffer
		pong.Decode()

		ping := uint64(manager.server.GetRunTime() - pong.PingSendTime)

		session.SetPing(ping)

	case identifiers.DisconnectNotification:
		session.Close()

	case MinecraftHeader:
		if !session.IsConnected() {
			return
		}
		session.AddProcessedEncapsulatedPacket(*packet)

	default:
		fmt.Println("Unknown encapsulated packet:", packet.Buffer[0])
		session.AddProcessedEncapsulatedPacket(*packet)
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
		newPacket.ResetStream()

		var packets = make([]*protocol.EncapsulatedPacket, packet.SplitCount)

		for len(session.splits[id]) != 0 {
			pk := <-session.splits[id]
			packets[pk.SplitIndex] = pk
		}

		for _, pk := range packets {
			newPacket.PutBytes(pk.Buffer)
		}

		session.lastSplitSize = packet.SplitCount - 1

		manager.HandleEncapsulated(newPacket, session)

		delete(session.splits, id)
	}
}

func (manager *SessionManager) SendPing(session *Session) {
	var ping = protocol.NewConnectedPing()
	ping.PingSendTime = manager.server.GetRunTime()

	session.SendConnectedPacket(ping, protocol.ReliabilityUnreliable, PriorityMedium)
}
