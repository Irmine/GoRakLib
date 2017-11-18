package identifiers

const (
	UnconnectedPing = 0x01
	UnconnectedPong = 0x1c

	OpenConnectionRequest1 = 0x05
	OpenConnectionResponse1 = 0x06

	OpenConnectionRequest2 = 0x07
	OpenConnectionResponse2 = 0x08

	ConnectionRequest = 0x09
	ConnectionAccept = 0x10

	NewIncomingConnection = 0x13
)