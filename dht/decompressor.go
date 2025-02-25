package dht

import (
	"errors"
	"fmt"
	"magnet-parser/globals"
	"strconv"
	"strings"
)

func decodeIp(ipBytes []byte) (string, error) { // todo rework
	if len(ipBytes) != 6 {
		return "", errors.New("ip has wrong number of bytes")
	}
	ip1, err1 := strconv.ParseInt(fmt.Sprintf("%x", ipBytes[0]), 16, 64)
	ip2, err2 := strconv.ParseInt(fmt.Sprintf("%x", ipBytes[1]), 16, 64)
	ip3, err3 := strconv.ParseInt(fmt.Sprintf("%x", ipBytes[2]), 16, 64)
	ip4, err4 := strconv.ParseInt(fmt.Sprintf("%x", ipBytes[3]), 16, 64)
	port, err5 := strconv.ParseInt(fmt.Sprintf("%x", ipBytes[4])+fmt.Sprintf("%x", ipBytes[5]), 16, 64)
	if err1 == nil && err2 == nil && err3 == nil && err4 == nil && err5 == nil {
		return fmt.Sprintf("%d.%d.%d.%d:%d", ip1, ip2, ip3, ip4, port), nil
	}
	return "", errors.Join(err1, err2, err3, err4)
}

func decodeHash(bytes []byte) string {
	var result strings.Builder
	for _, b := range bytes {
		str := fmt.Sprintf("%02x", b)
		result.WriteString(str)
	}
	return result.String()
}

func decodeNodes(literal string) (res string, err error) {
	var ip string
	var builder strings.Builder
	partN := 0
	builder.WriteByte('{')
	for i := 0; i < len(literal); {
		if partN% 2 == 1 {
			ip, err = decodeIp([]byte(literal[i:i+6]))
			if err != nil {
				return
			}
			builder.WriteByte('"')
			builder.WriteString(ip)
			builder.WriteByte('"')
			i += 6
			if i < len(literal) {
				builder.WriteByte(',')
			}
		} else {
			hash := decodeHash([]byte(literal[i:i+20]))
			builder.WriteByte('"')
			builder.WriteString(hash)
			builder.WriteByte('"')
			builder.WriteByte(':')
			i += 20
		}
		partN++
	}
	builder.WriteByte('}')
	return builder.String(), nil
}

func Decompress(obj *globals.PackageType) (*globals.PackageType, error) {
	if obj == nil {
		return nil, nil
	}
	if len(obj.V) > 0 {
		obj.V = decodeHash([]byte(obj.V))
	}
	if obj.R != nil {
		if len(obj.R.Id) > 0 {
			obj.R.Id = decodeHash([]byte(obj.R.Id))
		}
		if obj.R.Nodes != nil {
			nodes, err := decodeNodes(*obj.R.Nodes)
			if err != nil {
				return nil, err
			}
			obj.R.Nodes = globals.TakePointer(nodes)
		}
		if obj.R.Token != nil {
			obj.R.Token = globals.TakePointer(decodeHash([]byte(*obj.R.Token)))
		}
	}
	return obj, nil
}
