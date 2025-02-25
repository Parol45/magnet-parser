package globals

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func SetupLogger(logFileName string) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   "./logs/" + logFileName + ".log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, //days
		Compress:   true,
	}))
}

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

func TakePointer[T any](item T) *T {
	return &item
}

func ByteIsDigit(symbol byte) bool {
	return 47 < symbol && symbol < 58 || symbol == 45
}