package protocol

type UnconnectedPong struct {
	*UnconnectedMessage
	PingTime        int64
	ServerId		int64
	PongData		string
}

func NewUnconnectedPong() *UnconnectedPong {
	return &UnconnectedPong{NewUnconnectedMessage(NewPacket(
		IdUnconnectedPong,
	)), 0, 0, ""}
}

func (pong *UnconnectedPong) Encode() {
	pong.EncodeId()
	pong.PutLong(pong.PingTime)
	pong.PutLong(pong.ServerId)
	pong.PutMagic()
	pong.PutString(pong.PongData)
}

func (pong *UnconnectedPong) Decode() {
	pong.DecodeStep()
	pong.PingTime = pong.GetLong()
	pong.ServerId = pong.GetLong()
	pong.ReadMagic()
	pong.PongData = pong.GetString()
}
