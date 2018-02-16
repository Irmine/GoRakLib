package protocol

import "github.com/irmine/goraklib/protocol/identifiers"

type ConnectedPong struct {
	*Packet
	PingSendTime int64
	PongSendTime int64
}

func NewConnectedPong() *ConnectedPong {
	return &ConnectedPong{NewPacket(
		identifiers.ConnectedPong,
	), 0, 0}
}

func (pong *ConnectedPong) Encode() {
	pong.EncodeId()
	pong.PutLong(pong.PingSendTime)
	pong.PutLong(pong.PongSendTime)
}

func (pong *ConnectedPong) Decode() {
	pong.DecodeStep()
	pong.PingSendTime = pong.GetLong()
	pong.PongSendTime = pong.GetLong()
}
