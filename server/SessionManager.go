package server

import (
	"strconv"
	"errors"
	"goraklib/protocol"
	"time"
)

type SessionManager struct {
	server *GoRakLibServer
	sessions map[string]*Session
	disconnectedSessions chan *Session
}

func NewSessionManager(server *GoRakLibServer) *SessionManager {
	return &SessionManager{server, make(map[string]*Session), make(chan *Session, 512)}
}

func GetSessionIndex(session *Session) string {
	return session.String()
}

func (manager *SessionManager) GetSessions() map[string]*Session {
	return manager.sessions
}

func (manager *SessionManager) CreateSession(address string, port uint16) {
	var session = NewSession(manager, address, port)
	manager.sessions[GetSessionIndex(session)] = session
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
		go func(session *Session) {
			if (session.LastUpdate + 10) < time.Now().Unix() {
				session.Close(false)
			}

			if session.IsReadyForDeletion() {
				manager.Disconnect(session)
				return
			}

			if session.IsClosed() {
				return
			}
			session.queue.Flush()
		}(session)
	}
}

func (manager *SessionManager) HandlePacket(packet protocol.IPacket, session *Session) {
	if packet.HasMagic() {
		manager.HandleUnconnectedMessage(packet, session)
	} else if session.IsOpened() {
		if datagram, ok := packet.(*protocol.Datagram); ok {
			manager.HandleDatagram(datagram, session)
		} else if nak, ok := packet.(*protocol.NAK); ok {
			manager.HandleNak(nak, session)
		} else if ack, ok := packet.(*protocol.ACK); ok {
			manager.HandleAck(ack, session)
		}
	}
}

func (manager *SessionManager) SendPacket(packet protocol.IPacket, session *Session) {
	session.LastUpdate = time.Now().Unix()
	packet.Encode()

	manager.server.udp.WriteBuffer(packet.GetBuffer(), session.GetAddress(), session.GetPort())
}

func (manager *SessionManager) Disconnect(session *Session) {
	if !session.IsReadyForDeletion() {

	}

	delete(manager.sessions, GetSessionIndex(session))

	if session.IsConnected() {
		manager.disconnectedSessions <- session
	}
}

func (manager *SessionManager) GetDisconnectedSessions() map[string]*Session {
	var sessions = map[string]*Session{}
	for len(manager.disconnectedSessions) > 0 {
		session := <-manager.disconnectedSessions
		sessions[GetSessionIndex(session)] = session
	}
	return sessions
}