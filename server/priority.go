package server

import (
	"github.com/irmine/goraklib/protocol"
	"math"
)

const (
	// PriorityImmediate is the highest possible priority.
	// Packets with this priority get sent out immediately.
	PriorityImmediate Priority = iota
	// PriorityHigh is the highest possible priority that gets buffered.
	// High priority packets get sent out every tick.
	PriorityHigh
	// PriorityMedium is the priority most used.
	// Medium priority packets get sent out every other tick.
	PriorityMedium
	// PriorityLow is the lowest possible priority.
	// Low priority packets get sent out every fourth tick.
	PriorityLow
)

// Priority is the priority at which packets will be sent out when queued.
// PriorityImmediate will make packets get sent out immediately.
type Priority byte

// A PriorityQueue is used to send packets with a certain priority.
// Encapsulated packets can be queued in these queues.
type PriorityQueue chan *protocol.EncapsulatedPacket

// NewPriorityQueue returns a new priority queue with buffering size.
// The buffering size specifies the maximum amount of packets in the queue,
// and the queue becomes blocking once the amount of packets exceeds the buffering size.
func NewPriorityQueue(bufferingSize int) *PriorityQueue {
	queue := PriorityQueue(make(chan *protocol.EncapsulatedPacket, bufferingSize))
	return &queue
}

// AddEncapsulated adds an encapsulated packet to a priority queue.
// The packet will first be split into smaller sub packets if needed,
// after which all packets will be added to the queue.
func (queue *PriorityQueue) AddEncapsulated(packet *protocol.EncapsulatedPacket, session *Session) {
	for _, encapsulated := range queue.Split(packet, session) {
		*queue <- encapsulated
	}
}

// Flush flushes all encapsulated packets in the priority queue, and sends them to a session.
// All encapsulated packets will first be fetched from the channel,
// after which they will be put into datagrams.
// A new datagram is made once an encapsulated packet makes the size
// of a datagram exceed the MTU size of the session.
func (queue *PriorityQueue) Flush(session *Session) {
	if len(*queue) == 0 {
		return
	}
	ind := 0
	datagram := protocol.NewDatagram()
	datagram.NeedsBAndAs = true
	datagrams := map[int]*protocol.Datagram{0: datagram}
	datagram.SequenceNumber = session.Indexes.sendSequence
	session.Indexes.sendSequence++

	i := 0
	for len(*queue) > 0 && i < 16 {
		i++
		packet := <-*queue
		if datagrams[ind].GetLength()+packet.GetLength() > int(session.MTUSize-38) {
			ind++
			datagrams[ind] = protocol.NewDatagram()
			datagrams[ind].NeedsBAndAs = true
			datagrams[ind].SequenceNumber = session.Indexes.sendSequence
			session.Indexes.sendSequence++
		}
		datagrams[ind].AddPacket(packet)
	}

	l := len(datagrams)
	for j := 0; j < l; j++ {
		datagram := datagrams[j]
		datagram.Encode()
		session.RecoveryQueue.AddRecovery(datagram)
		session.Manager.Server.Write(datagram.Buffer, session.UDPAddr)
	}
}

// Split splits an encapsulated packet into smaller sub packets.
// Every encapsulated packet that exceeds the MTUSize of the session
// will be split into sub packets, and returned into a slice.
func (queue *PriorityQueue) Split(packet *protocol.EncapsulatedPacket, session *Session) []*protocol.EncapsulatedPacket {
	mtuSize := int(session.MTUSize - 60) // We subtract 60 to account for headers.
	var packets []*protocol.EncapsulatedPacket
	if packet.IsSequenced() {
		packet.OrderIndex = session.Indexes.orderIndex
		session.Indexes.orderIndex++
	}

	if packet.GetLength() > mtuSize {
		buffer := packet.GetBuffer()
		splitSize := mtuSize
		var b uint
		for i := 0; i < len(buffer)+splitSize; i += splitSize {
			if i + splitSize >= len(buffer) {
				splitSize = len(buffer) - i
				if splitSize == 0 {
					break
				}
			}
			split := buffer[i:i+splitSize]
			encapsulated := protocol.NewEncapsulatedPacket()
			encapsulated.HasSplit = true
			encapsulated.SplitId = session.Indexes.splitId
			encapsulated.SplitIndex = b
			encapsulated.SplitCount = uint(math.Ceil(float64(len(buffer)) / float64(splitSize)))
			encapsulated.Reliability = packet.Reliability
			b++
			encapsulated.Buffer = split
			encapsulated.OrderIndex = packet.OrderIndex
			if packet.IsReliable() {
				encapsulated.MessageIndex = session.Indexes.messageIndex
				session.Indexes.messageIndex++
			}
			packets = append(packets, encapsulated)
		}
		session.Indexes.splitId++
	} else {
		if packet.IsReliable() {
			packet.MessageIndex = session.Indexes.messageIndex
			session.Indexes.messageIndex++
		}
		packets = append(packets, packet)
	}
	return packets
}