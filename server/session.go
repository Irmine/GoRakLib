package server

import (
	"net"
	"fmt"
	"github.com/irmine/goraklib/protocol"
	"sync"
)

// Session is a manager of a connection between the client and the server.
// Sessions manage everything related to packet ordering and processing.
type Session struct {
	*net.UDPAddr
	Manager       *Manager
	ReceiveWindow *ReceiveWindow
	RecoveryQueue *RecoveryQueue

	MTUSize		  int16
	Indexes		  Indexes
}

// Indexes is used for the collection of indexes related to datagrams and encapsulated packets.
// It uses several maps and is therefore protected by a mutex.
type Indexes struct {
	sync.Mutex
	splits		map[int16][]*protocol.EncapsulatedPacket
	splitCounts map[int16]uint
}

// NewSession returns a new session with UDP address.
// The MTUSize provided is the maximum packet size of the session.
func NewSession(addr *net.UDPAddr, mtuSize int16, manager *Manager) *Session {
	session := &Session{addr, manager, NewReceiveWindow(), NewRecoveryQueue(), mtuSize, Indexes{sync.Mutex{}, make(map[int16][]*protocol.EncapsulatedPacket), make(map[int16]uint)}}
	session.ReceiveWindow.DatagramHandleFunction = func(datagram TimestampedDatagram) {
		session.SendACK(datagram.SequenceNumber)
		session.HandleDatagram(datagram)
	}
	return session
}

// Send sends the given buffer to the session over UDP.
// Returns an int describing the amount of bytes written,
// and an error if unsuccessful.
func (session *Session) Send(buffer []byte) (int, error) {
	return session.Manager.Server.Write(buffer, session.UDPAddr)
}

// SendACK sends an ACK packet to the session for the given sequence number.
// ACKs should only be sent once a datagram is received.
func (session *Session) SendACK(sequenceNumber uint32) {
	ack := protocol.NewACK()
	ack.Packets = []uint32{sequenceNumber}
	ack.Encode()
	session.Send(ack.Buffer)
}

// HandleDatagram handles an incoming datagram encapsulated by a timestamp.
// The actual receive time of the datagram can be checked.
func (session *Session) HandleDatagram(datagram TimestampedDatagram) {
	fmt.Println(datagram.SequenceNumber)
	for _, packet := range *datagram.GetPackets() {
		if packet.HasSplit {
			session.HandleSplitEncapsulated(packet, datagram.Timestamp)
		} else {
			session.HandleEncapsulated(packet, datagram.Timestamp)
		}
	}
}

// HandleEncapsulated handles an encapsulated packet from a datagram.
// A timestamp is passed, which is the timestamp of which the datagram received in the receive window.
func (session *Session) HandleEncapsulated(packet *protocol.EncapsulatedPacket, timestamp int64) {
	switch packet.Buffer[0] {
	case protocol.IdConnectionRequest:

	case protocol.IdNewIncomingConnection:

	case protocol.IdConnectedPing:

	case protocol.IdConnectedPong:

	case protocol.IdDisconnectNotification:

	default:
		session.Manager.EncapsulatedFunction(*packet, session)
	}
}

// HandleSplitEncapsulated handles a split encapsulated packet.
// Split encapsulated packets are first collected into an array,
// and are merged once all fragments of the encapsulated packets have arrived.
func (session *Session) HandleSplitEncapsulated(packet *protocol.EncapsulatedPacket, timestamp int64) {
	id := packet.SplitId
	session.Indexes.Lock()
	if session.Indexes.splits[id] == nil {
		session.Indexes.splits[id] = make([]*protocol.EncapsulatedPacket, packet.SplitCount)
		session.Indexes.splitCounts[id] = 0
	}
	if pk := session.Indexes.splits[id][packet.SplitIndex]; pk == nil {
		session.Indexes.splitCounts[id]++
	}
	session.Indexes.splits[id][packet.SplitIndex] = packet
	if session.Indexes.splitCounts[id] == packet.SplitCount {
		newPacket := protocol.NewEncapsulatedPacket()
		for _, pk := range session.Indexes.splits[id] {
			newPacket.PutBytes(pk.Buffer)
		}
		session.HandleEncapsulated(newPacket, timestamp)
		delete(session.Indexes.splits, id)
	}
	session.Indexes.Unlock()
}

// Tick ticks the session and processes the receive window.
func (session *Session) Tick() {
	session.ReceiveWindow.Tick()
}