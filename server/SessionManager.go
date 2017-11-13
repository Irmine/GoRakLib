package server

import (
	"strconv"
	"errors"
	"goraklib/protocol"
)

type SessionManager struct {
	server *GoRakLibServer
	sessions map[string]*Session
}

func NewSessionManager(server *GoRakLibServer) *SessionManager {
	return &SessionManager{server, make(map[string]*Session)}
}

func (manager *SessionManager) CreateSession(address string, port int) {
	var session = NewSession(address, port)
	manager.sessions[address + ":" + strconv.Itoa(port)] = session
}

func (manager *SessionManager) SessionExists(address string, port int) bool {
	var _, exists = manager.sessions[address + ":" + strconv.Itoa(port)]
	return exists
}

func (manager *SessionManager) GetSession(address string, port int) (*Session, error) {
	var session *Session
	if !manager.SessionExists(address, port) {
		return session, errors.New("session does not yet exist")
	}
	session = manager.sessions[address + ":" + strconv.Itoa(port)]
	return session, nil
}

func (manager *SessionManager) Tick() {
	for _, session := range manager.sessions {
		for !session.IsStackEmpty() {
			manager.HandlePacket(session.FetchFromStack(), session)
		}
	}
}

func (manager *SessionManager) HandlePacket(packetInterface protocol.IPacket, session *Session) {
	switch packet := packetInterface.(type) {
	case *protocol.UnconnectedPing:
		var unconnectedPong = protocol.NewUnconnectedPong()

		unconnectedPong.PingId = packet.PingId
		unconnectedPong.ServerId = manager.server.GetServerId()
		unconnectedPong.ServerName = manager.server.GetName()

		unconnectedPong.Encode()
		manager.SendPacket(unconnectedPong, session.address, session.port)
	}
}

func (manager *SessionManager) SendPacket(packet protocol.IPacket, ip string, port int) {
	manager.server.udp.WriteBuffer(packet.GetBuffer(), ip, port)
}