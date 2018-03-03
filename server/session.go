package server

import "net"

type Session struct {
	*net.UDPAddr

}
