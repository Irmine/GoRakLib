package server

import (
	"net"
	"fmt"
	"os"
)

type UDPServer struct {
	Conn *net.UDPConn
	pool *PacketPool
}

func NewServer(port int) UDPServer {

	server := UDPServer{}

	var addr, err = net.ResolveUDPAddr("udp", ":19132")
	addr.Port = port

	conn, err := net.ListenUDP("udp", addr)

	if err != nil {
		//todo
		fmt.Printf("An error has occurred: %v", err)
		os.Exit(1)
	}

	server.Conn = conn
	server.pool = NewPacketPool()

	return server
}

func (udp *UDPServer) ReadBuffer() string {

	var buffer = make([]byte, 1028)

	n, addr, err := udp.Conn.ReadFromUDP(buffer)

	if err != nil {
		fmt.Printf("An error has occurred: %v", err)
		os.Exit(1)
	}

	if n == 0 {
		return ""
	}

	var idBuffer = buffer
	var packetId = int(idBuffer[0])

	var packet = udp.pool.GetPacket(packetId)

	packet.SetBuffer(buffer)
	packet.Decode()

	return addr.IP.To4().String()
}

func (udp *UDPServer) writeBuffer(buffer *[]byte, ip string, port int) {

	addr := net.UDPAddr{
		IP: net.ParseIP(ip),
		Port: port,
	}

	_, err := udp.Conn.WriteToUDP(*buffer, &addr)

	if err != nil {
		fmt.Printf("An error has occurred: %v", err)
		os.Exit(1)
	}
}