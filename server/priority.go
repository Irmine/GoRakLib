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

// AddEncapsulated adds an encapsulated packet to a priority queue.
// The packet will first be split into smaller sub packets if needed,
// after which all packets will be added to the queue.
func (queue PriorityQueue) AddEncapsulated(packet *protocol.EncapsulatedPacket, session *Session) {
	for _, encapsulated := range queue.Split(packet, session) {
		queue <- encapsulated
	}
}

// Split splits an encapsulated packet into smaller sub packets.
// Every encapsulated packet that exceeds the MTUSize of the session
// will be split into sub packets, and returned into a slice.
func (queue PriorityQueue) Split(packet *protocol.EncapsulatedPacket, session *Session) []*protocol.EncapsulatedPacket {
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