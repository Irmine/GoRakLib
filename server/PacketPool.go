package server

import (
	"goraklib/protocol"
	"goraklib/protocol/identifiers"
)

type PacketPool struct {
	packets map[int]func() protocol.IPacket
}

func NewPacketPool() *PacketPool {
	var pool = PacketPool{}
	pool.packets = make(map[int]func() protocol.IPacket)

	pool.RegisterPacket(identifiers.UnconnectedPing, func() protocol.IPacket { return protocol.NewUnconnectedPing() })
	pool.RegisterPacket(identifiers.OpenConnectionRequest1, func() protocol.IPacket { return protocol.NewOpenConnectionRequest1() })
	pool.RegisterPacket(identifiers.OpenConnectionRequest2, func() protocol.IPacket { return protocol.NewOpenConnectionRequest2() })
	return &pool
}

func (pool *PacketPool) RegisterPacket(id int, packet func() protocol.IPacket) {
	pool.packets[id] = packet
}

func (pool *PacketPool) GetPacket(id int) protocol.IPacket {
	var packet, ok = pool.packets[id]
	if !ok {
		return protocol.NewDatagram()
	}
	return packet()
}