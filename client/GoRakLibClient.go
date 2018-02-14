package client

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/irmine/goraklib/protocol"
)

type GoRakLibClient struct {
	connection        net.Conn
	connectionAddress string
	connectionPort    uint16

	packetPool *PacketPool

	recoveryQueue *RecoveryQueue
	queue         *PriorityQueue

	isConnected          bool
	isDiscoveringMtuSize bool
	ticker               *time.Ticker

	messageIndex          uint32
	currentSequenceNumber uint32
	orderIndex            map[byte]uint32

	mtuSize int16
	splitId int16

	splits map[int]chan *protocol.EncapsulatedPacket

	online bool
}

func NewGoRakLibClient() *GoRakLibClient {
	var client = &GoRakLibClient{packetPool: NewPacketPool(), recoveryQueue: NewRecoveryQueue(), ticker: time.NewTicker(time.Second), mtuSize: 1500}
	client.queue = NewPriorityQueue(client)

	go func() {
		for {
			client.Tick()
		}
	}()

	go func() {
		var ticker = time.NewTicker(time.Second / 20)
		for range ticker.C {
			client.queue.Flush()
		}
	}()

	go func() {
		for !client.HasConnection() {
			time.Sleep(time.Second / 20)
		}
		if !client.isConnected {
			for range client.ticker.C {
				if client.isConnected {
					break
				}

				ping := protocol.NewUnconnectedPing()
				ping.PingTime = int64(time.Now().Second()) * 1000
				client.SendPacket(ping)
			}
		}
		if client.isDiscoveringMtuSize {
			for range client.ticker.C {
				if !client.isDiscoveringMtuSize {
					break
				}

				request := protocol.NewOpenConnectionRequest1()
				request.Protocol = 6
				request.MtuSize = client.mtuSize
				client.mtuSize--

				client.SendPacket(request)
			}
		}
	}()

	return client
}

func (client *GoRakLibClient) Tick() {
	if !client.HasConnection() {
		return
	}

	go func() {
		client.queue.Flush()
	}()

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
}

func (client *GoRakLibClient) WriteToUDP(buffer []byte) error {
	var _, err = client.connection.Write(buffer)
	return err
}

func (client *GoRakLibClient) ReadFromUDP() (protocol.IPacket, error) {
	var buffer = make([]byte, 4096)

	n, err := client.connection.Read(buffer)
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

	realAddresses, err := net.LookupIP(address)
	if err == nil {
		address = realAddresses[0].String()
	}

	fmt.Println("Client connecting to address:", address, "with port:", port, "...")
	client.connection, err = net.Dial("udp", address+":"+strconv.Itoa(int(port)))
	client.connectionAddress = address
	client.connectionPort = port

	if err != nil {
		return err
	}

	return nil
}

func (client *GoRakLibClient) Disconnect() {
	client.connection.Close()
	client.connection = nil
}

func (client *GoRakLibClient) HasConnection() bool {
	return client.connection != nil
}

func (client *GoRakLibClient) SendPacket(packet protocol.IPacket) {
	packet.Encode()

	client.WriteToUDP(packet.GetBuffer())
}

func (client *GoRakLibClient) SendConnectedPacket(packet protocol.IConnectedPacket, reliability byte, priority byte) {
	packet.Encode()

	var encapsulatedPacket = protocol.NewEncapsulatedPacket()
	encapsulatedPacket.Buffer = packet.GetBuffer()
	encapsulatedPacket.OrderChannel = 0
	encapsulatedPacket.Reliability = reliability

	if priority != PriorityImmediate {
		client.queue.AddEncapsulatedToQueue(encapsulatedPacket, priority)
		return
	}

	if encapsulatedPacket.IsReliable() {
		encapsulatedPacket.MessageIndex = client.messageIndex
		client.messageIndex++
	}
	if encapsulatedPacket.IsSequenced() {
		encapsulatedPacket.OrderIndex = client.orderIndex[encapsulatedPacket.OrderChannel]
		client.orderIndex[encapsulatedPacket.OrderChannel]++
	}

	var datagram = protocol.NewDatagram()
	datagram.NeedsBAndAs = true

	datagram.SequenceNumber = client.currentSequenceNumber
	client.currentSequenceNumber++

	datagram.AddPacket(encapsulatedPacket)
	client.SendPacket(datagram)

	client.recoveryQueue.AddRecoveryFor(datagram)
}

func (client *GoRakLibClient) IsOnline() bool {
	return client.online
}
