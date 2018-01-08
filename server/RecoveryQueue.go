package server

import (
	"goraklib/protocol"
)

type RecoveryQueue struct {
	recoveryMap map[uint32]*protocol.Datagram
}

func NewRecoveryQueue() *RecoveryQueue {
	return &RecoveryQueue{make(map[uint32]*protocol.Datagram)}
}

func (queue *RecoveryQueue) AddRecoveryFor(datagram *protocol.Datagram) {
	queue.recoveryMap[datagram.SequenceNumber] = datagram
}

func (queue *RecoveryQueue) CanBeRecovered(sequenceNumber uint32) bool {
	var _, canBeRecovered = queue.recoveryMap[sequenceNumber]
	return canBeRecovered
}

func (queue *RecoveryQueue) Recover(sequenceNumbers []uint32) []*protocol.Datagram {
	var datagrams []*protocol.Datagram
	for _, sequenceNum := range sequenceNumbers {
		if queue.CanBeRecovered(sequenceNum) {
			datagram, _ := queue.recoveryMap[sequenceNum]
			datagrams = append(datagrams, datagram)
		}
	}
	return datagrams
}

func (queue *RecoveryQueue) FlagForDeletion(sequenceNumbers []uint32) {
	for _, sequenceNum := range sequenceNumbers {
		delete(queue.recoveryMap, sequenceNum)
	}
}

func (queue *RecoveryQueue) IsClear() bool {
	return len(queue.recoveryMap) == 0
}