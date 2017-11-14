package server

import (
	"fmt"
	"strconv"
	"goraklib/protocol"
)

type Session struct {
	address string
	port int
	opened bool
	connected bool
	packets chan protocol.IPacket
}

func NewSession(address string, port int) *Session {
	fmt.Println("Session created for ip: " + address + ":" + strconv.Itoa(port))
	return &Session{address: address, port: port, opened: false, packets: make(chan protocol.IPacket, 20)}
}

func (session *Session) GetAddress() string {
	return session.address
}

func (session *Session) GetPort() int {
	return session.port
}

func (session *Session) Open() {
	session.opened = true
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

