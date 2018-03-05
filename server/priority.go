package server

import (
	"github.com/irmine/goraklib/protocol"
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
	datagramIndex := 0
	datagram := protocol.NewDatagram()
	datagram.NeedsBAndAs = true
	datagrams := map[int]*protocol.Datagram{0: datagram}
	for len(*queue) > 0 {
		packet := <-*queue
		if datagrams[datagramIndex].GetLength()+packet.GetLength() > int(session.MTUSize-36) {
			datagrams[datagramIndex].Encode()
			datagramIndex++
			datagrams[datagramIndex] = protocol.NewDatagram()
			datagrams[datagramIndex].NeedsBAndAs = true
			datagrams[datagramIndex].SequenceNumber = session.Indexes.sendSequence
			session.Indexes.sendSequence++
		}
		datagrams[datagramIndex].AddPacket(packet)
	}
	for _, datagram := range datagrams {
		datagram.Encode()
		session.RecoveryQueue.AddRecovery(datagram)
		session.Manager.Server.Write(datagram.Buffer, session.UDPAddr)
	}
}

// Split splits an encapsulated packet into smaller sub packets.
// Every encapsulated packet that exceeds the MTUSize of the session
// will be split into sub packets, and returned into a slice.
func (queue *PriorityQueue) Split(packet *protocol.EncapsulatedPacket, session *Session) []*protocol.EncapsulatedPacket {
	mtuSize := session.MTUSize - 60 // We subtract 60 to account for the headers of datagrams and encapsulated packets.
	var packets []*protocol.EncapsulatedPacket
	if int16(packet.GetLength()) > mtuSize {
		buffer := packet.GetBuffer()
		var splitBuffers [][]byte
		var split []byte
		for int16(len(buffer)) >= mtuSize {
			split, buffer = buffer[:mtuSize], buffer[mtuSize:]
			splitBuffers = append(splitBuffers, split)
		}
		splitBuffers = append(splitBuffers, buffer)
		for index, splitBuffer := range splitBuffers {
			encapsulated := protocol.NewEncapsulatedPacket()
			encapsulated.HasSplit = true
			encapsulated.SplitId = int16(index)
			encapsulated.SplitIndex = uint(index)
			encapsulated.SplitCount = uint(len(splitBuffers))
			encapsulated.Reliability = packet.Reliability
			encapsulated.MessageIndex = packet.MessageIndex + uint32(index)
			encapsulated.Buffer = splitBuffer
			encapsulated.OrderChannel = packet.OrderChannel
			encapsulated.OrderIndex = packet.OrderIndex

			packets = append(packets, encapsulated)
		}
	} else {
		packets = append(packets, packet)
	}
	return packets
}