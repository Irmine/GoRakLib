package server

import (
	"strconv"
	"errors"
	"goraklib/protocol"
)

type SessionManager struct {
	server *GoRakLibServer
	sessions map[string]*Session

	packetBatches chan protocol.EncapsulatedPacket
	splits map[int]map[int]*protocol.EncapsulatedPacket
}

func NewSessionManager(server *GoRakLibServer) *SessionManager {
	return &SessionManager{server, make(map[string]*Session), make(chan protocol.EncapsulatedPacket, 512), make(map[int]map[int]*protocol.EncapsulatedPacket)}
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

func (manager *SessionManager) AddProcessedEncapsulatedPacket(packet protocol.EncapsulatedPacket) {
	manager.packetBatches <- packet
}

func (manager *SessionManager) GetReadyEncapsulatedPackets() []protocol.EncapsulatedPacket {
	var packets = []protocol.EncapsulatedPacket{}
	for len(manager.packetBatches) != 0 {
		packets = append(packets, <-manager.packetBatches)
	}
	return packets
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

func (manager *SessionManager) SendPacket(packet protocol.IPacket, ip string, port uint16) {
	manager.server.udp.WriteBuffer(packet.GetBuffer(), ip, port)
}