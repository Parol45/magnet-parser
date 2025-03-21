package main

import (
	"magnet-parser/globals"
	"magnet-parser/listeners"
	"magnet-parser/requests"
	"time"
)

func main() {
	globals.SetupLogger("dht_listener")
	addr := "130.239.18.158:8603" //todo dns-resolve
	go listeners.ListenUDP()
	//listeners.UdpServer.WriteTo(requests.Ping(), globals.StringToUDPAddr(addr))
	listeners.UdpServer.WriteTo(requests.GetPeers("a84c21f7d1208ef74a1acc3fbe3704b102d46590"), globals.StringToUDPAddr(addr))
	//udpServer.WriteTo(requests.FindNode("0f43110b4bdf604508cb75dfd326dcd68ac936c7"), utils.StringToUDPAddr(addr))
	//udpServer.WriteTo(requests.AnnouncePeer("0f52fb0190dcd61381cdd0893173b1309ba9868f", ServerPort, "1dcbde3b"), utils.StringToUDPAddr(addr))
	time.Sleep(1*time.Minute)
}
