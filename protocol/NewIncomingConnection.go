package protocol

import "goraklib/protocol/identifiers"

type NewIncomingConnection struct {
	*Packet
	ServerAddress string
	ServerPort uint16
	PingSendTime uint64
	PongSendTime uint64
	SystemAddresses []string
	SystemPorts []uint16
	SystemIdVersions []byte
}

func NewNewIncomingConnection() *NewIncomingConnection {
	return &NewIncomingConnection{NewPacket(
		identifiers.NewIncomingConnection,
	), "", 0, 0, 0, []string{"127.0.0.1"}, []uint16{0}, []byte{4}}
}

func (request *NewIncomingConnection) Encode() {

}

func (request *NewIncomingConnection) Decode() {
	request.DecodeStep()
	request.ServerAddress, request.ServerPort, _ = request.GetAddress()
	for i := 0; i < 20; i++ {
		address, port, version := request.GetAddress()
		request.SystemAddresses = append(request.SystemAddresses, address)
		request.SystemPorts = append(request.SystemPorts, port)
		request.SystemIdVersions = append(request.SystemIdVersions, version)
	}

	request.PingSendTime = request.GetUnsignedLong()
	request.PongSendTime = request.GetUnsignedLong()
}