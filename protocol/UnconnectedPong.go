package protocol

import (
	"goraklib/protocol/identifiers"
	"strconv"
)

type UnconnectedPong struct {
	*UnconnectedMessage
	PingId int64
	ServerId int64
	ServerData string
}

func NewUnconnectedPong() *UnconnectedPong {
	return &UnconnectedPong{NewUnconnectedMessage(NewPacket(
		identifiers.UnconnectedPong,
	)), 0, 0, ""}
}

func (pong *UnconnectedPong) Encode() {
	pong.EncodeId()
	pong.PutLong(pong.PingId)
	pong.PutLong(pong.ServerId)
	pong.WriteMagic()
	pong.PutString("MCPE;GoMineServer;140;1.2.6.2;5;100;" + strconv.Itoa(int(pong.ServerId)) + ";GoMine;Creative;")
}

func (pong *UnconnectedPong) Decode() {
	pong.DecodeStep()
	pong.PingId = pong.GetLong()
	pong.ServerId = pong.GetLong()
	pong.ReadMagic()
}
