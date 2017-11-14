package server

import (
	"math/rand"
)

type GoRakLibServer struct {
	serverName string
	serverId int64
	port int

	udp *UDPServer
	sessionManager *SessionManager
}

func NewGoRakLibServer(serverName string, port int) *GoRakLibServer {
	var server = GoRakLibServer{}
	server.serverName = serverName
	var udp = NewUDPServer(port)

	server.serverId = int64(rand.Int())

	server.udp = &udp
	server.port = port
	server.sessionManager = NewSessionManager(&server)

	go func() {
		for {
			var packet, ip, port, err = udp.ReadBuffer()
			if err != nil {
				continue
			}
			if !server.sessionManager.SessionExists(ip, port) {
				server.sessionManager.CreateSession(ip, port)
			}
			var session, _ = server.sessionManager.GetSession(ip, port)
			go session.Forward(packet)
		}
	}()

	return &server
}

func (server *GoRakLibServer) GetName() string {
	return server.serverName
}

func (server *GoRakLibServer) GetServerId() int64 {
	return server.serverId
}

func (server *GoRakLibServer) GetUDP() *UDPServer {
	return server.udp
}

func (server *GoRakLibServer) GetPort() int {
	return server.port
}

func (server *GoRakLibServer) Tick() {
	server.sessionManager.Tick()
}