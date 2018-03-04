package server

import (
	"github.com/irmine/goraklib/protocol"
	"time"
)

// TimestampedDatagram is a datagram encapsulated by a timestamp.
// Every datagram added to the receive window gets its timestamp recorded immediately.
type TimestampedDatagram struct {
	*protocol.Datagram
	Timestamp int64
}

// ReceiveWindow is a window used to hold datagrams until they're read to be released.
// ReceiveWindow restores the order of datagrams that arrived out of order,
// and sends NAKs where needed.
type ReceiveWindow struct {
	// DatagramHandleFunction is a function that gets called once a datagram gets released from the receive window.
	// A timestamped datagram gets returned with the timestamp of the time the datagram entered the receive window.
	DatagramHandleFunction func(datagram TimestampedDatagram)

	pendingDatagrams       chan TimestampedDatagram
	datagrams              map[uint32]TimestampedDatagram
	expectedSequenceNumber uint32
	highestSequenceNumber  uint32
}

// NewReceiveWindow returns a new receive window.
func NewReceiveWindow() *ReceiveWindow {
	return &ReceiveWindow{func(datagram TimestampedDatagram){}, make(chan TimestampedDatagram, 128), make(map[uint32]TimestampedDatagram), 0, 0}
}

// AddDatagram adds a datagram to the receive window.
// The datagram is first encapsulated with a timestamp,
// and is added to a channel in order to await the next tick for further processing.
func (window *ReceiveWindow) AddDatagram(datagram *protocol.Datagram) {
	if datagram.SequenceNumber < window.expectedSequenceNumber {
		return
	}
	if datagram.SequenceNumber > window.highestSequenceNumber {
		window.highestSequenceNumber = datagram.SequenceNumber
	}
	window.pendingDatagrams <- TimestampedDatagram{datagram, time.Now().Unix()}
}

// Tick ticks the ReceiveWindow and releases any datagrams when possible.
// Tick also fetches all datagrams that are currently in the channel.
func (window *ReceiveWindow) Tick() {
	for len(window.pendingDatagrams) > 0 {
		var datagram = <-window.pendingDatagrams
		window.datagrams[datagram.SequenceNumber] = datagram
	}
	for i := window.expectedSequenceNumber;; i++ {
		if datagram, ok := window.datagrams[i]; ok {
			window.DatagramHandleFunction(datagram)
			window.expectedSequenceNumber++
			delete(window.datagrams, i)
		} else {
			break
		}
	}
}