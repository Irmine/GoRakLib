package server

import (
	"net"
	"strconv"
	"errors"
	"goraklib/protocol"
)

type UDPServer struct {
	port uint16
	address string
	Conn *net.UDPConn
	pool *PacketPool
}

func NewUDPServer(address string, port uint16) UDPServer {
	server := UDPServer{}

	server.port = port

	if address == "127.0.0.1" {
		address = ""
	}

	var addr, err = net.ResolveUDPAddr("udp", address + ":" + strconv.Itoa(int(port)))
	addr.Port = int(port)

	server.address = addr.IP.To4().String()

	conn, err := net.ListenUDP("udp", addr)

	if err != nil {
		panic(err)
	}

	server.Conn = conn
	server.pool = NewPacketPool()

	return server
}

func (udp *UDPServer) GetPort() uint16 {
	return udp.port
}

func (udp *UDPServer) ReadBuffer(server *GoRakLibServer) (protocol.IPacket, string, uint16, error) {
	var buffer = make([]byte, 4096)

	n, addr, err := udp.Conn.ReadFromUDP(buffer)

	if err != nil {
		panic(err)
	}

	buffer = buffer[:n]
	var packet protocol.IPacket

	if n == 0 {
		return packet, "", 0, errors.New("received null packet")
	}

	var ip = addr.IP.To4().String()
	var port = uint16(addr.Port)

	packet = udp.pool.GetPacket(buffer, server.sessionManager.SessionExists(ip, port))

	packet.SetBuffer(buffer)

	return packet, ip, port, nil
}

func (udp *UDPServer) WriteBuffer(buffer []byte, ip string, port uint16) {

	addr := net.UDPAddr{
		IP: net.ParseIP(ip),
		Port: int(port),
	}

	_, err := udp.Conn.WriteToUDP(buffer, &addr)

	if err != nil {
		panic(err)
	}
}