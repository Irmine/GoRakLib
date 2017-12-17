package client

import (
	"goraklib/protocol"
	"goraklib/protocol/identifiers"
	"fmt"
)

const (
	MinecraftHeader = 0xFE
)

func (client *GoRakLibClient) HandleDatagram(datagram *protocol.Datagram) {
	var ack = protocol.NewACK()
	ack.Packets = []uint32{datagram.SequenceNumber}
	client.SendPacket(ack)

	for _, packet := range *datagram.GetPackets() {
		client.HandleEncapsulated(packet, client)
	}
}

func (client *GoRakLibClient) HandleAck(ack *protocol.ACK) {
	client.recoveryQueue.FlagForDeletion(ack.Packets)
}

func (client *GoRakLibClient) HandleNak(nak *protocol.NAK) {
	var datagrams = client.recoveryQueue.Recover(nak.Packets)
	for _, datagram := range datagrams {
		client.WriteToUDP(datagram.GetBuffer())
	}
}

func (client *GoRakLibClient) HandleEncapsulated(packet *protocol.EncapsulatedPacket) {
	if packet.HasSplit {
		client.HandleSplitEncapsulated(packet, client)
		return
	}

	switch packet.Buffer[0] {
	case identifiers.ConnectionRequest:
		var request = protocol.NewConnectionRequest()
		request.Buffer = packet.GetBuffer()
		request.Decode()

		client.clientId = request.ClientId

		var accept = protocol.NewConnectionAccept()
		accept.ClientAddress = client.GetAddress()
		accept.ClientPort = client.GetPort()

		accept.PingSendTime = request.PingSendTime
		var pongTime = uint64(client.server.GetRunTime())
		accept.PongSendTime = pongTime

		client.SendConnectedPacket(accept, protocol.ReliabilityReliableOrdered, PriorityImmediate)

		client.SetPing(pongTime - request.PingSendTime)

	case identifiers.NewIncomingConnection:
		var connection = protocol.NewNewIncomingConnection()
		connection.Buffer = packet.Buffer
		connection.Decode()

		client.SetConnected(true)

	case identifiers.ConnectedPing:
		var ping = protocol.NewConnectedPing()
		ping.Buffer = packet.Buffer
		ping.Decode()

		var pong = protocol.NewConnectedPong()
		pong.PingSendTime = ping.PingSendTime
		var pongTime = client.server.GetRunTime()
		pong.PongSendTime = pongTime

		client.SendConnectedPacket(pong, protocol.ReliabilityUnreliable, PriorityLow)

		client.SendPing()

	case identifiers.ConnectedPong:
		var pong = protocol.NewConnectedPong()
		pong.Buffer = packet.Buffer
		pong.Decode()

		ping := uint64(client.server.GetRunTime() - pong.PingSendTime)

		client.SetPing(ping)

	case identifiers.DisconnectNotification:
		client.Disconnect()

	case MinecraftHeader:
		if !client.IsConnected() {
			return
		}
		client.AddProcessedEncapsulatedPacket(*packet)

	default:
		fmt.Println("Unknown encapsulated packet:", packet.Buffer[0])
		client.AddProcessedEncapsulatedPacket(*packet)
	}
}

func (client *GoRakLibClient) HandleSplitEncapsulated(packet *protocol.EncapsulatedPacket) {
	var id = int(packet.SplitId)

	if client.splits[id] == nil {
		client.splits[id] = make(chan *protocol.EncapsulatedPacket, packet.SplitCount)
	}

	client.splits[id] <- packet

	if len(client.splits[id]) == int(packet.SplitCount) {
		var newPacket = protocol.NewEncapsulatedPacket()
		newPacket.ResetStream()

		var packets = make([]*protocol.EncapsulatedPacket, packet.SplitCount)

		for len(client.splits[id]) != 0 {
			pk := <-client.splits[id]
			packets[pk.SplitIndex] = pk
		}

		for _, pk := range packets {
			newPacket.PutBytes(pk.Buffer)
		}

		client.HandleEncapsulated(newPacket)

		delete(client.splits, id)
	}
}

func (client *GoRakLibClient) SendPing() {
	var ping = protocol.NewConnectedPing()
	ping.PingSendTime = 0

	client.SendConnectedPacket(ping, protocol.ReliabilityUnreliable, PriorityMedium)
}

