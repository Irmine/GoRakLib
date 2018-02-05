package server

import (
	"goraklib/protocol"
	"sync"
)

type RecoveryQueue struct {
	recoveryMap sync.Map
}

func NewRecoveryQueue() *RecoveryQueue {
	return &RecoveryQueue{sync.Map{}}
}

func (queue *RecoveryQueue) AddRecoveryFor(datagram *protocol.Datagram) {
	queue.recoveryMap.Store(datagram.SequenceNumber, datagram)
}

func (queue *RecoveryQueue) CanBeRecovered(sequenceNumber uint32) bool {
	var _, canBeRecovered = queue.recoveryMap.Load(sequenceNumber)
	return canBeRecovered
}

func (queue *RecoveryQueue) Recover(sequenceNumbers []uint32) []*protocol.Datagram {
	var datagrams []*protocol.Datagram
	for _, sequenceNum := range sequenceNumbers {
		if queue.CanBeRecovered(sequenceNum) {
			datagram, _ := queue.recoveryMap.Load(sequenceNum)
			datagrams = append(datagrams, datagram.(*protocol.Datagram))
		}
	}
	return datagrams
}

func (queue *RecoveryQueue) FlagForDeletion(sequenceNumbers []uint32) {
	for _, sequenceNum := range sequenceNumbers {
		queue.recoveryMap.Delete(sequenceNum)
	}
}

func (queue *RecoveryQueue) IsClear() bool {
	var length = 0
	queue.recoveryMap.Range(func(key, value interface{}) bool {
		length++
		return true
	})
	return length == 0
}