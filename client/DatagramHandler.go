package client

import (
	"goraklib/protocol"
	"goraklib/protocol/identifiers"
	"fmt"
	"encoding/hex"
)

const (
	MinecraftHeader = 0xFE
)

func (client *GoRakLibClient) HandleDatagram(datagram *protocol.Datagram) {
	var ack = protocol.NewACK()
	ack.Packets = []uint32{datagram.SequenceNumber}
	client.SendPacket(ack)

	for _, packet := range *datagram.GetPackets() {
		client.HandleEncapsulated(packet)
	}
}

func (client *GoRakLibClient) HandleAck(ack *protocol.ACK) {
	fmt.Println(hex.EncodeToString(ack.GetBuffer()))
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
		client.HandleSplitEncapsulated(packet)
		return
	}

	switch packet.Buffer[0] {
	case identifiers.ConnectionAccept:
		println(hex.EncodeToString(packet.Buffer))

		var incoming = protocol.NewNewIncomingConnection()
		incoming.ServerPort = client.connectionPort
		incoming.ServerAddress = client.connectionAddress

		for i := 0; i < 20; i++ {
			incoming.SystemAddresses = append(incoming.SystemAddresses, "0.0.0.0")
			incoming.SystemPorts = append(incoming.SystemPorts, 0)
			incoming.SystemIdVersions = append(incoming.SystemIdVersions, 4)
		}

		client.online = true
		client.SendConnectedPacket(incoming, protocol.ReliabilityUnreliable, PriorityHigh)

	case identifiers.ConnectedPing:
		var ping = protocol.NewConnectedPing()
		ping.Buffer = packet.Buffer
		ping.Decode()

		var pong = protocol.NewConnectedPong()
		pong.PingSendTime = ping.PingSendTime
		pong.PongSendTime = 0

		client.SendConnectedPacket(pong, protocol.ReliabilityUnreliable, PriorityLow)

		client.SendPing()

	case identifiers.ConnectedPong:
		var pong = protocol.NewConnectedPong()
		pong.Buffer = packet.Buffer
		pong.Decode()

		//ping := uint64(client.server.GetRunTime() - pong.PingSendTime)

		//client.SetPing(ping)

	case identifiers.DisconnectNotification:
		client.Disconnect()

	case MinecraftHeader:
		/*if !client.IsConnected() {
			return
		}
		client.AddProcessedEncapsulatedPacket(*packet)*/

	default:
		fmt.Println("Unknown encapsulated packet:", packet.Buffer[0])
		//client.AddProcessedEncapsulatedPacket(*packet)
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

	//client.SendConnectedPacket(ping, protocol.ReliabilityUnreliable, PriorityMedium)
}

