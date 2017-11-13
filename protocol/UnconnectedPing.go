package protocol

import (
	"goraklib/protocol/identifiers"
	"fmt"
)

type UnconnectedPing struct {
	*UnconnectedMessage
	pingId int64
}

func NewUnconnectedPing() UnconnectedPing {
	return UnconnectedPing{NewUnconnectedMessage(NewPacket(
		identifiers.UnconnectedPing,
	)), 0}
}

func (ping UnconnectedPing) Encode() {

}

func (ping UnconnectedPing) Decode() {
	ping.DecodeStep()
	ping.pingId = ping.GetLong()
	ping.ReadMagic()
	if ping.HasValidMagic() {
		fmt.Println("Received a valid UnconnectedPing")
	}
}