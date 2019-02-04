package server

import (
	"net"
	"fmt"
	"github.com/irmine/goraklib/protocol"
	"sync"
	"time"
)

// Session is a manager of a connection between the client and the server.
// Sessions manage everything related to packet ordering and processing.
type Session struct {
	*net.UDPAddr
	Manager       *Manager
	ReceiveWindow *ReceiveWindow
	RecoveryQueue *RecoveryQueue

	// MTUSize is the maximum size of packets sent and received to and from this sessoin.
	MTUSize 	int16
	// Indexes holds all datagram and encapsulated packet indexes.
	Indexes 	Indexes
	// Queues holds all send queues of the session.
	Queues 		Queues
	// ClientId is the unique client ID of the session.
	ClientId 	uint64
	// CurrentPing is the current latency of the session.
	CurrentPing int64
	// LastUpdate is the last update time of the session.
	LastUpdate	time.Time
	// FlaggedForClose indicates if this session has been flagged to close.
	// Sessions flagged for close will be closed next tick safely.
	FlaggedForClose bool
}

// Queues is a container of four priority queues.
// Immediate priority, high priority, medium priority and low priority queues.
type Queues struct {
	Immediate *PriorityQueue
	High      *PriorityQueue
	Medium    *PriorityQueue
	Low       *PriorityQueue
}

// Indexes is used for the collection of indexes related to datagrams and encapsulated packets.
// It uses several maps and is therefore protected by a mutex.
type Indexes struct {
	sync.Mutex
	splits        map[int16][]*protocol.EncapsulatedPacket
	splitCounts   map[int16]uint
	splitId       int16
	sendSequence  uint32
	messageIndex  uint32
	orderIndex    uint32 // TODO: Implement proper order channels and indexes.
	sequenceIndex uint32
}

// NewSession returns a new session with UDP address.
// The MTUSize provided is the maximum packet size of the session.
func NewSession(addr *net.UDPAddr, mtuSize int16, manager *Manager) *Session {
	session := &Session{addr,
		manager,
		NewReceiveWindow(),
		NewRecoveryQueue(),
		mtuSize,
		Indexes{sync.Mutex{}, make(map[int16][]*protocol.EncapsulatedPacket), make(map[int16]uint), 0, 0, 0, 0, 0},
		Queues{NewPriorityQueue(1), NewPriorityQueue(256), NewPriorityQueue(256), NewPriorityQueue(256)},
		0,
		0,
		time.Now(),
		false,
	}
	session.ReceiveWindow.DatagramHandleFunction = func(datagram TimestampedDatagram) {
		session.LastUpdate = time.Now()
		session.SendACK(datagram.SequenceNumber)
		session.HandleDatagram(datagram)
	}
	return session
}

// Close removes all references of the session,
// and removes the capability to send and handle packets.
// Sessions can not be opened once closed.
// It is strongly unrecommended to use this function directly.
// Use FlagForClose instead.
func (session *Session) Close() {
	session.Manager.DisconnectFunction(session)
	session.UDPAddr = nil
	session.Manager = nil
	session.ReceiveWindow = nil
	session.RecoveryQueue = nil
	session.Indexes = Indexes{}
	session.Queues = Queues{}
}

// FlagForClose flags the session for close.
// It is always recommended to use this function over direct Close.
// Sessions flagged for close will be closed the next tick.
func (session *Session) FlagForClose() {
	session.FlaggedForClose = true
}

// IsClosed checks if the session is closed.
// Sending and handling packets for a session is
// impossible once the session is closed.
func (session *Session) IsClosed() bool {
	return session.UDPAddr == nil
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
	for _, packet := range *datagram.GetPackets() {
		if packet.HasSplit {
			session.HandleSplitEncapsulated(packet, datagram.Timestamp)
		} else {
			session.HandleEncapsulated(packet, datagram.Timestamp)
		}
	}
}

// HandleACK handles an incoming ACK packet.
// Recovery gets removed for every datagram with a sequence number in the ACK.
func (session *Session) HandleACK(ack *protocol.ACK) {
	session.RecoveryQueue.RemoveRecovery(ack.Packets)
}

// HandleNACK handles an incoming NACK packet.
// All packets requested in the NACK get resent to the client.
func (session *Session) HandleNACK(nack *protocol.NAK) {
	for _, seq := range nack.Packets {
		if !session.RecoveryQueue.IsRecoverable(seq) {
			fmt.Println("Unrecoverable datagram:", seq, "(", nack.Packets, ")")
		}
	}
	datagrams, _ := session.RecoveryQueue.Recover(nack.Packets)
	for _, datagram := range datagrams {
		session.Manager.Server.Write(datagram.Buffer, session.UDPAddr)
	}
}

// HandleEncapsulated handles an encapsulated packet from a datagram.
// A timestamp is passed, which is the timestamp of which the datagram received in the receive window.
func (session *Session) HandleEncapsulated(packet *protocol.EncapsulatedPacket, timestamp int64) {
	session.LastUpdate = time.Now()
	switch packet.Buffer[0] {
	case protocol.IdConnectionRequest:
		session.HandleConnectionRequest(packet)
	case protocol.IdNewIncomingConnection:
		session.Manager.ConnectFunction(session)
	case protocol.IdConnectedPing:
		session.HandleConnectedPing(packet, timestamp)
	case protocol.IdConnectedPong:
		session.HandleConnectedPong(packet, timestamp)
	case protocol.IdDisconnectNotification:
		session.FlagForClose()
	default:
		session.Manager.PacketFunction(packet.Buffer, session)
	}
}

// HandleConnectedPong handles a pong reply of our own sent ping.
func (session *Session) HandleConnectedPong(packet *protocol.EncapsulatedPacket, timestamp int64) {
	pong := protocol.NewConnectedPong()
	pong.Buffer = packet.Buffer
	pong.Decode()
	session.CurrentPing = (timestamp - pong.PongSendTime) / int64(time.Millisecond)
}

// HandleConnectedPing handles a connected ping from the client.
// A pong is sent back at low priority.
func (session *Session) HandleConnectedPing(packet *protocol.EncapsulatedPacket, timestamp int64) {
	ping := protocol.NewConnectedPing()
	ping.Buffer = packet.Buffer
	ping.Decode()

	pong := protocol.NewConnectedPong()
	pong.PingSendTime = ping.PingSendTime
	pong.PongSendTime = timestamp

	session.SendPacket(pong, protocol.ReliabilityUnreliable, PriorityLow)
}

// HandleConnectionRequest handles a connection request from the session.
// A connection accept gets sent back to the client.
func (session *Session) HandleConnectionRequest(packet *protocol.EncapsulatedPacket) {
	request := protocol.NewConnectionRequest()
	request.Buffer = packet.GetBuffer()
	request.Decode()

	session.ClientId = request.ClientId

	accept := protocol.NewConnectionAccept()
	accept.ClientAddress = session.UDPAddr.IP.String()
	accept.ClientPort = uint16(session.UDPAddr.Port)

	accept.PingSendTime = uint64(time.Now().Unix())
	accept.PongSendTime = uint64(time.Now().Unix())

	session.SendPacket(accept, protocol.ReliabilityReliableOrdered, PriorityImmediate)
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

// Tick ticks the session and processes the receive window and priority queues.
// currentTick is the current tick of the server, which increments every time this function is ran.
func (session *Session) Tick(currentTick int64) {
	session.Queues.High.Flush(session)
	if currentTick % 400 == 0 {
		ping := protocol.NewConnectedPing()
		ping.PingSendTime = time.Now().Unix()
		session.SendPacket(ping, protocol.ReliabilityUnreliable, PriorityImmediate)
	}
	if currentTick % 2 == 0 {
		session.ReceiveWindow.Tick()
		session.Queues.Medium.Flush(session)
	}
	if currentTick % 4 == 0 {
		session.Queues.Low.Flush(session)
	}
}

// SendPacket sends an external packet to a session.
// The reliability given will be added to the encapsulated packet.
// The packet will be added with the given priority. Immediate priority packets are sent out immediately.
func (session *Session) SendPacket(packet protocol.IConnectedPacket, reliability byte, priority Priority) {
	packet.Encode()
	encapsulated := protocol.NewEncapsulatedPacket()
	encapsulated.Reliability = reliability
	encapsulated.Buffer = packet.GetBuffer()
	session.Queues.AddEncapsulated(encapsulated, priority, session)
}

// AddEncapsulated adds an encapsulated packet at the given priority.
// The queue gets flushed immediately if the priority is immediate priority.
func (queues Queues) AddEncapsulated(packet *protocol.EncapsulatedPacket, priority Priority, session *Session) {
	if session.IsClosed() {
		return
	}
	var queue *PriorityQueue
	switch priority {
	case PriorityImmediate:
		queue = queues.Immediate
	case PriorityHigh:
		queue = queues.High
	case PriorityMedium:
		queue = queues.Medium
	case PriorityLow:
		queue = queues.Low
	}
	queue.AddEncapsulated(packet, session)
	if priority == PriorityImmediate {
		queue.Flush(session)
	}
}