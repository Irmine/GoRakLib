package server

import (
	"goraklib/protocol"
	"goraklib/protocol/identifiers"
)

type PacketPool struct {
	packets map[int]protocol.IPacket
}

func NewPacketPool() *PacketPool {
	var pool = PacketPool{}
	pool.packets = make(map[int]protocol.IPacket)

	pool.RegisterPacket(identifiers.UnconnectedPing, protocol.NewUnconnectedPing())
	pool.RegisterPacket(identifiers.OpenConnectionRequest1, protocol.NewOpenConnectionRequest1())
	pool.RegisterPacket(identifiers.OpenConnectionRequest2, protocol.NewOpenConnectionRequest2())
	return &pool
}

func (pool *PacketPool) RegisterPacket(id int, packet protocol.IPacket) {
	pool.packets[id] = packet
}

func (pool *PacketPool) GetPacket(id int) protocol.IPacket {
	var packet, ok = pool.packets[id]
	if !ok {
		return protocol.NewDatagram()
	}
	return packet
}