package server

import (
	"strconv"
	"sync"
	"time"

	"github.com/irmine/goraklib/protocol"
)

type Session struct {
	manager       *SessionManager
	queue         *PriorityQueue
	recoveryQueue *RecoveryQueue

	address string
	port    uint16

	opened    bool
	connected bool
	closed    bool

	currentSequenceNumber uint32
	mtuSize               int16

	packetBatches chan protocol.EncapsulatedPacket

	splits      map[int][]*protocol.EncapsulatedPacket
	splitCounts map[int]uint

	clientId uint64

	ping uint64

	orderIndex   sync.Map
	messageIndex uint32
	splitId      int16

	LastUpdate int64

	forceClose  bool
	timeoutTick int

	receiveWindow   *ReceiveWindow
	receiveSequence uint32
}

func NewSession(manager *SessionManager, address string, port uint16) *Session {
	var session = &Session{splitCounts: make(map[int]uint), recoveryQueue: NewRecoveryQueue(), orderIndex: sync.Map{}, manager: manager, address: address, port: port, opened: false, connected: false, splits: make(map[int][]*protocol.EncapsulatedPacket), packetBatches: make(chan protocol.EncapsulatedPacket, 512), currentSequenceNumber: 1}
	session.queue = NewPriorityQueue(session)
	session.LastUpdate = time.Now().Unix()
	session.receiveWindow = NewReceiveWindow(session)
	return session
}

func (session *Session) SendUnconnectedPacket(packet protocol.IPacket) {
	packet.Encode()

	session.manager.server.udp.WriteBuffer(packet.GetBuffer(), session.GetAddress(), session.GetPort())
}

func (session *Session) SendConnectedPacket(packet protocol.IConnectedPacket, reliability byte, priority byte) {
	var encapsulatedPacket = protocol.NewEncapsulatedPacket()
	encapsulatedPacket.Pk = packet
	encapsulatedPacket.OrderChannel = 0
	encapsulatedPacket.Reliability = reliability

	if priority != PriorityImmediate {
		session.queue.AddEncapsulatedToQueue(encapsulatedPacket, priority)
		return
	}

	packet.Encode()
	encapsulatedPacket.Buffer = packet.GetBuffer()
	encapsulatedPacket.Pk = nil

	if encapsulatedPacket.IsReliable() {
		encapsulatedPacket.MessageIndex = session.messageIndex
		session.messageIndex++
	}

	if encapsulatedPacket.IsSequenced() {
		i, _ := session.orderIndex.Load(encapsulatedPacket.OrderChannel)
		if i == nil {
			i = uint32(0)
		}
		encapsulatedPacket.OrderIndex = i.(uint32)

		session.orderIndex.Store(encapsulatedPacket.OrderChannel, i.(uint32)+1)
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

	session.manager.HandlePacket(packet, session)
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
