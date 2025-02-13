package bencode

import (
	"errors"
	"fmt"
	"magnet-parser/utils"
	"strconv"
	"strings"
)

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

func decodeNextElement(bytes []byte, index int, key string) (string, int, error) {
	if bytes[index] == 'i' {
		return decodeNumber(bytes, index)
	} else if bytes[index] == 'd' {
		return decodeDict(bytes, index)
	} else if bytes[index] == 'l' {
		return decodeList(bytes, index, key)
	} else if byteIsDigit(bytes[index]) {
		return decodeStringLiteral(bytes, index, key)
	} else {
		return "", len(bytes), errors.New(fmt.Sprintf("No known entity start for index: %d, symbol: '%s'", index, string(bytes[index])))
	}
}

func decodeList(bytes []byte, index int, key string) (string, int, error) {
	startIndex := index
	index++
	builder := strings.Builder{}
	builder.WriteString("[")
	var tempStr string
	var err error
	for index < len(bytes) && bytes[index] != 'e' {
		tempStr, index, err = decodeNextElement(bytes, index, key)
		if err != nil {
			return "", 0, err
		}
		builder.WriteString(tempStr + ",")
	}
	if index >= len(bytes) {
		return "", 0, errors.New(fmt.Sprintf("No closing 'e' symbol for list starting at index: %d, symbol: '%s'", startIndex, string(bytes[startIndex])))
	}
	temp := builder.String()
	temp = temp[:len(temp)-1] + "]"
	return temp, index + 1, err
}

func decodeNumber(bytes []byte, index int) (string, int, error) {
	startIndex := index
	index++
	result := strings.Builder{}
	for index < len(bytes) && bytes[index] != 'e' {
		result.WriteByte(bytes[index])
		index++
	}
	if index >= len(bytes) {
		return "", 0, errors.New(fmt.Sprintf("No closing 'e' symbol for number starting at index: %d, symbol: '%s'", startIndex, string(bytes[startIndex])))
	}
	return result.String(), index + 1, nil
}

func decodeStringLiteral(bytes []byte, index int, key string) (string, int, error) {
	strLen := 0
	for byteIsDigit(bytes[index]) {
		strLen = strLen*10 + int(bytes[index]-'0')
		index++
	}
	if bytes[index] != ':' {
		return "", 0, errors.New(fmt.Sprintf("Missing string semicolon. Index: %d, Symbol: '%s'", index, string(bytes[index])))
	} else if index+strLen > len(bytes) {
		return "", 0, errors.New(fmt.Sprintf("Wrong string length. Index: %d, Symbol: '%s'", index, string(bytes[index])))
	}
	index++
	var resultingStr string
	if key == `"values"` {
		resultingStr = "\"" + decodeIp(bytes[index:index+strLen]) + "\""
	} else if !utils.IsItemInArray(key, []string{"\"y\"", "\"q\"", "", "\"e\""}) {
		resultingStr = "\"" + tryDecodeIpOrHash(bytes[index:index+strLen], key) + "\""
	} else {
		resultingStr = "\"" + string(bytes[index:index+strLen]) + "\""
	}
	return resultingStr, index + strLen, nil
}

func decodeDict(bytes []byte, index int) (string, int, error) {
	startIndex := index
	builder := strings.Builder{}
	builder.WriteString("{")
	index++
	readingKey := true
	var tempStr, key string
	var err error
	for index < len(bytes) && bytes[index] != 'e' {
		if readingKey && byteIsDigit(bytes[index]) {
			key, index, err = decodeStringLiteral(bytes, index, "")
			if err != nil {
				return "", 0, err
			}
			builder.WriteString(key + ":")
			readingKey = !readingKey
		} else if readingKey {
			return "", 0, errors.New(fmt.Sprintf("Dictionary key can be nothing but string. Index: %d, Symbol: '%s'", index, string(bytes[index])))
		} else if !readingKey {
			tempStr, index, err = decodeNextElement(bytes, index, key)
			if err != nil {
				return "", 0, err
			}
			builder.WriteString(tempStr + ",")
			readingKey = !readingKey
		}
	}
	if index >= len(bytes) {
		return "", 0, errors.New(fmt.Sprintf("No closing 'e' symbol for dictionary starting at index: %d, symbol: '%s'", startIndex, string(bytes[startIndex])))
	}
	temp := builder.String()
	temp = temp[:len(temp)-1] + "}"
	return temp, index + 1, err
}

func BencodeToJSON(encodedStr []byte) (string, error) {
	var decodedStr string
	var err error
	if len(encodedStr) > 0 {
		decodedStr, _, err = decodeNextElement(encodedStr, 0, "")
		return decodedStr, err
	}
	return decodedStr, err
}