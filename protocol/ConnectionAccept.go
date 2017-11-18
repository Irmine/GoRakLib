package protocol

import "goraklib/protocol/identifiers"

type ConnectionAccept struct {
	*Packet
	ClientAddress string
	ClientPort uint16
	PingSendTime uint64
	PongSendTime uint64
	SystemAddresses []string
	SystemPorts []uint16
	SystemIdVersions []byte
}

func NewConnectionAccept() *ConnectionAccept {
	return &ConnectionAccept{NewPacket(
		identifiers.ConnectionAccept,
	), "", 0, 0, 0, []string{"127.0.0.1"}, []uint16{0}, []byte{4}}
}

func (request *ConnectionAccept) Encode() {
	request.EncodeId()
	request.PutAddress(request.ClientAddress, request.ClientPort, 4)
	request.PutShort(0)

	for i := 0; i < 20; i++ {
		if i < len(request.SystemAddresses) {
			request.PutAddress(request.SystemAddresses[i], request.SystemPorts[i], request.SystemIdVersions[i])
		} else {
			request.PutAddress("0.0.0.0", 0, 4)
		}
	}

	request.PutUnsignedLong(request.PingSendTime)
	request.PutUnsignedLong(request.PongSendTime)
}

func (request *ConnectionAccept) Decode() {

}
