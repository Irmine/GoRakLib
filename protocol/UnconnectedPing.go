package protocol

import (
	"goraklib/protocol/identifiers"
)

type UnconnectedPing struct {
	*UnconnectedMessage
	PingId int64
}

func NewUnconnectedPing() *UnconnectedPing {
	return &UnconnectedPing{NewUnconnectedMessage(NewPacket(
		identifiers.UnconnectedPing,
	)), 0}
}

func (ping *UnconnectedPing) Encode() {

}

func (ping *UnconnectedPing) Decode() {
	ping.DecodeStep()
	ping.PingId = ping.GetLong()
	ping.ReadMagic()
}