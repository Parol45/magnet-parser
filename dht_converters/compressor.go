package dht_converters

import (
	"fmt"
	"strconv"
	"strings"
)

// a - arguments (string)
// q - method name (string)
// y - message type: r - response, q - query, e - error (string)
// v - version (hex)
// t - transaction id (hex)
// id - id of node (hex)
// ip - ... (hex)
func decodeIp(ipBytes []byte) string {
	ip1, err1 := strconv.ParseInt(fmt.Sprintf("%x", ipBytes[0]), 16, 64)
	ip2, err2 := strconv.ParseInt(fmt.Sprintf("%x", ipBytes[1]), 16, 64)
	ip3, err3 := strconv.ParseInt(fmt.Sprintf("%x", ipBytes[2]), 16, 64)
	ip4, err4 := strconv.ParseInt(fmt.Sprintf("%x", ipBytes[3]), 16, 64)
	port, err5 := strconv.ParseInt(fmt.Sprintf("%x", ipBytes[4])+fmt.Sprintf("%x", ipBytes[5]), 16, 64)
	if err1 == nil && err2 == nil && err3 == nil && err4 == nil && err5 == nil {
		return fmt.Sprintf("%d.%d.%d.%d:%d", ip1, ip2, ip3, ip4, port)
	}
	return ""
}

func decodeHash(bytes []byte) string {
	var result strings.Builder
	for _, b := range bytes {
		str := fmt.Sprintf("%02x", b)
		result.WriteString(str)
	}
	return result.String()
}

func tryDecodeIpOrHash(bytes []byte, key string) string {
	if key == "\"ip\"" {
		return decodeIp(bytes)
	} else if key == "\"nodes\"" {
		var result strings.Builder
		index := 0
		for index < len(bytes) {
			result.WriteString(decodeHash(bytes[index:index+20]))
			result.WriteString("   ")
			index += 20
			result.WriteString(decodeIp(bytes[index:index+6]))
			result.WriteString("   ")
			index += 6
		}
		str := result.String()
		lastIndex := len(str)-3
		if lastIndex < 0 {
			lastIndex = 0
		}
		return str[:lastIndex]
	}
	return decodeHash(bytes)
}