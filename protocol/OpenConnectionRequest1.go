package protocol

import (
	"goraklib/protocol/identifiers"
)

type OpenConnectionRequest1 struct {
	*UnconnectedMessage
}

func NewOpenConnectionRequest1() *OpenConnectionRequest1 {
	return &OpenConnectionRequest1{NewUnconnectedMessage(NewPacket(
		identifiers.OpenConnectionRequest1,
	))}
}

func (request *OpenConnectionRequest1) Encode() {

}

func (request *OpenConnectionRequest1) Decode() {
	request.DecodeStep()
}
