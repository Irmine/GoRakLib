package protocol

import (
	"goraklib/protocol/identifiers"
)

type OpenConnectionResponse2 struct {
	*UnconnectedMessage
	ServerId int64
	MtuSize int16
	ClientAddress string
	ClientPort uint16
	UseEncryption bool
}

func NewOpenConnectionResponse2() *OpenConnectionResponse2 {
	return &OpenConnectionResponse2{NewUnconnectedMessage(NewPacket(
		identifiers.OpenConnectionResponse2,
	)), 0, 0, "", 0, false}
}

func (response *OpenConnectionResponse2) Encode() {
	response.EncodeId()
	response.PutMagic()
	response.PutLong(response.ServerId)
	response.PutAddress(response.ClientAddress, response.ClientPort, 4)
	response.PutShort(response.MtuSize)
	response.PutBool(response.UseEncryption)
}

func (response *OpenConnectionResponse2) Decode() {
	response.DecodeStep()
	response.ReadMagic()
	response.ServerId = response.GetLong()
	response.ClientAddress, response.ClientPort, _ = response.GetAddress()
	response.MtuSize = response.GetShort()
	response.UseEncryption = response.GetBool()
}

