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

func (pool *PacketPool) GetPacket(buffer []byte) protocol.IPacket {
	var packet, ok = pool.packets[int(buffer[0])]
	if !ok {
		var header = buffer[0]
		if header & protocol.BitFlagValid == 0 {
			return protocol.NewDatagram() // TODO: Error for invalid.
		}

		if header & protocol.BitFlagIsAck != 0 {
			return protocol.NewACK()
		} else if header & protocol.BitFlagIsNak != 0 {
			return protocol.NewNAK()
		}
		return protocol.NewDatagram()
	}
	return packet()
}