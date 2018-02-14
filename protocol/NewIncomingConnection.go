package protocol

import "github.com/irmine/goraklib/protocol/identifiers"

type NewIncomingConnection struct {
	*Packet
	ServerAddress    string
	ServerPort       uint16
	PingSendTime     uint64
	PongSendTime     uint64
	SystemAddresses  []string
	SystemPorts      []uint16
	SystemIdVersions []byte
}

func NewNewIncomingConnection() *NewIncomingConnection {
	return &NewIncomingConnection{NewPacket(
		identifiers.NewIncomingConnection,
	), "", 0, 0, 0, []string{"127.0.0.1"}, []uint16{0}, []byte{4}}
}

func (request *NewIncomingConnection) Encode() {
	request.EncodeId()
	request.PutAddress(request.ServerAddress, request.ServerPort, 4)

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
