package server

import (
	"github.com/irmine/goraklib/protocol"
	"net"
	"time"
	"fmt"
)

// HandleUnconnectedMessage handles an incoming unconnected message from a UDPAddr.
// A response will be made for every packet, which gets sent back to the sender.
// A session gets created for the sender once the OpenConnectionRequest2 gets sent.
func HandleUnconnectedMessage(packetInterface protocol.IPacket, addr *net.UDPAddr, manager *Manager) {
	switch packet := packetInterface.(type) {
	case *protocol.UnconnectedPing:
		handleUnconnectedPing(addr, manager)
	case *protocol.OpenConnectionRequest1:
		handleOpenConnectionRequest1(packet, addr, manager)
	case *protocol.OpenConnectionRequest2:
		handleOpenConnectionRequest2(packet, addr, manager)
	}
}

// handleUnconnectedPing handles an unconnected ping.
// An unconnected pong is sent back with the server's pong data.
func handleUnconnectedPing(addr *net.UDPAddr, manager *Manager) {
	pong := protocol.NewUnconnectedPong()
	pong.PingTime = time.Now().Unix()
	pong.ServerId = manager.ServerId
	pong.PongData = manager.PongData
	pong.Encode()
	manager.Server.Write(pong.Buffer, addr)
}

// handleOpenConnectionRequest1 handles an open connection request 1.
// An open connection response 1 is sent back with the MTU size and security.
func handleOpenConnectionRequest1(request *protocol.OpenConnectionRequest1, addr *net.UDPAddr, manager *Manager) {
	reply := protocol.NewOpenConnectionReply1()
	reply.ServerId = manager.ServerId
	reply.MtuSize = request.MtuSize
	reply.Security = manager.Security
	reply.Encode()
	manager.Server.Write(reply.Buffer, addr)
}

// handleOpenConnectionRequest2 handles an open connection request 2.
// An open connection response 2 is sent back, with the definite MTU size and encryption.
func handleOpenConnectionRequest2(request *protocol.OpenConnectionRequest2, addr *net.UDPAddr, manager *Manager) {
	reply := protocol.NewOpenConnectionReply2()
	reply.ServerId = manager.ServerId
	if request.MtuSize < MinimumMTUSize {
		request.MtuSize = MinimumMTUSize
	} else if request.MtuSize > MaximumMTUSize {
		request.MtuSize = MaximumMTUSize
	}
	reply.MtuSize = request.MtuSize
	reply.UseEncryption = manager.Encryption
	reply.ClientAddress = addr.IP.String()
	reply.ClientPort = uint16(addr.Port)

	reply.Encode()

	manager.Sessions[fmt.Sprint(addr)] = NewSession(addr, request.MtuSize, manager)
	manager.Server.Write(reply.Buffer, addr)
}