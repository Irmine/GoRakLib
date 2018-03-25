package test

import (
	"testing"
	"strconv"
	"fmt"
	"time"
	"encoding/hex"
	"github.com/irmine/goraklib/server"
)

func Test(t *testing.T) {
	manager := server.NewManager()
	manager.Start("0.0.0.0", 19132)
	manager.PongData = "MCPE;§e§lTesting §bServer;201;;0;20;" + strconv.Itoa(int(manager.ServerId)) + ";§aVersionless Minecraft Server MOTD;Creative;"

	manager.PacketFunction = func(packet []byte, session *server.Session) {
		fmt.Println("Packet:", hex.EncodeToString(packet[0:1]))
	}
	manager.ConnectFunction = func(session *server.Session) {
		manager.BlockIP(session.UDPAddr, time.Second * 20)
		fmt.Println(session, "connected!")
	}
	manager.DisconnectFunction = func(session *server.Session) {
		fmt.Println(session, "disconnected!")
	}

	time.Sleep(time.Minute)
	manager.Stop()
}
