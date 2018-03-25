package server

import (
	"net"
	"time"
	"math/rand"
	"github.com/irmine/goraklib/protocol"
	"fmt"
	"sync"
)

const (
	// Maximum MTU size is the maximum packet size.
	// Any MTU size above this will get limited to the maximum.
	MaximumMTUSize = 1492
	// MinimumMTUSize is the minimum packet size.
	// Any MTU size below this will get set to the minimum.
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
	// ServerId is a random ID to identify servers. It is randomly generated for each manager.
	ServerId int64
	// Running specifies the running state of the manager.
	// The manager will automatically stop working if the running state is false.
	Running bool
	// CurrentTick is the current tick of the manager. This current Tick increments for every
	// time the manager ticks.
	CurrentTick int64
	// TimeoutDuration is the duration after which a session gets timed out.
	// Timed out sessions get closed and removed immediately.
	// The default timeout duration is 6 seconds.
	TimeoutDuration time.Duration

	// RawPacketFunction gets called when a raw packet is processed.
	// The address given is the address of the sender, and the byte array the buffer of the packet.
	RawPacketFunction    func(packet []byte, addr *net.UDPAddr)
	// PacketFunction gets called once an encapsulated packet is fully processed.
	// This function only gets called for encapsulated packets not recognized as RakNet internal packets.
	// A byte array argument gets passed, which is the buffer of the buffer in the encapsulated packet.
	PacketFunction 		 func(packet []byte, session *Session)
	// ConnectFunction gets called once a session is fully connected to the server,
	// and packets of the game protocol start to get sent.
	ConnectFunction		 func(session *Session)
	// DisconnectFunction gets called with the associated session on a disconnect.
	// This disconnect may be either client initiated or server initiated.
	DisconnectFunction	 func(session *Session)

	*sync.RWMutex
	// ipBlocks is a field containing all blocked addresses.
	// Blocked addresses are ignored completely; Their packets are not processed.
	ipBlocks map[string]*net.UDPAddr
}

// NewManager returns a new Manager for a UDP Server.
// A random server ID gets generated.
func NewManager() *Manager {
	rand.Seed(time.Now().Unix())
	return &Manager{Server: NewUDPServer(), Sessions: NewSessionManager(), ServerId: rand.Int63(),
		RawPacketFunction: func(packet []byte, addr *net.UDPAddr) {},
		PacketFunction: func(packet []byte, session *Session) {},
		ConnectFunction: func(session *Session) {},
		DisconnectFunction: func(session *Session) {},
		ipBlocks: make(map[string]*net.UDPAddr),
		RWMutex: &sync.RWMutex{},
		TimeoutDuration: time.Second * 6,
	}
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
	go manager.tickSessions()

	return err
}

// Stop makes the manager stop processing incoming packets.
func (manager *Manager) Stop() {
	manager.Running = false
}

// BlockIP blocks the IP of the given UDP address,
// ignoring any further packets until the duration runs out.
func (manager *Manager) BlockIP(addr *net.UDPAddr, duration time.Duration) {
	manager.Lock()
	manager.ipBlocks[fmt.Sprint(addr.IP)] = addr
	manager.Unlock()
	time.AfterFunc(duration, func() {
		manager.UnblockIP(addr)
	})
}

// UnblockIP unblocks the IP of the address and allows packets from the address once again.
func (manager *Manager) UnblockIP(addr *net.UDPAddr) {
	manager.Lock()
	delete(manager.ipBlocks, fmt.Sprint(addr.IP))
	manager.Unlock()
}

// IsIPBlocked checks if the IP of a UDP address is blocked.
// If true, packets are not processed of the address.
func (manager *Manager) IsIPBlocked(addr *net.UDPAddr) bool {
	manager.RLock()
	_, ok := manager.ipBlocks[fmt.Sprint(addr.IP)]
	manager.RUnlock()
	return ok
}

// tickSessions makes the server start ticking its sessions.
// Sessions get ticked on an interval of 80 ticks per second.
func (manager *Manager) tickSessions() {
	ticker := time.NewTicker(time.Duration(time.Second / 80))
	for range ticker.C {
		if !manager.Running {
			return
		}
		for index, session := range manager.Sessions {
			manager.updateSession(session, index)
		}
		manager.CurrentTick++
	}
}

// updateSession updates a session with session index.
// The session will be ticked while it's open.
// Sessions that have not responded for too long are timed out and
// flagged for closing, and sessions flagged for closing will be cleaned up.
func (manager *Manager) updateSession(session *Session, index string) {
	session.Tick(manager.CurrentTick)
	if time.Now().Sub(session.LastUpdate) > manager.TimeoutDuration {
		session.FlagForClose()
	}
	if session.FlaggedForClose {
		session.Close()
		delete(manager.Sessions, index)
	}
}

// processIncomingPacket processes any incoming packet from the UDP server.
// Unconnected messages get handled freely, while any other packet gets passed to its owner session.
func (manager *Manager) processIncomingPacket() {
	buffer := make([]byte, 2048)
	n, addr, err := manager.Server.Read(buffer)
	if manager.IsIPBlocked(addr) {
		return
	}
	buffer = buffer[:n]

	if err != nil {
		return
	}
	manager.Sessions.SessionExists(addr)

	defer func() {
		if err := recover(); err != nil {
			manager.BlockIP(addr, time.Second * 5)
			fmt.Println("IP blocked of", addr, "for 5 seconds:", err)
		}
	}()
	packet := getPacketFor(buffer, manager.Sessions.SessionExists(addr))

	if raw, ok := packet.(RawPacket); ok {
		manager.RawPacketFunction(raw.Buffer, addr)
		return
	}
	if packet.HasMagic() {
		HandleUnconnectedMessage(packet, addr, manager)
	} else {
		session, ok := manager.Sessions.GetSession(addr)
		if !ok {
			return
		}
		if datagram, ok := packet.(*protocol.Datagram); ok {
			session.ReceiveWindow.AddDatagram(datagram)
		} else if ack, ok := packet.(*protocol.ACK); ok {
			session.HandleACK(ack)
		} else if nack, ok := packet.(*protocol.NAK); ok {
			session.HandleNACK(nack)
		}
	}
}

// SessionManager is a manager of all sessions in the Manager.
type SessionManager map[string]*Session

// NewSessionManager returns a new session manager.
func NewSessionManager() SessionManager {
	return SessionManager{}
}

// SessionExists checks if the session manager has a session with a UDPAddr.
func (manager SessionManager) SessionExists(addr *net.UDPAddr) bool {
	_, ok := manager[fmt.Sprint(addr)]
	return ok
}

// GetSession returns a session by a UDP address.
// GetSession also returns a bool indicating success of the call.
func (manager SessionManager) GetSession(addr *net.UDPAddr) (*Session, bool) {
	session, ok := manager[fmt.Sprint(addr)]
	return session, ok
}