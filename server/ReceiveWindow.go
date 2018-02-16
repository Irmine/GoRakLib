package server

import (
	"sync"

	"github.com/irmine/goraklib/protocol"
)

type ReceiveWindow struct {
	session *Session

	lowestIndex  uint32
	highestIndex uint32
	datagrams    sync.Map

	canFlush bool
}

func NewReceiveWindow(session *Session) *ReceiveWindow {
	return &ReceiveWindow{session, 0, 0, sync.Map{}, false}
}

func (window *ReceiveWindow) SetHighestIndex(index uint32) {
	window.highestIndex = index
}

func (window *ReceiveWindow) SetLowestIndex(index uint32) {
	window.lowestIndex = index
}

func (window *ReceiveWindow) SubmitDatagram(datagram *protocol.Datagram) {
	defer func() {
		recover()
	}()
	window.datagrams.Store(datagram.SequenceNumber, datagram)
	if datagram.SequenceNumber > window.highestIndex {
		window.highestIndex = datagram.SequenceNumber
	}

	for i := window.lowestIndex + 1; i <= window.highestIndex; i++ {
		if _, ok := window.datagrams.Load(uint32(i)); !ok {
			break
		}
		if i == window.highestIndex {
			window.canFlush = true
		}
	}
}

func (window *ReceiveWindow) Tick() {
	if window.canFlush {
		window.Flush()
		window.canFlush = false
		return
	}
	var length = 0
	window.datagrams.Range(func(key, value interface{}) bool {
		length++
		return true
	})
	if length > 0 {
		var nak = protocol.NewNAK()
		for i := window.lowestIndex + 1; i < window.highestIndex; i++ {
			if _, ok := window.datagrams.Load(uint32(i)); ok {
				continue
			}
			nak.Packets = append(nak.Packets, i)
		}
		window.session.SendUnconnectedPacket(nak)
	}
}

func (window *ReceiveWindow) Flush() {
	window.Release()
	window.datagrams = sync.Map{}
}

func (window *ReceiveWindow) Release() {
	for i := window.lowestIndex; i <= window.highestIndex; i++ {
		datagram, _ := window.datagrams.Load(uint32(i))
		if datagram == nil {
			continue
		}
		window.session.manager.HandleDatagram(datagram.(*protocol.Datagram), window.session)
	}

	for i := window.lowestIndex + 1; i <= window.highestIndex; i++ {
		window.lowestIndex++
	}
}
