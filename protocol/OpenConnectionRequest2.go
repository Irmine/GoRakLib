package protocol

import (
	"goraklib/protocol/identifiers"
)

type OpenConnectionRequest2 struct {
	*UnconnectedMessage
	ServerAddress string
	ServerPort uint16
	MtuSize int16
	ClientId int64
}

func NewOpenConnectionRequest2() *OpenConnectionRequest2 {
	return &OpenConnectionRequest2{NewUnconnectedMessage(NewPacket(
		identifiers.OpenConnectionRequest2,
	)), "", 0, 0, 0}
}

func (request *OpenConnectionRequest2) Encode() {

}

func (request *OpenConnectionRequest2) Decode() {
	request.DecodeStep()
	request.ReadMagic()
	var address, port, _ = request.GetAddress()
	request.ServerAddress = address
	request.ServerPort = port
	request.MtuSize = request.GetShort()
	request.ClientId = request.GetLong()
}