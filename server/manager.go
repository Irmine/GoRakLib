package server

import (
	"net"
	"time"
	"math/rand"
	"github.com/irmine/goraklib/protocol"
)

const (
	MaximumMTUSize = 1492
	MinimumMTUSize = 400
)

// Manager manages a UDP server and its components.
type Manager struct {
	Server   *UDPServer
	Sessions SessionManager

	// PongData is the data returned when the server gets an unconnected ping.
	PongData string
	// Security ensures a secure connection between a pair of systems.
	// Setting this to false is often the best idea for mobile devices.
	Security bool
	// Encryption encrypts all packets sent over RakNet.
	// Encryption should be disabled if used for Minecraft.
	Encryption bool
	// serverId is a random ID to identify servers. It is randomly generated for each manager.
	serverId int64
	// Running
	Running bool

	// RawPacketFunction gets called when a raw packet is processed.
	// The address given is the address of the sender.
	RawPacketFunction    func(packet RawPacket, addr *net.UDPAddr)
	// EncapsulatedFunction gets called once an encapsulated packet is fully processed.
	// Only encapsulated packets not recognized as RakNet internal packets get called.
	EncapsulatedFunction func(packet protocol.EncapsulatedPacket, session *Session)
}

// NewManager returns a new Manager for a UDP Server.
// A random server ID gets generated.
func NewManager() *Manager {
	rand.Seed(time.Now().Unix())
	return &Manager{Server: NewUDPServer(), Sessions: NewSessionManager(), serverId: rand.Int63()}
}

// Start starts the UDP server on the given address and port.
// Start returns an error if any might have occurred during starting.
// The manager will keep processing incoming packets until it has been Stop()ed.
func (manager *Manager) Start(address string, port int) error {
	manager.Running = true
	err := manager.Server.Start(address, port)

	go func() {
		for manager.Running {
			manager.processIncomingPacket()
		}
	}()

	return err
}

// Stop makes the manager stop processing incoming packets.
func (manager *Manager) Stop() {
	manager.Running = false
}

// processIncomingPacket processes any incoming packet from the UDP server.
func (manager *Manager) processIncomingPacket() {
	buffer := make([]byte, 8192)
	addr, err := manager.Server.Read(buffer)
	if err != nil {
		return
	}
	packet := getPacketFor(buffer, manager.Sessions.SessionExists(addr))
	if raw, ok := packet.(RawPacket); ok {
		manager.RawPacketFunction(raw, addr)
		return
	}
	if packet.HasMagic() {
		HandleUnconnectedMessage(packet, addr, manager)
	} else {
		if !manager.Sessions.SessionExists(addr) {
			return
		}
	}
}

// SessionManager is a manager of all sessions in the Manager.
type SessionManager map[*net.UDPAddr]*Session

// NewSessionManager returns a new session manager.
func NewSessionManager() SessionManager {
	return SessionManager{}
}

// SessionExists checks if the session manager has a session with a UDPAddr.
func (manager SessionManager) SessionExists(addr *net.UDPAddr) bool {
	_, ok := manager[addr]
	return ok
}