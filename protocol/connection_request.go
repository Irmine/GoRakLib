package protocol

type ConnectionRequest struct {
	*Packet
	ClientId     uint64
	PingSendTime uint64
	Security     byte
}

func NewConnectionRequest() *ConnectionRequest {
	return &ConnectionRequest{NewPacket(
		IdConnectionRequest,
	), 0, 0, 0}
}

func (request *ConnectionRequest) Encode() {
	request.EncodeId()
	request.PutUnsignedLong(request.ClientId)
	request.PutUnsignedLong(request.PingSendTime)
}

func (request *ConnectionRequest) Decode() {
	request.DecodeStep()
	request.ClientId = request.GetUnsignedLong()
	request.PingSendTime = request.GetUnsignedLong()
}
