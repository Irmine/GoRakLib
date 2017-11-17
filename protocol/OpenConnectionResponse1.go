package protocol

import "goraklib/protocol/identifiers"

type OpenConnectionResponse1 struct {
	*UnconnectedMessage
	ServerId int64
	MtuSize int16
	Security byte
}

func NewOpenConnectionResponse1() *OpenConnectionResponse1 {
	return &OpenConnectionResponse1{NewUnconnectedMessage(NewPacket(
		identifiers.OpenConnectionResponse1,
	)), 0, 0, 0}
}

func (response *OpenConnectionResponse1) Encode() {
	response.EncodeId()
	response.WriteMagic()
	response.PutLong(response.ServerId)
	response.PutByte(response.Security)
	response.PutShort(response.MtuSize)
}

func (response *OpenConnectionResponse1) Decode() {

}
