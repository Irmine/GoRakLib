package client

import (
	"goraklib/protocol"
	"time"
)

func (client *GoRakLibClient) HandleUnconnectedMessage(packetInterface protocol.IPacket) {
	switch packet := packetInterface.(type) {
	case *protocol.UnconnectedPong:
		client.isConnected = true
		client.isDiscoveringMtuSize = true

		println(packet.ServerData)

	case *protocol.OpenConnectionResponse1:
		client.isDiscoveringMtuSize = false

		var request = protocol.NewOpenConnectionRequest2()
		request.ServerAddress = client.connectionAddress
		request.ServerPort = client.connectionPort
		request.MtuSize = packet.MtuSize
		request.ClientId = 0

		println(packet.MtuSize)

		client.SendPacket(request)

	case *protocol.OpenConnectionResponse2:
		client.mtuSize = packet.MtuSize
		// var encryption = packet.UseEncryption // Encryption...

		var request = protocol.NewConnectionRequest()
		request.ClientId = 0
		request.PingSendTime = uint64(time.Now().Second()) * 1000

		client.SendConnectedPacket(request, protocol.ReliabilityUnreliable, PriorityImmediate)
	}
}
