package protocol

const (
	FlagDatagramAck  = 0xc0
	FlagDatagramNack = 0xa0
)

const (
	IdConnectedPing = 0x00
	IdConnectedPong = 0x03

	IdUnconnectedPing = 0x01
	IdUnconnectedPong = 0x1c

	IdOpenConnectionRequest1  = 0x05
	IdOpenConnectionReply1 = 0x06
	IdOpenConnectionRequest2  = 0x07
	IdOpenConnectionReply2 = 0x08

	IdConnectionRequest = 0x09
	IdConnectionAccept  = 0x10

	IdNewIncomingConnection = 0x13

	IdDisconnectNotification = 0x15
)
