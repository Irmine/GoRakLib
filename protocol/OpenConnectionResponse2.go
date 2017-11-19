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
	Security byte
}

func NewOpenConnectionResponse2() *OpenConnectionResponse2 {
	return &OpenConnectionResponse2{NewUnconnectedMessage(NewPacket(
		identifiers.OpenConnectionResponse2,
	)), 0, 0, "", 0, 0}
}

func (response *OpenConnectionResponse2) Encode() {
	response.EncodeId()
	response.WriteMagic()
	response.PutLong(response.ServerId)
	response.PutAddress(response.ClientAddress, response.ClientPort, 4)
	response.PutShort(response.MtuSize)
	response.PutByte(response.Security)
}

func (response *OpenConnectionResponse2) Decode() {

}

