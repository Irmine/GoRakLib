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

func (manager *SessionManager) GetSessions() map[string]*Session {
	return manager.sessions
}

func (manager *SessionManager) CreateSession(address string, port uint16) {
	var session = NewSession(manager, address, port)
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
	for index, session := range manager.sessions {
		go func(session *Session) {
			if session.IsReadyForDeletion() {
				delete(manager.sessions, index)
				return
			}

			for !session.IsStackEmpty() {
				packet := session.FetchFromStack()
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
			session.queue.Flush()
		}(session)
	}
}

func (manager *SessionManager) SendPacket(packet protocol.IPacket, session *Session) {
	packet.Encode()

	manager.server.udp.WriteBuffer(packet.GetBuffer(), session.GetAddress(), session.GetPort())
}