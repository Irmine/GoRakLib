package server

import (
	"strconv"
	"goraklib/protocol"
	"time"
)

type Session struct {
	manager *SessionManager
	queue *PriorityQueue
	recoveryQueue *RecoveryQueue

	address string
	port uint16

	opened bool
	connected bool
	closed bool

	currentSequenceNumber uint32
	mtuSize int16

	packets chan protocol.IPacket
	packetBatches chan protocol.EncapsulatedPacket

	splits map[int]chan *protocol.EncapsulatedPacket

	clientId uint64

	ping uint64

	orderIndex map[byte]uint32
	messageIndex uint32
	splitId int16

	LastUpdate int64

	forceClose bool
	timeoutTick int
}

func NewSession(manager *SessionManager, address string, port uint16) *Session {
	var session = &Session{recoveryQueue: NewRecoveryQueue(), orderIndex: make(map[byte]uint32), manager: manager, address: address, port: port, opened: false, connected: false, splits: make(map[int]chan *protocol.EncapsulatedPacket), packets: make(chan protocol.IPacket, 20), packetBatches: make(chan protocol.EncapsulatedPacket, 512), currentSequenceNumber: 1}
	session.queue = NewPriorityQueue(session)
	session.LastUpdate = time.Now().Unix()
	return session
}

func (session *Session) SendUnconnectedPacket(packet protocol.IPacket) {
	session.LastUpdate = time.Now().Unix()
	packet.Encode()

	session.manager.server.udp.WriteBuffer(packet.GetBuffer(), session.GetAddress(), session.GetPort())
}

func (session *Session) SendConnectedPacket(packet protocol.IConnectedPacket, reliability byte, priority byte) {
	packet.Encode()

	var encapsulatedPacket = protocol.NewEncapsulatedPacket()
	encapsulatedPacket.Buffer = packet.GetBuffer()
	encapsulatedPacket.OrderChannel = 0
	encapsulatedPacket.Reliability = reliability

	if priority != PriorityImmediate {
		session.queue.AddEncapsulatedToQueue(encapsulatedPacket, priority)
		return
	}

	if encapsulatedPacket.IsReliable() {
		encapsulatedPacket.MessageIndex = session.messageIndex
		session.messageIndex++
	}
	if encapsulatedPacket.IsSequenced() {
		encapsulatedPacket.OrderIndex = session.orderIndex[encapsulatedPacket.OrderChannel]
		session.orderIndex[encapsulatedPacket.OrderChannel]++
	}

	var datagram = protocol.NewDatagram()
	datagram.NeedsBAndAs = true

	datagram.SequenceNumber = session.currentSequenceNumber
	session.currentSequenceNumber++

	datagram.AddPacket(encapsulatedPacket)
	session.manager.SendPacket(datagram, session)

	session.recoveryQueue.AddRecoveryFor(datagram)
}

func (session *Session) Open() {
	session.opened = true
}

func (session *Session) Close(force bool) {
	session.queue.Wipe()

	session.closed = true
	if force {
		session.forceClose = true
	}
	session.timeoutTick++
}

func (session *Session) IsOpened() bool {
	return session.opened
}

func (session *Session) IsClosed() bool {
	return session.closed
}

func (session *Session) IsReadyForDeletion() bool {
	return session.closed && (session.recoveryQueue.IsClear() || session.forceClose || session.timeoutTick > 20)
}

func (session *Session) SetConnected(value bool) {
	session.connected = value
}

func (session *Session) IsConnected() bool {
	return session.connected
}

func (session *Session) Forward(packet protocol.IPacket) {
	session.LastUpdate = time.Now().Unix()
	packet.Decode()

	session.packets <- packet
}

func (session *Session) IsStackEmpty() bool {
	return len(session.packets) == 0
}

func (session *Session) FetchFromStack() protocol.IPacket {
	return <- session.packets
}

func (session *Session) GetPort() uint16 {
	return session.port
}

func (session *Session) GetAddress() string {
	return session.address
}

func (session *Session) GetMtuSize() int16 {
	return session.mtuSize
}

func (session *Session) GetClientId() uint64 {
	return session.clientId
}

func (session *Session) GetReadyEncapsulatedPackets() []protocol.EncapsulatedPacket {
	var packets []protocol.EncapsulatedPacket
	for len(session.packetBatches) != 0 {
		packets = append(packets, <-session.packetBatches)
	}
	return packets
}

func (session *Session) AddProcessedEncapsulatedPacket(packet protocol.EncapsulatedPacket) {
	session.packetBatches <- packet
}

func (session *Session) GetPing() uint64 {
	return session.ping
}

func (session *Session) SetPing(ping uint64) {
	session.ping = ping
}

func (session *Session) String() string {
	return session.address + ":" + strconv.Itoa(int(session.port))
}