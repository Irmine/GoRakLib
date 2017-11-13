package server

import (
	"net"
	"fmt"
	"os"
	"strconv"
	"errors"
	"goraklib/protocol"
)

type UDPServer struct {
	port int
	address string
	Conn *net.UDPConn
	pool *PacketPool
	packets []protocol.IPacket
}

func NewUDPServer(port int) UDPServer {

	server := UDPServer{}

	server.port = port
	var addr, err = net.ResolveUDPAddr("udp", ":" + strconv.Itoa(port))
	addr.Port = port

	server.address = addr.IP.To4().String()

	conn, err := net.ListenUDP("udp", addr)

	if err != nil {
		fmt.Printf("An error has occurred: %v", err)
		os.Exit(1)
	}

	server.Conn = conn
	server.pool = NewPacketPool()

	return server
}

func (udp *UDPServer) GetPort() int {
	return udp.port
}

func (udp *UDPServer) ReadBuffer() (protocol.IPacket, string, int, error) {

	var buffer = make([]byte, 4096)

	n, addr, err := udp.Conn.ReadFromUDP(buffer)

	if err != nil {
		fmt.Printf("An error has occurred: %v", err)
		os.Exit(1)
	}

	var packet protocol.IPacket

	if n == 0 {
		return packet, "", 0, errors.New("received null packet")
	}

	var ip = addr.IP.To4().String()
	var port = addr.Port

	var idBuffer = buffer
	var packetId = int(idBuffer[0])

	packet = udp.pool.GetPacket(packetId)
	if packet == nil {
		fmt.Println("Unknown package with ID:", packetId)
	}

	packet.SetBuffer(buffer)

	return packet, ip, port, nil
}

func (udp *UDPServer) WriteBuffer(buffer []byte, ip string, port int) {

	addr := net.UDPAddr{
		IP: net.ParseIP(ip),
		Port: port,
	}

	_, err := udp.Conn.WriteToUDP(buffer, &addr)

	if err != nil {
		fmt.Printf("An error has occurred: %v", err)
		os.Exit(1)
	}
}