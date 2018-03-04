package server

import (
	"sync"
	"github.com/irmine/goraklib/protocol"
)

// A RecoveryQueue manages the recovery of lost datagrams over the connection.
// Datagrams get restored by the client sending a NAK.
// The recovery queue holds every datagram sent and releases them,
// once an ACK is received with the datagram's sequence number.
type RecoveryQueue struct {
	mutex sync.Mutex
	datagrams map[uint32]*protocol.Datagram
}

// NewRecoveryQueue returns a new recovery queue.
func NewRecoveryQueue() *RecoveryQueue {
	return &RecoveryQueue{sync.Mutex{}, make(map[uint32]*protocol.Datagram)}
}

// AddRecovery adds recovery for the given datagram.
// The recovery will consist until an ACK gets sent by the client,
// and the datagram is safe to be removed.
func (queue *RecoveryQueue) AddRecovery(datagram *protocol.Datagram) {
	queue.mutex.Lock()
	queue.datagrams[datagram.SequenceNumber] = datagram
	queue.mutex.Unlock()
}

// IsRecoverable checks if the datagram with the given sequence number is recoverable.
func (queue *RecoveryQueue) IsRecoverable(sequenceNumber uint32) bool {
	queue.mutex.Lock()
	_, ok := queue.datagrams[sequenceNumber]
	queue.mutex.Unlock()
	return ok
}

// RemoveRecovery removes recovery for all sequence numbers given.
// Removed datagrams can not be retrieved in anyway,
// therefore this function should only be used once the client sends an ACK to ensure arrival.
func (queue *RecoveryQueue) RemoveRecovery(sequenceNumbers []uint32) {
	queue.mutex.Lock()
	for _, sequenceNumber := range sequenceNumbers {
		delete(queue.datagrams, sequenceNumber)
	}
	queue.mutex.Unlock()
}

// Recover recovers all datagrams associated with the sequence numbers in the array given.
// Every recoverable datagram with sequence number in the array will be returned,
// along with an array containing all recovered sequence numbers.
func (queue *RecoveryQueue) Recover(sequenceNumbers []uint32) ([]*protocol.Datagram, []uint32) {
	var datagrams []*protocol.Datagram
	var recoveredSequenceNumbers []uint32
	for _, sequenceNumber := range sequenceNumbers {
		if datagram, ok := queue.datagrams[sequenceNumber]; ok {
			datagrams = append(datagrams, datagram)
			recoveredSequenceNumbers = append(recoveredSequenceNumbers, sequenceNumber)
		}
	}
	return datagrams, sequenceNumbers
}