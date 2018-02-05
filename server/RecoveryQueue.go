package server

import (
	"goraklib/protocol"
	"sync"
)

type RecoveryQueue struct {
	mutex sync.Mutex
	recoveryMap map[uint32]*protocol.Datagram
}

func NewRecoveryQueue() *RecoveryQueue {
	return &RecoveryQueue{sync.Mutex{}, make(map[uint32]*protocol.Datagram)}
}

func (queue *RecoveryQueue) AddRecoveryFor(datagram *protocol.Datagram) {
	queue.mutex.Lock()
	queue.recoveryMap[datagram.SequenceNumber] = datagram
	queue.mutex.Unlock()
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
	queue.mutex.Lock()
	for _, sequenceNum := range sequenceNumbers {
		delete(queue.recoveryMap, sequenceNum)
	}
	queue.mutex.Unlock()
}

func (queue *RecoveryQueue) IsClear() bool {
	return len(queue.recoveryMap) == 0
}