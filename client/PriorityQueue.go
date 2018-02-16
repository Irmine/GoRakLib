package client

import (
	"math"

	"github.com/irmine/goraklib/protocol"
)

const (
	PriorityImmediate = 0
	PriorityLow       = 1
	PriorityMedium    = 2
	PriorityHigh      = 3
)

type PriorityQueue struct {
	client *GoRakLibClient
	Low    chan *protocol.EncapsulatedPacket
	Medium chan *protocol.EncapsulatedPacket
	High   chan *protocol.EncapsulatedPacket
}

func NewPriorityQueue(client *GoRakLibClient) *PriorityQueue {
	return &PriorityQueue{client, make(chan *protocol.EncapsulatedPacket, 16), make(chan *protocol.EncapsulatedPacket, 16), make(chan *protocol.EncapsulatedPacket, 16)}
}

func (queue *PriorityQueue) Wipe() {
	queue.Low = make(chan *protocol.EncapsulatedPacket, 16)
	queue.Medium = make(chan *protocol.EncapsulatedPacket, 16)
	queue.High = make(chan *protocol.EncapsulatedPacket, 16)
}

func (queue *PriorityQueue) AddEncapsulatedToQueue(packet *protocol.EncapsulatedPacket, priority byte) {
	if packet.IsReliable() {
		packet.MessageIndex = queue.client.messageIndex
		queue.client.messageIndex++
	}
	if packet.IsSequenced() {
		packet.OrderIndex = queue.client.orderIndex[packet.OrderChannel]
		queue.client.orderIndex[packet.OrderChannel]++
	}

	var maximumEncapsulatedSize = int(queue.client.mtuSize - 60)

	if packet.GetLength() > maximumEncapsulatedSize {
		var buffer = packet.GetBuffer()
		var splitBuffers [][]byte
		var split []byte

		for len(buffer) >= maximumEncapsulatedSize {
			split, buffer = buffer[:maximumEncapsulatedSize], buffer[maximumEncapsulatedSize:]
			splitBuffers = append(splitBuffers, split)
		}
		splitBuffers = append(splitBuffers, buffer)

		var splitId = queue.client.splitId % math.MaxInt16
		queue.client.splitId++

		for index, splitBuffer := range splitBuffers {
			encapsulated := protocol.NewEncapsulatedPacket()
			encapsulated.ResetStream()

			encapsulated.HasSplit = true
			encapsulated.SplitId = splitId
			encapsulated.SplitIndex = uint(index)
			encapsulated.SplitCount = uint(len(splitBuffers))

			encapsulated.Reliability = packet.Reliability
			encapsulated.MessageIndex = packet.MessageIndex + uint32(index)
			if index != 0 {
				queue.client.messageIndex++
			}

			encapsulated.Buffer = splitBuffer

			encapsulated.OrderChannel = packet.OrderChannel
			encapsulated.OrderIndex = packet.OrderIndex

			queue.AddToQueue(encapsulated, priority)
		}
		return
	}

	queue.AddToQueue(packet, priority)
}

func (queue *PriorityQueue) AddToQueue(packet *protocol.EncapsulatedPacket, priority byte) {
	switch priority {
	case PriorityLow:
		queue.Low <- packet
	case PriorityMedium:
		queue.Medium <- packet
	case PriorityHigh:
		queue.High <- packet
	}
}

func (queue *PriorityQueue) Flush() {
	queue.FlushHighPriority()
	queue.FlushMediumPriority()
	queue.FlushHighPriority()
}

func (queue *PriorityQueue) FlushHighPriority() {
	if len(queue.High) == 0 {
		return
	}

	var datagramIndex = 0
	var datagrams = map[int]*protocol.Datagram{datagramIndex: protocol.NewDatagram()}
	datagrams[datagramIndex].NeedsBAndAs = true

	datagrams[datagramIndex].SequenceNumber = queue.client.currentSequenceNumber
	queue.client.currentSequenceNumber++

	for len(queue.High) > 0 {
		var encapsulated = <-queue.High
		if datagrams[datagramIndex].GetLength()+encapsulated.GetLength() > int(queue.client.mtuSize-36) {
			datagramIndex++

			var newDatagram = protocol.NewDatagram()
			newDatagram.SequenceNumber = queue.client.currentSequenceNumber
			queue.client.currentSequenceNumber++

			newDatagram.NeedsBAndAs = true
			datagrams[datagramIndex] = newDatagram
		}
		datagrams[datagramIndex].AddPacket(encapsulated)
	}

	for _, datagram := range datagrams {
		if len(*datagram.GetPackets()) == 0 {
			break
		}
		queue.client.SendPacket(datagram)

		queue.client.recoveryQueue.AddRecoveryFor(datagram)
	}
}

func (queue *PriorityQueue) FlushMediumPriority() {
	if len(queue.Medium) == 0 {
		return
	}

	var datagramIndex = 0
	var datagrams = map[int]*protocol.Datagram{datagramIndex: protocol.NewDatagram()}
	datagrams[datagramIndex].NeedsBAndAs = true

	datagrams[datagramIndex].SequenceNumber = queue.client.currentSequenceNumber
	queue.client.currentSequenceNumber++

	for len(queue.Medium) > 0 {
		var encapsulated = <-queue.Medium
		if datagrams[datagramIndex].GetLength()+encapsulated.GetLength() > int(queue.client.mtuSize-36) {
			datagramIndex++

			var newDatagram = protocol.NewDatagram()
			newDatagram.SequenceNumber = queue.client.currentSequenceNumber
			queue.client.currentSequenceNumber++

			newDatagram.NeedsBAndAs = true
			datagrams[datagramIndex] = newDatagram
		}
		datagrams[datagramIndex].AddPacket(encapsulated)
	}

	for _, datagram := range datagrams {
		if len(*datagram.GetPackets()) == 0 {
			break
		}
		queue.client.SendPacket(datagram)
		queue.client.recoveryQueue.AddRecoveryFor(datagram)
	}
}

func (queue *PriorityQueue) FlushLowPriority() {
	if len(queue.Low) == 0 {
		return
	}

	var datagramIndex = 0
	var datagrams = map[int]*protocol.Datagram{datagramIndex: protocol.NewDatagram()}
	datagrams[datagramIndex].NeedsBAndAs = true

	datagrams[datagramIndex].SequenceNumber = queue.client.currentSequenceNumber
	queue.client.currentSequenceNumber++

	for len(queue.Low) > 0 {
		var encapsulated = <-queue.Low
		if datagrams[datagramIndex].GetLength()+encapsulated.GetLength() > int(queue.client.mtuSize-36) {
			datagramIndex++

			var newDatagram = protocol.NewDatagram()
			newDatagram.SequenceNumber = queue.client.currentSequenceNumber
			queue.client.currentSequenceNumber++

			newDatagram.NeedsBAndAs = true
			datagrams[datagramIndex] = newDatagram
		}
		datagrams[datagramIndex].AddPacket(encapsulated)
	}

	for _, datagram := range datagrams {
		if len(*datagram.GetPackets()) == 0 {
			break
		}
		queue.client.SendPacket(datagram)
		queue.client.recoveryQueue.AddRecoveryFor(datagram)
	}
}
