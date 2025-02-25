package dht

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"magnet-parser/globals"
	"strconv"
	"strings"
)

func encodeNodes(js string) ([]byte, error) {
	var result []byte
	var nodes map[string]string
	err := json.Unmarshal([]byte(js), &nodes)
	if err != nil {
		return nil, err
	}
	for hash, ip := range nodes {
		var encodedIp, encodedHash []byte
		encodedIp, err = encodeIp(ip)
		if err != nil {
			return nil, err
		}
		encodedHash, err = encodeHash(hash)
		if err != nil {
			return nil, err
		}
		result = append(result, encodedIp...)
		result = append(result, encodedHash...)
	}
	return result, nil
}

func encodeIp(literal string) ([]byte, error) { // todo rework
	var result []byte
	parts := strings.Split(strings.ReplaceAll(literal, ":", "."), ".")
	i := 0
	for i < 4 {
		ipPart := parts[i]
		ipInt, err := strconv.Atoi(ipPart) // ahoy
		if err != nil || ipInt > 255 || ipInt < 0 {
			return nil, errors.New(fmt.Sprintf("Problems with parsing ip: '%s', part: '%s'", literal, ipPart))
		}
		result = append(result, byte(ipInt))
		i++
	}
	ipPart := parts[4]
	portInt, err := strconv.Atoi(ipPart) // ahoy
	if err != nil || portInt > 65535 {
		return nil, errors.New(fmt.Sprintf("Problems with parsing ip: '%s', part: '%s'", literal, ipPart))
	}
	portStr := fmt.Sprintf("%x", portInt)
	first, _ := hex.DecodeString(portStr[:2])
	second, _ := hex.DecodeString(portStr[2:])
	result = append(result, first...)
	result = append(result, second...)
	return result, nil
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

func Compress(obj *globals.PackageType) (*globals.PackageType, error) {
	if obj == nil {
		return nil, nil
	}
	if len(obj.V) > 0 {
		V, err := encodeHash(obj.V)
		if err != nil {
			return nil, err
		}
		obj.V = string(V)
	}

	if obj.A != nil {
		if len(obj.A.Id) > 0 {
			Id, err := encodeHash(obj.A.Id)
			if err != nil {
				return nil, err
			}
			obj.A.Id = string(Id)
		}
		if obj.A.InfoHash != nil {
			Hash, err := encodeHash(*obj.A.InfoHash)
			if err != nil {
				return nil, err
			}
			obj.A.Id = string(Hash)
		}
	}

	if obj.R != nil {
		if len(obj.R.Id) > 0 {
			Id, err := encodeHash(obj.R.Id)
			if err != nil {
				return nil, err
			}
			obj.R.Id = string(Id)
		}
	}
	return obj, nil
}