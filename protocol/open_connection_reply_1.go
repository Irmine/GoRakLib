package protocol

type OpenConnectionReply1 struct {
	*UnconnectedMessage
	ServerId int64
	MtuSize  int16
	Security bool
}

func NewOpenConnectionReply1() *OpenConnectionReply1 {
	return &OpenConnectionReply1{NewUnconnectedMessage(NewPacket(
		IdOpenConnectionReply1,
	)), 0, 0, false}
}

func (response *OpenConnectionReply1) Encode() {
	response.EncodeId()
	response.PutMagic()
	response.PutLong(response.ServerId)
	response.PutBool(response.Security)
	response.PutShort(response.MtuSize)
}

func (response *OpenConnectionReply1) Decode() {
	response.DecodeStep()
	response.ReadMagic()
	response.ServerId = response.GetLong()
	response.Security = response.GetBool()
	response.MtuSize = response.GetShort()
}
