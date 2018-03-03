package server

import (
	"github.com/irmine/binutils"
)

// A RawPacket is used to identify packets that do not follow the RakNet protocol.
// Raw packets may be of different protocols, query protocols as an example.
// They are simply ignored and forwarded to the managing program.
// Raw packets serve no purpose other than simply forwarding data.
type RawPacket struct {
	*binutils.Stream
}

// NewRawPacket returns a new raw packet.
func NewRawPacket() RawPacket {
	return RawPacket{Stream: binutils.NewStream()}
}

func (pk RawPacket) Encode() {}

func (pk RawPacket) Decode() {}

func (pk RawPacket) GetId() int {return -1}

func (pk RawPacket) HasMagic() bool {return false}
