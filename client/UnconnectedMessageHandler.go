package client

import "goraklib/protocol"

func (client *GoRakLibClient) HandleUnconnectedMessage(packetInterface protocol.IPacket) {
	switch packet := packetInterface.(type) {
	case *protocol.UnconnectedPong:
		var request = protocol.NewOpenConnectionRequest1()

		client.SendPacket(request)

	case *protocol.OpenConnectionResponse1:
		var request = protocol.NewOpenConnectionResponse2()

		client.SendPacket(request)

	case *protocol.OpenConnectionResponse2:

	}
}
