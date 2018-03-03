package protocol

type ConnectedPing struct {
	*Packet
	PingSendTime int64
}

func NewConnectedPing() *ConnectedPing {
	return &ConnectedPing{NewPacket(
		IdConnectedPing,
	), 0}
}

func (ping *ConnectedPing) Encode() {
	ping.EncodeId()
	ping.PutLong(ping.PingSendTime)
}

func (ping *ConnectedPing) Decode() {
	ping.DecodeStep()
	ping.PingSendTime = ping.GetLong()
}
