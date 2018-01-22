package server

import (
	"goraklib/binary"
)

type RawPacket struct {
	*binary.BinaryStream
	Address string
	Port uint16
}

func NewRawPacket() RawPacket {
	return RawPacket{BinaryStream: binary.NewStream()}
}

func (pk RawPacket) Encode() {

}

func (pk RawPacket) Decode() {

}

func (pk RawPacket) GetId() int {
	return -1
}

func (pk RawPacket) HasMagic() bool {
	return false
}