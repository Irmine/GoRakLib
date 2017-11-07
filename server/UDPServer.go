package server

import (
	"net"
	"fmt"
	"os"
)

type UDPServer struct {
	Conn *net.UDPConn
}

func NewServer(port int) UDPServer {

	server := UDPServer{}

	addr := net.UDPAddr{
		IP: net.ParseIP("127.0.0.1"),
		Port: port,
	}

	conn, err := net.ListenUDP("udp", &addr)

	if err != nil {
		//todo
		fmt.Printf("An error has occurred: %v", err)
		os.Exit(1)
	}

	server.Conn = conn

	return server
}

func (udp *UDPServer) readBuffer(buffer *[]byte) string {

	_, addr, err := udp.Conn.ReadFromUDP(*buffer)

	if err != nil {
		fmt.Printf("An error has occurred: %v", err)
		os.Exit(1)
	}

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