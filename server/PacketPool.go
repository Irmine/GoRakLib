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
	return &pool
}

func (pool *PacketPool) RegisterPacket(id int, packet protocol.IPacket) {
	pool.packets[id] = packet
}

func (pool *PacketPool) GetPacket(id int) protocol.IPacket {
	return pool.packets[id]
}