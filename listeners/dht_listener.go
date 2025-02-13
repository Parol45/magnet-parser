package main

import (
	"fmt"
	"log/slog"
	"magnet-parser/bencode"
	"magnet-parser/requests"
	"magnet-parser/utils"
	"net"
	"strings"
	"time"
)

const ServerPort = 14888
var udpServer, _ = net.ListenPacket("udp", fmt.Sprintf(":%d", ServerPort))

func listenUDP() {
	defer udpServer.Close()
	for {
		buf := make([]byte, 65507)
		n, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			continue
		}
		go processInputRequest(udpServer, addr, buf[:n])
	}
}

// todo ответы на разные пакеты (сохранение токенов и сидов)
func processInputRequest(udpServer net.PacketConn, addr net.Addr, buf []byte) {
	if buf[0] == 0 {
		slog.Error("Unexpected response: " + string(buf))
		return
	}
	slog.Info(fmt.Sprintf("Encoded json is: %s\n", strings.ReplaceAll(string(buf), "\r", "")))
	json, err := bencode.BencodeToJSON(buf)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while converting bencode to json: %v\n", err))
	} else {
		slog.Info(fmt.Sprintf("Decoded json is: %s\n", json))
	}
	var bytes []byte
	bytes, err = bencode.JSONToBencode(json)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while converting bencode to json: %v\n", err))
	} else {
		slog.Info(fmt.Sprintf("Encoded json is: %s\n", strings.ReplaceAll(string(bytes), "\r", "")))
	}
	json, _ = bencode.BencodeToJSON(bytes)
	slog.Info(fmt.Sprintf("Decoded json is: %s\n", json))
	bytes, err = bencode.JSONToBencode(json)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while converting bencode to json: %v\n", err))
	} else {
		slog.Info(fmt.Sprintf("Encoded json is: %s\n", strings.ReplaceAll(string(bytes), "\r", "")))
	}
}

func main() {
	utils.SetupLogger("dht_listener")
	addr := "93.158.213.92:1337" //todo dns-resolve
	go listenUDP()
	udpServer.WriteTo(requests.Ping(), utils.StringToUDPAddr(addr))
	//udpServer.WriteTo(requests.GetPeers("0f52fb0190dcd61381cdd0893173b1309ba9868f"), utils.StringToUDPAddr(addr))
	//udpServer.WriteTo(requests.FindNode("0f43110b4bdf604508cb75dfd326dcd68ac936c7"), utils.StringToUDPAddr(addr))
	//udpServer.WriteTo(requests.AnnouncePeer("0f52fb0190dcd61381cdd0893173b1309ba9868f", ServerPort, "1dcbde3b"), utils.StringToUDPAddr(addr))
	time.Sleep(10*time.Minute)
}
