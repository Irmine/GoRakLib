package protocol

import (
	"goraklib/protocol/identifiers"
	"strconv"
)

type UnconnectedPong struct {
	*UnconnectedMessage
	PingId int64
	ServerId int64
	ServerName string
	ServerProtocol uint
	ServerVersion string
	OnlinePlayers uint
	MaximumPlayers uint
	Motd string
	DefaultGameMode string
}

func NewUnconnectedPong() *UnconnectedPong {
	return &UnconnectedPong{NewUnconnectedMessage(NewPacket(
		identifiers.UnconnectedPong,
	)), 0, 0, "", 0, "", 0, 20, "", ""}
}

func (pong *UnconnectedPong) Encode() {
	pong.EncodeId()
	pong.PutLong(pong.PingId)
	pong.PutLong(pong.ServerId)
	pong.WriteMagic()
	pong.PutString(
		"MCPE;" +
		pong.ServerName + ";" +
		strconv.Itoa(int(pong.ServerProtocol)) + ";" +
		pong.ServerVersion + ";" +
		strconv.Itoa(int(pong.OnlinePlayers)) + ";" +
		strconv.Itoa(int(pong.MaximumPlayers)) + ";" +
		strconv.Itoa(int(pong.ServerId)) + ";" +
		pong.Motd + ";" +
		pong.DefaultGameMode + ";",
	)
}

func (pong *UnconnectedPong) Decode() {
	pong.DecodeStep()
	pong.PingId = pong.GetLong()
	pong.ServerId = pong.GetLong()
	pong.ReadMagic()
}
