package dht_converters

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//func encodeNodes(json string, index int) ([]byte, int, error) {
//	var result []byte
//	literalEnd := strings.IndexByte(json[index:], '"')
//	literal := json[index : index+literalEnd]
//	parts := strings.Split(literal, "   ")
//	for i, p := range parts {
//		if p == "" {
//			continue
//		}
//		if i%2 == 1 {
//			ip, err := encodeIp(p)
//			if err != nil {
//				return nil, err
//			}
//			result = append(result, ip...)
//		} else {
//			hash, err := encodeHash(p)
//			if err != nil {
//				return nil, err
//			}
//			result = append(result, hash...)
//		}
//	}
//}

func encodeIp(json string, index int) ([]byte, int, error) {
	var result []byte
	literalEnd := strings.IndexByte(json[index:], '"')
	literal := json[index : index+literalEnd]
	parts := strings.Split(strings.ReplaceAll(literal, ":", "."), ".")
	i := 0
	for i < 4 {
		ipPart := parts[i]
		ipInt, err := strconv.Atoi(ipPart) // ahoy
		if err != nil || ipInt > 255 || ipInt < 0 {
			return nil, 0, errors.New(fmt.Sprintf("Problems with parsing ip: '%s', part: '%s'", literal, ipPart))
		}
		result = append(result, byte(ipInt))
		i++
	}
	ipPart := parts[4]
	portInt, err := strconv.Atoi(ipPart) // ahoy
	if err != nil || portInt > 65535 {
		return nil, 0, errors.New(fmt.Sprintf("Problems with parsing ip: '%s', part: '%s'", literal, ipPart))
	}
	portStr := fmt.Sprintf("%x", portInt)
	first, _ := hex.DecodeString(portStr[:2])
	second, _ := hex.DecodeString(portStr[2:])
	result = append(result, first...)
	result = append(result, second...)
	return result, index + literalEnd, nil
}

func encodeHash(literal string) ([]byte, error) {
	var result []byte
	var i int
	if len(literal)%2 != 0 {
		return nil, errors.New(fmt.Sprintf("Even number of ascii symbols while converting to []bytes string: '%s'", literal))
	}
	for 2*i+1 < len(literal) {
		byteAt0, err := hex.DecodeString(literal[2*i : 2*i+2])
		if err != nil || len(byteAt0) > 1 || byteAt0[0] > 255 {
			return nil, errors.New(fmt.Sprintf("Problem with encoding hex string: '%s'", literal))
		}
		result = append(result, byteAt0[0])
		i++
	}
	return result, nil
}

func DHTPackageToJson() {

}