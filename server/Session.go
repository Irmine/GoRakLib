package server

import (
	"fmt"
	"strconv"
	"goraklib/protocol"
)

type Session struct {
	address string
	port uint16
	opened bool
	connected bool
	currentSequenceNumber int
	mtuSize int16
	packets chan protocol.IPacket

	sendDatagram chan protocol.Datagram
	clientId uint64
}

func NewSession(address string, port uint16) *Session {
	fmt.Println("Session created for ip: " + address + ":" + strconv.Itoa(int(port)))
	return &Session{address: address, port: port, opened: false, connected: false, packets: make(chan protocol.IPacket, 20), currentSequenceNumber: 0, sendDatagram: make(chan protocol.Datagram, 3)}
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