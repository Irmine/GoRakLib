package protocol

import "goraklib/protocol/identifiers"

type OpenConnectionResponse1 struct {
	*UnconnectedMessage
	ServerId int64
	MtuSize int16
	Security bool
}

func NewOpenConnectionResponse1() *OpenConnectionResponse1 {
	return &OpenConnectionResponse1{NewUnconnectedMessage(NewPacket(
		identifiers.OpenConnectionResponse1,
	)), 0, 0, false}
}

func (response *OpenConnectionResponse1) Encode() {
	response.EncodeId()
	response.PutMagic()
	response.PutLong(response.ServerId)
	response.PutBool(response.Security)
	response.PutShort(response.MtuSize)
}

func (response *OpenConnectionResponse1) Decode() {
	response.DecodeStep()
	response.ReadMagic()
	response.ServerId = response.GetLong()
	response.Security = response.GetBool()
	response.MtuSize = response.GetShort()
}
