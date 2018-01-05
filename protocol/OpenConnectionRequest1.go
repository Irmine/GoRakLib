package protocol

import (
	"goraklib/protocol/identifiers"
)

type OpenConnectionRequest1 struct {
	*UnconnectedMessage
	Protocol byte
	MtuSize int16
}

func NewOpenConnectionRequest1() *OpenConnectionRequest1 {
	return &OpenConnectionRequest1{NewUnconnectedMessage(NewPacket(
		identifiers.OpenConnectionRequest1,
	)), 0, 0}
}

func (request *OpenConnectionRequest1) Encode() {
	request.EncodeId()
	request.PutMagic()
	request.PutByte(request.Protocol)

	var bytes = make([]byte, request.MtuSize + 18)
	request.PutBytes(bytes)
}

func (request *OpenConnectionRequest1) Decode() {
	request.DecodeStep()
	request.ReadMagic()
	request.Protocol = request.GetByte()
	request.MtuSize = int16(len(request.Get(-1)) + 18)
}
