package server

import (
	"strconv"
	"errors"
	"goraklib/protocol"
)

type SessionManager struct {
	server *GoRakLibServer
	sessions map[string]*Session

	splits map[int]map[int]*protocol.EncapsulatedPacket
}

func NewSessionManager(server *GoRakLibServer) *SessionManager {
	return &SessionManager{server, make(map[string]*Session), make(map[int]map[int]*protocol.EncapsulatedPacket)}
}

func (manager *SessionManager) GetSessions() map[string]*Session {
	return manager.sessions
}

func (manager *SessionManager) CreateSession(address string, port uint16) {
	var session = NewSession(address, port)
	manager.sessions[address + ":" + strconv.Itoa(int(port))] = session
}

func (manager *SessionManager) SessionExists(address string, port uint16) bool {
	var _, exists = manager.sessions[address + ":" + strconv.Itoa(int(port))]
	return exists
}

func (manager *SessionManager) GetSession(address string, port uint16) (*Session, error) {
	var session *Session
	if !manager.SessionExists(address, port) {
		return session, errors.New("session does not yet exist")
	}
	session = manager.sessions[address + ":" + strconv.Itoa(int(port))]
	return session, nil
}

func (manager *SessionManager) Tick() {
	for _, session := range manager.sessions {
		for !session.IsStackEmpty() {
			packet := session.FetchFromStack()
			if packet.HasMagic() {
				go manager.HandleUnconnectedMessage(packet, session)
			} else {
				if packet, ok := packet.(*protocol.Datagram); ok {
					go manager.HandleDatagram(packet, session)
				}
			}
		}
	}
}

func (manager *SessionManager) SendPacket(packet protocol.IPacket, session *Session) {
	manager.server.udp.WriteBuffer(packet.GetBuffer(), session.GetAddress(), session.GetPort())
}