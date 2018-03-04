package server

import (
	"net"
	"errors"
)

// UDPServer is a wrapper around a UDPConn.
// It can be started on a given address and port,
// and provides functions to read and write packets to the connection.
type UDPServer struct {
	*net.UDPConn
}

// NotStarted is an error returned for the Read and Write functions if the server has not yet been started.
var NotStarted = errors.New("udp server has not started")

// NewUDPServer returns a new UDP server.
// The UDPServer will not have a default connection,
// and all actions executed on it before starting will fail.
func NewUDPServer() *UDPServer {
	return &UDPServer{}
}

// Start starts the UDP server on the given address and port.
// An error is returned if ListenUDP is not successful.
// Actions can be used on the UDP server once started.
func (server *UDPServer) Start(address string, port int) error {
	addr := &net.UDPAddr{IP: net.ParseIP(address), Port: port}
	var err error
	server.UDPConn, err = net.ListenUDP("udp", addr)
	return err
}

// HasStarted checks if a UDPServer has been started.
// No actions can be executed on the UDPServer while not started.
func (server *UDPServer) HasStarted() bool {
	return server.UDPConn != nil
}

// Read reads any data from the UDP connection into the given byte array.
// The IP address and port of the client that sent the data will be returned,
// along with an error that might have occurred during reading.
func (server *UDPServer) Read(buffer []byte) (bytesRead int, addr *net.UDPAddr, err error) {
	if !server.HasStarted() {
		return 0, nil, NotStarted
	}
	bytesRead, addr, err = server.UDPConn.ReadFromUDP(buffer)
	return
}

// Write writes a byte array to a UDP connection.
// Write returns the amount of bytes written and an error that might have occurred.
func (server *UDPServer) Write(buffer []byte, addr *net.UDPAddr) (int, error) {
	if !server.HasStarted() {
		return 0, NotStarted
	}
	return server.UDPConn.WriteToUDP(buffer, addr)
}