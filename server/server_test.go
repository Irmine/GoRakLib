package server

import (
	"testing"
	"strconv"
	"fmt"
	"time"
	"encoding/hex"
)

func Test(t *testing.T) {
	manager := NewManager()
	manager.Start("0.0.0.0", 19132)
	manager.PongData = "MCPE;This is the MOTD;201;1.2.10.2;0;20;" + strconv.Itoa(int(manager.ServerId)) + ";GoMine;Creative;"

	manager.PacketFunction = func(packet []byte, session *Session) {
		fmt.Println("Packet:", hex.EncodeToString(packet[0:1]))
	}
	manager.ConnectFunction = func(session *Session) {
		fmt.Println(session, "connected!")
	}
	manager.DisconnectFunction = func(session *Session) {
		fmt.Println(session, "disconnected!")
	}

	time.Sleep(time.Minute)
	manager.Stop()
}