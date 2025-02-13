package bencode

import (
	"encoding/hex"
	"errors"
	"fmt"
	"magnet-parser/utils"
	"strconv"
	"strings"
)

func byteIsDigit(symbol byte) bool {
	return 47 < symbol && symbol < 58
}

func encodeNextElement(json string, index int, key string) ([]byte, int, error) {
	if byteIsDigit(json[index]) {
		return encodeNumber(json, index)
	} else if json[index] == '{' {
		return encodeDict(json, index)
	} else if json[index] == '[' {
		return encodeList(json, index, key)
	} else if json[index] == '"' {
		return encodeStringLiteral(json, index, key)
	} else {
		return nil, 0, errors.New(fmt.Sprintf("No known entity start for index: %d, symbol: '%s'", index, string(json[index])))
	}
}

func encodeIp(literal string) ([]byte, error) {
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

func encodeIpOrHash(literal string, key string) ([]byte, error) {
	if key == "ip" {
		return encodeIp(literal)
	} else if key == "nodes" {
		var result []byte
		parts := strings.Split(literal, "   ")
		for i, p := range(parts) {
			if p == "" {
				continue
			}
			if i % 2 == 1 {
				ip, err := encodeIp(p)
				if err != nil {
					return nil, err
				}
				result = append(result, ip...)
			} else {
				hash, err := encodeHash(p)
				if err != nil {
					return nil, err
				}
				result = append(result, hash...)
			}
		}
		return result, nil
	}
	return encodeHash(literal)
}

func encodeDict(json string, index int) ([]byte, int, error) {
	startIndex := index
	index++
	readingKey := true
	var tempBytes, resultingBytes []byte
	var err error
	var key string
	resultingBytes = append(resultingBytes, 'd')
	for index < len(json) && json[index] != '}' {
		if readingKey && json[index] == ':' {
			readingKey = !readingKey
		} else if readingKey && json[index] == '"' {
			tempBytes, index, err = encodeStringLiteral(json, index, "")
			if err != nil {
				return nil, 0, err
			}
			newKey := string(tempBytes)
			newKey = newKey[strings.IndexByte(newKey, ':')+1:]
			if key != "" && strings.Compare(key, newKey) == 1 {
				return nil, 0, errors.New(fmt.Sprintf("No alphabetical order in dictionary key at index: %d", index))
			}
			key = newKey
			resultingBytes = append(resultingBytes, tempBytes...)
		} else if readingKey && json[index] != ',' {
			return nil, 0, errors.New(fmt.Sprintf("Dictionary key can be nothing but string. Index: %d, Symbol: '%s'", index, string(json[index])))
		} else if !readingKey {
			tempBytes, index, err = encodeNextElement(json, index, key)
			if err != nil {
				return nil, 0, err
			}
			resultingBytes = append(resultingBytes, tempBytes...)
			readingKey = !readingKey
		}
		index++
	}
	if index >= len(json) {
		return nil, 0, errors.New(fmt.Sprintf("No closing '}' symbol for dictionary starting at index: %d", startIndex))
	}
	resultingBytes = append(resultingBytes, 'e')
	return resultingBytes, index, err
}

func encodeList(json string, index int, key string) ([]byte, int, error) {
	startIndex := index
	index++
	var tempBytes, resultingBytes []byte
	var err error
	resultingBytes = append(resultingBytes, 'l')
	for index < len(json) && json[index] != ']' {
		if json[index] != ',' {
			tempBytes, index, err = encodeNextElement(json, index, key)
			resultingBytes = append(resultingBytes, tempBytes...)
			if err != nil {
				return nil, 0, err
			}
		}
		index++
	}
	if index >= len(json) {
		return nil, 0, errors.New(fmt.Sprintf("No closing ']' symbol for list starting at index: %d", startIndex))
	}
	resultingBytes = append(resultingBytes, 'e')
	return resultingBytes, index, err
}

func encodeStringLiteral(json string, index int, key string) ([]byte, int, error) {
	index++
	literalEnd := strings.IndexByte(json[index:], '"')
	literal := json[index : index+literalEnd]
	var resultingBytes, temp []byte
	var err error
	var length int
	if key == "values" {
		temp, err = encodeIp(literal)
		if err != nil {
			return nil, 0, err
		}
		length = len(temp)
		resultingBytes = append(resultingBytes, fmt.Sprintf("%d:", length)...)
		resultingBytes = append(resultingBytes, temp...)
	} else if utils.IsItemInArray(key, []string{"ip", "nodes"}) {
		temp, err = encodeIpOrHash(literal, key)
		if err != nil {
			return nil, 0, err
		}
		length = len(temp)
		resultingBytes = append(resultingBytes, fmt.Sprintf("%d:", length)...)
		resultingBytes = append(resultingBytes, temp...)
	} else {
		length = len(literal)
		resultingBytes = append(resultingBytes, fmt.Sprintf("%d:", length)...)
		resultingBytes = append(resultingBytes, []byte(literal)...)
	}
	return resultingBytes, index + literalEnd, nil
}

func encodeNumber(json string, index int) ([]byte, int, error) {
	result := strings.Builder{}
	result.WriteByte('i')
	for index < len(json) && (byteIsDigit(json[index]) || json[index] == '-') {
		result.WriteByte(json[index])
		index++
	}
	if index < len(json) && json[index] == '.' {
		return nil, 0, errors.New(fmt.Sprintf("Numbers must be integer. Index: %d", index))
	}
	result.WriteByte('e')
	return []byte(result.String()), index-1, nil
}

// a - arguments (string)
// q - method name (string)
// y - message type: r - response, q - query, e - error (string)
// v - version (hex)
// t - transaction id (hex)
// id - id of node (hex)
// ip - ... (hex)
func JSONToBencode(json string) ([]byte, error) {
	var formattedStr string
	formattedStr = strings.ReplaceAll(json, "\t", "")
	var encodedBytes []byte
	var err error
	if len(formattedStr) > 0 {
		encodedBytes, _, err = encodeNextElement(formattedStr, 0, "")
		return encodedBytes, err
	}
	return encodedBytes, err
}
