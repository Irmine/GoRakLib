package client

import (
	"net"
	"strconv"
	"goraklib/protocol"
	"errors"
	"time"
)

type GoRakLibClient struct {
	localAddress *net.UDPAddr
	serverAddress *net.UDPAddr

	connection *net.UDPConn

	packetPool *PacketPool

	recoveryQueue *RecoveryQueue
	priorityQueue *PriorityQueue

	isConnected bool
	ticker *time.Ticker

	messageIndex uint32
	currentSequenceNumber uint32
	orderIndex map[byte]uint32

	mtuSize int16
	splitId int16

	splits map[int]chan *protocol.EncapsulatedPacket
}

func NewGoRakLibClient(address string, port uint16) *GoRakLibClient {
	var connection, _ = net.ResolveUDPAddr("udp", address + strconv.Itoa(int(port)))
	var client = &GoRakLibClient{localAddress:connection, packetPool:NewPacketPool(), recoveryQueue:NewRecoveryQueue(), ticker:time.NewTicker(time.Second)}
	client.priorityQueue = NewPriorityQueue(client)

	go func() {
		for {
			client.Tick()
		}
	}()

	return client
}

func (client *GoRakLibClient) Tick() {
	if !client.HasConnection() {
		return
	}

	if !client.isConnected {
		for range client.ticker.C {
			ping := protocol.NewUnconnectedPing()
			ping.PingTime = int64(time.Now().Second()) * 1000
			client.SendPacket(ping)
		}
	}

	go func() {
		client.priorityQueue.Flush()
	}()

	go func() {
		var packet, err = client.ReadFromUDP()
		if err != nil {
			return
		}
		packet.Decode()
		if packet.HasMagic() {
			client.HandleUnconnectedMessage(packet)
		} else {
			if datagram, ok := packet.(*protocol.Datagram); ok {
				client.HandleDatagram(datagram)
			} else if nak, ok := packet.(*protocol.NAK); ok {
				client.HandleNak(nak)
			} else if ack, ok := packet.(*protocol.ACK); ok {
				client.HandleAck(ack)
			}
		}
		//var session, _ = server.sessionManager.GetSession(ip, port)
		//go session.Forward(packet)
	}()
}

func (client *GoRakLibClient) WriteToUDP(buffer []byte) error {
	var _, err = client.connection.WriteToUDP(buffer, client.localAddress)
	return err
}

func (client *GoRakLibClient) ReadFromUDP() (protocol.IPacket, error) {
	var buffer = make([]byte, 4096)

	n, _, err := client.connection.ReadFromUDP(buffer)
	if err != nil {
		panic(err)
	}

	buffer = buffer[:n]
	var packet protocol.IPacket

	if n == 0 {
		return packet, errors.New("received null packet")
	}

	packet = client.packetPool.GetPacket(buffer)
	packet.SetBuffer(buffer)

	return packet, nil
}

func (client *GoRakLibClient) Connect(address string, port uint16) error {
	var err error
	client.serverAddress, err = net.ResolveUDPAddr("udp", address + strconv.Itoa(int(port)))
	if err != nil {
		return err
	}

	client.connection, err = net.DialUDP("udp", client.localAddress, client.serverAddress)

	return nil
}

func (client *GoRakLibClient) Disconnect() {
	client.connection.Close()
	client.connection = nil
	client.serverAddress = nil
}

func (client *GoRakLibClient) HasConnection() bool {
	return client.connection != nil && client.serverAddress != nil
}

func (client *GoRakLibClient) SendPacket(packet protocol.IPacket) {
	packet.Encode()

	client.WriteToUDP(packet.GetBuffer())
}