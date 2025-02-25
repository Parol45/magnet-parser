package listeners

import (
	"fmt"
	"log/slog"
	"magnet-parser/bencode_converters/json"
	"magnet-parser/bencode_converters/obj"
	"net"
	"strings"
)

const ServerPort = 14888
var UdpServer, _ = net.ListenPacket("udp", fmt.Sprintf(":%d", ServerPort))

func ListenUDP() {
	defer UdpServer.Close()
	for {
		buf := make([]byte, 65507)
		n, addr, err := UdpServer.ReadFrom(buf)
		if err != nil {
			continue
		}
		go processInputRequest(UdpServer, addr, buf[:n])
	}
}

// todo ответы на разные пакеты (сохранение токенов и сидов)
func processInputRequest(udpServer net.PacketConn, addr net.Addr, buf []byte) {
	if buf[0] == 0 {
		slog.Error("Unexpected response: " + string(buf))
		return
	}
	slog.Info(fmt.Sprintf("Encoded json is: %s", strings.ReplaceAll(string(buf), "\r", "")))
	o, err := obj.Decode(buf)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while converting bencode to json: %v\n", err))
	}
	println(o)
	js, err := json.Decode(buf)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while converting bencode to json: %v\n", err))
	} else {
		slog.Info(fmt.Sprintf("Decoded json is: %s", js))
	}
	var bytes []byte
	bytes, err = json.Encode(js)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while converting bencode to json: %v\n", err))
	} else {
		slog.Info(fmt.Sprintf("Encoded json is: %s", strings.ReplaceAll(string(bytes), "\r", "")))
	}
	js, _ = json.Decode(bytes)
	slog.Info(fmt.Sprintf("Decoded json is: %s", js))
	bytes, err = json.Encode(js)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while converting bencode to json: %v\n", err))
	} else {
		slog.Info(fmt.Sprintf("Encoded json is: %s", strings.ReplaceAll(string(bytes), "\r", "")))
	}
}
