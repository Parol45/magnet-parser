package utils

import (
	"net"
	"strconv"
	"strings"
)

func IsItemInArray[T comparable](val T, array []T) bool {
	for _, v := range array {
		if val == v {
			return true
		}
	}
	return false
}

func StringToUDPAddr(ip string) *net.UDPAddr {
	ipStrings := strings.Split(ip, ".")
	byteAndPort := strings.Split(ipStrings[3], ":")
	ipStrings[3] = byteAndPort[0]
	var ipBytes []byte
	i := 0
	for i < 4 {
		integer, _ := strconv.Atoi(ipStrings[i])
		ipBytes = append(ipBytes, byte(integer))
		i++
	}
	port, _ := strconv.Atoi(byteAndPort[1])
	return &net.UDPAddr{IP:ipBytes,Port:port,Zone:""}
}