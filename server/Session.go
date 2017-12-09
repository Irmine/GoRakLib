package server

import (
	"fmt"
	"strconv"
	"goraklib/protocol"
)

type Session struct {
	address string
	port uint16

	opening bool
	opened bool
	connected bool

	currentSequenceNumber uint32
	mtuSize int16

	packets chan protocol.IPacket
	packetBatches chan protocol.EncapsulatedPacket

	splits map[int]chan *protocol.EncapsulatedPacket

	clientId uint64

	ping uint64
}

func NewSession(address string, port uint16) *Session {
	var session = &Session{address: address, port: port, opened: false, connected: false, splits: make(map[int]chan *protocol.EncapsulatedPacket), packets: make(chan protocol.IPacket, 20), packetBatches: make(chan protocol.EncapsulatedPacket, 512), currentSequenceNumber: 1}
	fmt.Println("Session created for ip:", session)
	return session
}

func (session *Session) Open() {
	session.opened = true
}

func (session *Session) Close() {
	session.opened = false
	session.SetConnected(false)
}

func (session *Session) IsOpened() bool {
	return session.opened
}

func (session *Session) SetConnected(value bool) {
	session.connected = value
}

func (session *Session) IsConnected() bool {
	return session.connected
}

func (session *Session) Forward(packet protocol.IPacket) {
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